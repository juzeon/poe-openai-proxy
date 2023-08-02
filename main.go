package main

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juzeon/poe-openai-proxy/conf"
	"github.com/juzeon/poe-openai-proxy/poe"
	"github.com/juzeon/poe-openai-proxy/router"
)

func main() {
	conf.Setup()
	poe.Setup()
	engine := gin.Default()
	router.Setup(engine)

	// 检查 ssl 目录下的证书文件是否存在
	_, certErr := os.Stat("ssl/cert.pem")
	_, keyErr := os.Stat("ssl/key.pem")

	if certErr == nil && keyErr == nil {
		// SSL 证书和私钥文件都存在，启动 HTTPS 服务器
		err := engine.RunTLS(":"+strconv.Itoa(conf.Conf.Port), "ssl/cert.pem", "ssl/key.pem")
		if err != nil {
			panic(err)
		}
	} else {
		// SSL 证书或私钥文件不存在，启动 HTTP 服务器
		err := engine.Run(":" + strconv.Itoa(conf.Conf.Port))
		if err != nil {
			panic(err)
		}
	}
}
