package conf

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigStruct struct {
	Port          int
	Tokens        []string
	Bot           map[string]string
	SimulateRoles int
	RateLimit     int
	CoolDown      int
	Timeout       int
}

type ModelDef struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelsResp struct {
	Object string     `json:"object"`
	Data   []ModelDef `json:"data"`
}

var Conf ConfigStruct

var Models ModelsResp

func loadEnvVar(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func loadEnvVarAsInt(key string, defaultValue int) int {
	valueStr := loadEnvVar(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func loadEnvVarAsSlice(key string) []string {
	valueStr := loadEnvVar(key, "")
	return strings.Split(valueStr, ",")
}

func Setup() {
	// Load environment variables from .env file (for development purposes)
	_ = godotenv.Load()

	Conf.Port = loadEnvVarAsInt("PORT", 8080)
	Conf.Tokens = loadEnvVarAsSlice("TOKENS")
	Conf.SimulateRoles = loadEnvVarAsInt("SIMULATE_ROLES", 2)
	Conf.RateLimit = loadEnvVarAsInt("RATE_LIMIT", 10)
	Conf.CoolDown = loadEnvVarAsInt("COOL_DOWN", 5)
	Conf.Timeout = loadEnvVarAsInt("TIMEOUT", 60)

	Conf.Bot = map[string]string{
		"gpt-3.5-turbo":       "chinchilla",
		"gpt-4":               "beaver",
		"gpt-3.5-turbo-0301":  "a2",
		"gpt-4-32k":           "a2_100k",
		"gpt-4-0314":          "a2_2",
		"Sage":                "capybara",
		"ChatGPT":             "chinchilla",
		"GPT-4":               "beaver",
		"Claude-instant":      "a2",
		"Claude-instant-100k": "a2_100k",
		"Claude+":             "a2_2",
	}

	Models.Object = ""

	for key := range Conf.Bot {
		Models.Data = append(Models.Data, ModelDef{
			ID:      key,
			Object:  "",
			Created: 0,
			OwnedBy: "",
		})
	}
}