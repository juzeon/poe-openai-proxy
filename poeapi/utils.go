package poeapi

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func GetFinalResponse(ch <-chan map[string]interface{}) string {
	var m map[string]interface{}
	for message := range ch {
		m = message
		if message["state"] != "complete" {
			continue
		}
		return message["text"].(string)
	}
	return m["text"].(string)
}

func GetTextStream(ch <-chan map[string]interface{}) <-chan string {
	stream := make(chan string, 1)
	go func() {
		for message := range ch {
			stream <- message["text_new"].(string)
		}
		close(stream)
	}()
	return stream
}

func generatePayload(queryName string, variables map[string]interface{}) interface{} {
	if queryName == "recv" {
		if rand.Float64() > 0.9 {
			return []map[string]interface{}{
				{
					"category": "poe/bot_response_speed",
					"data":     variables,
				},
				{
					"category": "poe/statsd_event",
					"data": map[string]interface{}{
						"key":        "poe.speed.web_vitals.INP",
						"value":      rand.Intn(26) + 100,
						"category":   "time",
						"path":       "/[handle]",
						"extra_data": map[string]interface{}{},
					},
				},
			}
		} else {
			return []map[string]interface{}{
				{
					"category": "poe/bot_response_speed",
					"data":     variables,
				},
			}
		}
	}
	return map[string]interface{}{
		"query":     queries[queryName],
		"variables": variables,
	}
}

func generateNonce(length int) string {
	if length == 0 {
		length = 16
	}
	const lettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var nonce = make([]rune, length)
	for i := range nonce {
		nonce[i] = rune(lettersAndDigits[rand.Intn(len(lettersAndDigits))])
	}
	return string(nonce)
}

func getConfigPath() string {
	var configPath string
	if os.PathSeparator == '\\' {
		configPath = filepath.Join(os.Getenv("APPDATA"), "poe-api")
	} else {
		configPath = filepath.Join(os.Getenv("HOME"), ".config", "poe-api")
	}
	return configPath
}

func setSavedDeviceID(userID, deviceID string) {
	deviceIDPath := filepath.Join(getConfigPath(), "device_id.json")
	deviceIDs := make(map[string]string)

	if _, err := os.Stat(deviceIDPath); !os.IsNotExist(err) {
		deviceIDBytes, err := os.ReadFile(deviceIDPath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(deviceIDBytes, &deviceIDs)
		if err != nil {
			panic(err)
		}
	}

	deviceIDs[userID] = deviceID
	err := os.MkdirAll(filepath.Dir(deviceIDPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	deviceIDBytes, err := json.MarshalIndent(deviceIDs, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(deviceIDPath, deviceIDBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func getSavedDeviceID(userID string) string {
	deviceIDPath := filepath.Join(getConfigPath(), "device_id.json")
	deviceIDs := make(map[string]string)

	if _, err := os.Stat(deviceIDPath); !os.IsNotExist(err) {
		deviceIDBytes, err := os.ReadFile(deviceIDPath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(deviceIDBytes, &deviceIDs)
		if err != nil {
			panic(err)
		}
	}

	if deviceID, ok := deviceIDs[userID]; ok {
		return deviceID
	}

	deviceID := uuid.New().String()
	deviceIDs[userID] = deviceID
	err := os.MkdirAll(filepath.Dir(deviceIDPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	deviceIDBytes, err := json.MarshalIndent(deviceIDs, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(deviceIDPath, deviceIDBytes, 0644)
	if err != nil {
		panic(err)
	}

	return deviceID
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func reverseSlice(s []map[string]interface{}) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func containKey(key string, m map[string]interface{}) bool {
	_, ok := m[key]
	return ok
}
