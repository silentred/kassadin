package router

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	session "github.com/silentred/echo-session"
	"github.com/silentred/template/controllers"
	"github.com/silentred/template/filter"
	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/spf13/viper"
)

type routeInfo struct {
	method  string
	pattern string
	handler echo.HandlerFunc
	filters []echo.MiddlewareFunc
}

func InitRoutes(e *echo.Echo) {
	// initialize controlllers
	user := controllers.NewUserController(&service.UserSV{})

	// custom middleware
	metrics := filter.Metrics()

	// v1 group
	v1Group := e.Group("/v1")
	routes := []routeInfo{
		{echo.GET, "/user/:id", user.GetByID, []echo.MiddlewareFunc{metrics}},
		{echo.POST, "/user/:id", user.GetByID, nil},
		{echo.PUT, "/user/:id", user.GetByID, nil},
		{echo.DELETE, "/user/:id", user.GetByID, nil},
	}

	for _, route := range routes {
		// if route.filters == nil {
		// 	route.filters = []echo.MiddlewareFunc{}
		// }

		switch route.method {
		case echo.GET:
			v1Group.GET(route.pattern, route.handler, route.filters...)
		case echo.POST:
			v1Group.POST(route.pattern, route.handler, route.filters...)
		case echo.PUT:
			v1Group.PUT(route.pattern, route.handler, route.filters...)
		case echo.DELETE:
			v1Group.DELETE(route.pattern, route.handler, route.filters...)
		}
	}

}

func InitMiddleware(e *echo.Echo) {
	// use logger middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// session middleware
	setupSession(e)
}

func setupSession(e *echo.Echo) {
	var enable bool
	var provider string
	enable = viper.GetBool("app.sessionEnable")
	provider = viper.GetString("app.sessionProvider")

	if enable {
		switch provider {
		case "file":
			err := setupFileSession(e)
			if err != nil {
				panic(err)
			}
		default:
			panic("only support file provider for session")
		}
	}
}

func setupFileSession(e *echo.Echo) error {
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	sessionPath := filepath.Join(workDir, "storage", "session")
	if !util.FileExists(sessionPath) {
		fmt.Printf("path not exists: %s, use /tmp instead", sessionPath)
		sessionPath = "/tmp"
	}

	store := session.NewFileSystemStoreStore(sessionPath, []byte("secret"))

	e.Use(session.Sessions("GSESSION", store))

	return nil
}
