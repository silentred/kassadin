package main

import (
	"flag"

	"github.com/labstack/echo"
	"github.com/silentred/kassadin"
	"github.com/silentred/kassadin/example/model"
	"github.com/golang/go/src/pkg/net/http"

)

func main() {
	flag.Parse()

	app := kassadin.NewApp()
	app.RegisterConfigHook(initConfig)
	app.RegisterRouteHook(initRoute)
	app.Start()
}

func initConfig(app *kassadin.App) error {
	return nil
}

func initRoute(app *kassadin.App) error {
	app.Route.GET("/", func(ctx echo.Context) error {

		return ctx.String(200, "hello world")
	})
	app.Route.GET("/user/id/:id", getUser)
	return nil
}


func getUser(app *kassadin.App, c echo.Context) error {
	user:= model.User{}

	c.JSON(http.StatusOK, user)
}