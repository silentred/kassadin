package main

import "github.com/silentred/kassadin"

func main() {
	app := kassadin.NewApp()
	app.RegisterConfigHook(initConfig)
	app.Start()
}

func initConfig(app *kassadin.App) error {
	return nil
}
