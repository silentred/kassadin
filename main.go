package main

import (
	"path/filepath"

	"github.com/labstack/echo"
	elog "github.com/silentred/echo-log"
	"github.com/silentred/template/router"
	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/spf13/viper"
)

// Echo is the web engine
var Echo *echo.Echo

func init() {
	Echo = echo.New()
	initConfig()
	initLogger()
}

func main() {
	service.InitDBInfo()
	service.InitMysqlORM(service.MysqlConfig)

	router.InitRoutes(Echo)
	router.InitMiddleware(Echo)

	Echo.Start(":8090")
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath(util.SelfDir())
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot find config file in working directory.")
	}
}

func initLogger() {
	rotate := viper.GetBool("app.logRotate")
	provider := viper.GetString("app.logProvider")
	mode := viper.GetString("app.runMode")
	appName := viper.GetString("app.name")
	if rotate && provider == "file" {
		path := filepath.Join(util.SelfDir(), "storage", "log")
		limitSize := 100 << 20 // 100 MB
		Echo.Logger = elog.NewLogger(path, appName, limitSize)
	}

	switch mode {
	case "dev":
		Echo.Logger.SetLevel(elog.DEBUG)
	case "prod":
		Echo.Logger.SetLevel(elog.WARN)
	}

}
