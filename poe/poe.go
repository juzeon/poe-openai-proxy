package poe

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/juzeon/poe-openai-proxy/conf"
	"github.com/juzeon/poe-openai-proxy/util"
	"sync"
	"time"
)

var httpClient *resty.Client
var clients []*Client
var clientIx = 0
var clientLock = &sync.Mutex{}

func Setup() {
	httpClient = resty.New().SetBaseURL(conf.Conf.Gateway).SetTimeout(60 * time.Second)
	for _, token := range conf.Conf.Tokens {
		client, err := NewClient(token)
		if err != nil {
			panic(err)
		}
		clients = append(clients, client)
	}
}

type Client struct {
	Token string
}

func NewClient(token string) (*Client, error) {
	util.Logger.Info("registering client: " + token)
	resp, err := httpClient.R().SetFormData(map[string]string{
		"token": token,
	}).Post("/add_token")
	if err != nil {
		return nil, errors.New("registering client error: " + err.Error())
	}
	util.Logger.Info("registering client: " + resp.String())
	return &Client{Token: token}, nil
}
func (c Client) getContentToSend(messages []Message) string {
	leadingMap := map[string]string{
		"system":    "Instructions",
		"user":      "User",
		"assistant": "Assistant",
	}
	content := ""
	for _, message := range messages {
		content += "||>" + leadingMap[message.Role] + ":\n" + message.Content + "\n"
	}
	content += "||>Assistant:\n"
	util.Logger.Debug("Generated content to send: " + content)
	return content
}
func (c Client) Stream(messages []Message) (<-chan string, error) {
	channel := make(chan string, 1024)
	content := c.getContentToSend(messages)
	conn, _, err := websocket.DefaultDialer.Dial(conf.Conf.GetGatewayWsURL()+"/stream", nil)
	if err != nil {
		return nil, err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(c.Token))
	if err != nil {
		return nil, err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(conf.Conf.Bot))
	if err != nil {
		return nil, err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(content))
	if err != nil {
		return nil, err
	}
	go func(conn *websocket.Conn, channel chan string) {
		defer close(channel)
		defer conn.Close()
		for {
			_, v, err := conn.ReadMessage()
			channel <- string(v)
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					util.Logger.Error(err)
				}
				channel <- "[DONE]"
				break
			}
		}
	}(conn, channel)
	return channel, nil
}
func (c Client) Ask(messages []Message) (*Message, error) {
	content := c.getContentToSend(messages)
	resp, err := httpClient.R().SetFormData(map[string]string{
		"token":   c.Token,
		"bot":     conf.Conf.Bot,
		"content": content,
	}).Post("/ask")
	if err != nil {
		return nil, err
	}
	return &Message{
		Role:    "assistant",
		Content: resp.String(),
		Name:    "",
	}, nil
}

func GetClient() *Client {
	clientLock.Lock()
	defer clientLock.Unlock()
	client := clients[clientIx%len(clients)]
	clientIx++
	return client
}
