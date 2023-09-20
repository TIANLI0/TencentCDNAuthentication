package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	loadconfig()
	loadWhitelist()
	loadBlacklist()
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// 添加CORS中间件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// 设置请求处理函数
	router.Any("/", handleRequest)
	router.Any("/:path/*filepath", handleRequest)

	router.Run(":5276")
}
