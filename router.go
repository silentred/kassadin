package main

import (
	"beegotest/filter"

	"github.com/labstack/echo/middleware"
)

func initRoutes() {
	initMiddleware()

	// v1 group
	v1Group := Echo.Group("/v1")

	user := initControllers()
	v1Group.GET("/user/:id", user.GetByID, filter.Metrics())
}

func initMiddleware() {
	// use logger middleware
	Echo.Use(middleware.Logger())
}
