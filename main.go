package main

import (
	"awesomeProject/handlers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"awesomeProject/extension/global"
)


func main() {
	global.Conf = global.NewConf()
	global.RedisPool = global.InitRedisPool()
	app:=iris.New()
	app.Use(logger.New())
	app.Get("/adMobSuccess", handlers.AdMobSuccessGetHandler) // 接收admob 的回调 接口
	app.Get("/getJson", handlers.GetJson)
	app.Get("/getRiskControl", handlers.GetRiskControl)
	app.Get("/nextRound", handlers.NextRound)
	_ = app.Run(iris.Addr(":9080"), iris.WithoutServerError(iris.ErrServerClosed))
}
