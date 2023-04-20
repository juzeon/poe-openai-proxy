package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/juzeon/poe-openai-proxy/poe"
	"github.com/juzeon/poe-openai-proxy/util"
	"io"
	"net/http"
	"time"
)

func Setup(engine *gin.Engine) {
	engine.POST("/v1/chat/completions", func(c *gin.Context) {
		var req poe.CompletionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, "bad request")
			return
		}
		client := poe.GetClient()
		if req.Stream {
			util.Logger.Info("stream using client: " + client.Token)
			Stream(c, req, client)
		} else {
			util.Logger.Info("ask using client: " + client.Token)
			Ask(c, req, client)
		}
	})
}
func Stream(c *gin.Context, req poe.CompletionRequest, client *poe.Client) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	w := c.Writer
	flusher, _ := w.(http.Flusher)
	ticker := time.NewTimer(60 * time.Second)
	defer ticker.Stop()
	channel, err := client.Stream(req.Messages)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	createSSEResponse := func(role string, content string, done bool) {
		finishReason := ""
		if done {
			finishReason = "stop"
		}
		data := poe.CompletionSSEResponse{
			Choices: []poe.SSEChoice{{
				Index: 0,
				Delta: poe.Delta{
					Role:    role,
					Content: content,
				},
				FinishReason: finishReason,
			}},
			Created: time.Now().Unix(),
			Id:      "chatcmpl-" + util.RandStringRunes(8),
			Model:   req.Model,
			Object:  "chat.completion.chunk",
		}
		dataV, _ := json.Marshal(&data)
		_, err := io.WriteString(w, "data: "+string(dataV)+"\r\n\r\n")
		if err != nil {
			util.Logger.Error(err)
		}
		flusher.Flush()
		if done {
			_, err := io.WriteString(w, "data: "+string("[DONE]")+"\r\n\r\n")
			if err != nil {
				util.Logger.Error(err)
			}
			flusher.Flush()
		}
	}
	createSSEResponse("assistant", "", false)
forLoop:
	for {
		select {
		case <-ticker.C:
			c.SSEvent("error", "timeout")
			break forLoop
		case d := <-channel:
			if d == "[DONE]" {
				createSSEResponse("", "", true)
				break forLoop
			}
			createSSEResponse("", d, false)
		}
	}
}
func Ask(c *gin.Context, req poe.CompletionRequest, client *poe.Client) {
	message, err := client.Ask(req.Messages)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, poe.CompletionResponse{
		ID:      "chatcmpl-" + util.RandStringRunes(8),
		Object:  "chat.completion",
		Created: int(time.Now().Unix()),
		Choices: []poe.Choice{{
			Index:        0,
			Message:      *message,
			FinishReason: "stop",
		}},
		Usage: poe.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	})
}
