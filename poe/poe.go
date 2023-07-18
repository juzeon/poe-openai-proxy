package poe

import (
	"errors"
	"sync"
	"time"

	"github.com/juzeon/poe-openai-proxy/conf"
	"github.com/juzeon/poe-openai-proxy/util"
	poe_api "github.com/lwydyby/poe-api"
)

var clients []*Client
var clientLock sync.Mutex
var clientIx = 0
var invalidError *poe_api.InvalidToken

func Setup() {
	for i, token := range conf.Conf.Tokens {
		client, err := NewClient(i, token)
		if err != nil {
			panic(err)
		}
		clients = append(clients, client)
	}
}

type Client struct {
	Token  string
	client *poe_api.Client
	Usage  []time.Time
	Lock   bool
	index  int
}

func NewClient(i int, token string) (*Client, error) {
	util.Logger.Info("registering client: " + token)
	client := poe_api.NewClient(token, nil)
	return &Client{index: i, Token: token, Usage: nil, Lock: false, client: client}, nil
}
func (c *Client) getContentToSend(messages []Message) string {
	leadingMap := map[string]string{
		"system":    "Instructions",
		"user":      "User",
		"assistant": "Assistant",
	}
	content := ""
	var simulateRoles bool
	switch conf.Conf.SimulateRoles {
	case 0:
		simulateRoles = false
	case 1:
		simulateRoles = true
	case 2:
		if len(messages) == 1 && messages[0].Role == "user" ||
			len(messages) == 1 && messages[0].Role == "system" ||
			len(messages) == 2 && messages[0].Role == "system" && messages[1].Role == "user" {
			simulateRoles = false
		} else {
			simulateRoles = true
		}
	}
	for _, message := range messages {
		if simulateRoles {
			content += "||>" + leadingMap[message.Role] + ":\n" + message.Content + "\n"
		} else {
			content += message.Content + "\n"
		}
	}
	if simulateRoles {
		content += "||>Assistant:\n"
	}
	util.Logger.Debug("Generated content to send: " + content)
	return content
}
func (c *Client) Stream(messages []Message, model string) (<-chan string, error) {
	defer func() {
		if err := recover(); err != nil {
			if errors.As(err.(error), &invalidError) {
				removeClient(c)
			}
		}
	}()
	channel := make(chan string, 1024)
	content := c.getContentToSend(messages)
	bot, ok := conf.Conf.Bot[model]
	if !ok {
		bot = "capybara"
	}
	util.Logger.Info("Stream using bot", bot)
	resp, err := c.client.SendMessage(bot, content, true, time.Duration(conf.Conf.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(channel)
		defer func() {
			if err := recover(); err != nil {
				channel <- "\n\n[ERROR] " + err.(error).Error()
			}
		}()
		for message := range poe_api.GetTextStream(resp) {
			channel <- message
		}
		channel <- "[DONE]"
	}()
	return channel, nil
}
func (c *Client) Ask(messages []Message, model string) (*Message, error) {
	defer func() {
		if err := recover(); err != nil {
			if errors.As(err.(error), &invalidError) {
				removeClient(c)
			}
		}
	}()
	content := c.getContentToSend(messages)

	bot, ok := conf.Conf.Bot[model]
	if !ok {
		bot = "capybara"
	}
	util.Logger.Info("Ask using bot", bot)

	resp, err := c.client.SendMessage(bot, content, true, time.Duration(conf.Conf.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	return &Message{
		Role:    "assistant",
		Content: poe_api.GetFinalResponse(resp),
		Name:    "",
	}, nil
}
func (c *Client) Release() {
	clientLock.Lock()
	defer clientLock.Unlock()
	c.Lock = false
}

func removeClient(client *Client) {
	clientLock.Lock()
	defer clientLock.Unlock()
	util.Logger.Error("remove client: " + client.Token)
	clients = append(clients[:client.index], clients[client.index+1:]...)
}

func GetClient() (*Client, error) {
	clientLock.Lock()
	defer clientLock.Unlock()
	if len(clients) == 0 {
		return nil, errors.New("no client is available")
	}
	for i := 0; i < len(clients); i++ {
		client := clients[clientIx%len(clients)]
		clientIx++
		if client.Lock {
			continue
		}
		if len(client.Usage) > 0 {
			lastUsage := client.Usage[len(client.Usage)-1]
			if time.Since(lastUsage) < time.Duration(conf.Conf.CoolDown)*time.Second {
				continue
			}
		}
		if len(client.Usage) < conf.Conf.RateLimit {
			client.Usage = append(client.Usage, time.Now())
			client.Lock = true
			return client, nil
		} else {
			usage := client.Usage[len(client.Usage)-conf.Conf.RateLimit]
			if time.Since(usage) <= 1*time.Minute {
				continue
			}
			client.Usage = append(client.Usage, time.Now())
			client.Lock = true
			return client, nil
		}
	}
	return nil, errors.New("no available client")
}
