package main

import (
	"github.com/gin-gonic/gin"
	"github.com/juzeon/poe-openai-proxy/conf"
	"github.com/juzeon/poe-openai-proxy/poe"
	"github.com/juzeon/poe-openai-proxy/router"
	"strconv"
	 
)

func main() {
	conf.Setup()
	poe.Setup()
	engine := gin.Default()
	router.Setup(engine)
	err := engine.Run(":" + strconv.Itoa(conf.Conf.Port))
	if err != nil {
		panic(err)
	}
}
