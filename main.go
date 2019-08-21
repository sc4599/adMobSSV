package main

import (
	"awesomeProject/handlers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)


func main() {
	app:=iris.New()
	app.Use(logger.New())
	app.Get("/adMobSuccess", handlers.AdMobSuccessGetHandler) // 接收admob 的回调 接口
	app.Get("/testAdMob", handlers.TestAdMobHandler) // 接收admob 的回调 接口
	_ = app.Run(iris.Addr(":9080"), iris.WithoutServerError(iris.ErrServerClosed))
}
