package conf

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"strings"
)

type ConfigStruct struct {
	Port          int      `toml:"port"`
	Tokens        []string `toml:"tokens"`
	Gateway       string   `toml:"gateway"`
	Bot           string   `toml:"bot"`
	SimulateRoles int      `toml:"simulate-roles"`
	RateLimit     int      `toml:"rate-limit"`
	CoolDown      int      `toml:"cool-down"`
	Timeout       int      `toml:"timeout"`
}

func (c ConfigStruct) GetGatewayWsURL() string {
	str := strings.ReplaceAll(c.Gateway, "http://", "ws://")
	str = strings.ReplaceAll(str, "https://", "wss://")
	return str
}

var Conf ConfigStruct

func Setup() {
	v, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(v, &Conf)
	if err != nil {
		panic(err)
	}
	if Conf.Port == 0 {
		Conf.Port = 3700
	}
	if Conf.Bot == "" {
		Conf.Bot = "capybara"
	}
	if Conf.RateLimit == 0 {
		Conf.RateLimit = 10
	}
}
