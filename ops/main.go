package main

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/conf"
	_ "github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/opslog"
	"github.com/chujieyang/commonops/ops/presentation"
	"github.com/chujieyang/commonops/ops/presentation/middleware"
	_ "github.com/chujieyang/commonops/ops/cron"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(middleware.OptionsMiddleware)
	presentation.RegisterRouter(engine)
	if err := engine.Run(":" + conf.Port); err != nil {
		panic(fmt.Errorf("服务启动异常: %s", err))
	}
	opslog.Info().Println("commonops start success!")
}
