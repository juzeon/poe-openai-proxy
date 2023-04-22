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
}
