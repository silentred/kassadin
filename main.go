package main

import (
	"github.com/silentred/template/router"
	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/labstack/echo"
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
	service.InitDBs()
	service.InitServices()

	router.InitRoutes(Echo)
	router.InitMiddleware(Echo)

	Echo.Start(":8090")
}
