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
	userSV := service.NewUserSV()
	ituneSV := service.NewItunesSV(service.AdToken)

	user := controllers.NewUserController(userSV, ituneSV)

	// custom middleware
	metrics := filter.Metrics()

	// v1 group
	v1Group := e.Group("/promotion")
	v1Group.Use(metrics)

	routes := []routeInfo{
		{echo.POST, "/generatelink", user.GenerateLink, nil},
		{echo.POST, "/getpoints", user.GetPoint, nil},
		{echo.POST, "/usepoints", user.UsePoint, nil},
		{echo.POST, "/log", user.Log, nil},
	}

	applyGroupRoutes(v1Group, routes)

	e.GET("/ping", func(c echo.Context) error {
		c.String(200, "")
		return nil
	})
}

func applyGroupRoutes(g *echo.Group, routes []routeInfo) {
	for _, route := range routes {
		switch route.method {
		case echo.GET:
			g.GET(route.pattern, route.handler, route.filters...)
		case echo.POST:
			g.POST(route.pattern, route.handler, route.filters...)
		case echo.PUT:
			g.PUT(route.pattern, route.handler, route.filters...)
		case echo.DELETE:
			g.DELETE(route.pattern, route.handler, route.filters...)
		}
	}
}

func InitMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	// use logger middleware
	// middleware.DefaultLoggerConfig.Output = e.Logger.Output()
	// e.Use(middleware.Logger())
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
