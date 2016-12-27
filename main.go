package main

import (
	"beegotest/router"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

var Echo *echo.Echo

func init() {
	Echo = echo.New()
}

func main() {
	initLogger()
	router.InitRoutes(Echo)

	Echo.Start(":8090")
}

func initConfig() {
}

func initLogger() {
	Echo.Debug = true
	Echo.Logger.SetLevel(log.DEBUG)
}
