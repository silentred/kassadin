package main

import (
	"flag"

	"github.com/golang/go/src/pkg/net/http"
	"github.com/labstack/echo"
	"github.com/silentred/kassadin"
	"github.com/silentred/kassadin/example/model"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/silentred/kassadin/db"
	"github.com/silentred/kassadin/redis"
)

func main() {
	flag.Parse()

	app := kassadin.NewApp()
	app.RegisterConfigHook(initConfig)
	app.RegisterRouteHook(initRoute)
	app.Start()
}

func initConfig(app *kassadin.App) error {
	configName := app.GetConfigFile()
	configName = fmt.Sprint(configName, ".toml")
	// mysql config
	dbmap, err := db.InitDB(configName)
	if err != nil {
		logrus.Fatal(err)
	}
	app.Store.Set("mysql", dbmap)
	// redis config
	redis := redis.New(configName)
	if redis == nil {
		logrus.Fatal("redis instance is nil")
	}
	app.Store.Set("redis", dbmap)

	return nil
}

func initRoute(app *kassadin.App) error {
	app.Route.GET("/", func(ctx echo.Context) error {

		return ctx.String(200, "hello world")
	})
	app.Route.GET("/user/id/:id", func(c echo.Context) error {
		engine := app.Store.Get("mysql").(db.DBMap)["testdb"].W()
		engine.ShowSQL(true)
		user := model.User{}
		id := c.Param("id")
		has, err := engine.Where("id=?", id).Get(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return err
		}
		if !has {
			c.JSON(http.StatusBadRequest, "not fond record")
		}
		return c.JSON(http.StatusOK, user)
	})
	return nil
}
