package router

import (
	"beegotest/controllers"
	"beegotest/filter"
	"beegotest/service"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type routeInfo struct {
	method  string
	pattern string
	handler echo.HandlerFunc
	filters []echo.MiddlewareFunc
}

func InitRoutes(e *echo.Echo) {
	initMiddleware(e)

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

func initMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	// use logger middleware
	e.Use(middleware.Logger())
}
