package main

import (
	"flag"

	"github.com/labstack/echo"
	"github.com/silentred/kassadin"
	"github.com/silentred/kassadin/db"
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

func initService(app *kassadin.App) error {
	mm := db.NewMysqlManager(app, app.Config.Mysql)
	if mm != nil {
		app.Set("mysql", mm, nil)
	}

	return nil
}

func initRoute(app *kassadin.App) error {
	app.Route.GET("/", func(ctx echo.Context) error {
		return ctx.String(200, "hello world")
	})

	return nil
}
