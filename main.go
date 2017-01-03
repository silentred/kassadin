package main

import (
	"github.com/labstack/echo"
	"github.com/silentred/template/router"
	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
)

// Echo is the web engine
var Echo *echo.Echo

func init() {
	Echo = echo.New()
	util.InitConfig()
	util.InitLogger(Echo)
}

func main() {
	service.InitDBInfo()
	//service.InitMysqlORM(service.MysqlConfig)
	service.InitRedisClient(service.RedisConfig)

	router.InitRoutes(Echo)
	router.InitMiddleware(Echo)

	Echo.Start(":8090")
}
