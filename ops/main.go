package main

import (
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/controllers"
	_ "github.com/chujieyang/commonops/ops/cron"
	"github.com/gin-gonic/gin"
)

func optionsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
		"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, "+
		"Token, Language, From, Cookie, OperationCode, ClusterId")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	// 处理 react axios 的 options 请求
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func main() {
	engine := gin.Default()
	engine.Use(optionsMiddleware)
	controllers.RegisterRouter(engine)
	err := engine.Run(":" + conf.Port)
	if err != nil {
		panic("启动异常: " + err.Error())
	}
}
