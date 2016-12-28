package main

import (
	"beegotest/router"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// Echo is the web engine
var Echo *echo.Echo

func init() {
	Echo = echo.New()
	initConfig()
}

func main() {
	initLogger()
	router.InitRoutes(Echo)

	Echo.Start(":8090")
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot find config file in working directory.")
	}
}

func initLogger() {
	Echo.Debug = true
	Echo.Logger.SetLevel(log.DEBUG)
}
