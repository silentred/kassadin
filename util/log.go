package util

import (
	"github.com/labstack/echo"
	elog "github.com/silentred/echo-log"
)

var Logger echo.Logger

func InitFileLogger(path, appName string, limitSize int) {
	Logger = elog.NewLogger(path, appName, limitSize)
}
