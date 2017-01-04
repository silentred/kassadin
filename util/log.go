package util

import (
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/silentred/rotator"
	"github.com/spf13/viper"
)

var Logger echo.Logger

func setFileRotatorWriter(e *echo.Echo, path, appName string, limitSize int) {
	r := rotator.NewFileSizeRotator(path, appName, "log", limitSize)
	l := e.Logger
	l.SetOutput(r)
	Logger = l
}

func InitLogger(e *echo.Echo) {
	rotate := viper.GetBool("app.logRotate")
	logLimit := viper.GetString("app.logLimit")
	provider := viper.GetString("app.logProvider")
	mode := viper.GetString("app.runMode")
	appName := viper.GetString("app.name")
	if rotate && provider == "file" {
		path := filepath.Join(SelfDir(), "storage", "log")
		limitSize, err := ParseByteSize(logLimit) // 100 MB
		if err != nil {
			panic(err)
		}
		if e != nil {
			setFileRotatorWriter(e, path, appName, limitSize)
		}
	} else {
		Logger = e.Logger
	}

	if e != nil {
		switch mode {
		case "dev":
			e.Logger.SetLevel(log.DEBUG)
		case "prod":
			e.Logger.SetLevel(log.INFO)
		default:
			e.Logger.SetLevel(log.DEBUG)
		}
	}
}
