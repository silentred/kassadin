package kassadin

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/silentred/kassadin/util/container"
)

// Map stores objects
type Map map[string]interface{}

// HookFunc when app starting and tearing down
type HookFunc func(*App) error

// App represents the application
type App struct {
	Store    *Map
	Injector container.Injector

	route   *echo.Echo
	loggers map[string]*logrus.Logger
	config  appConfig

	configHook  HookFunc
	loggerHook  HookFunc
	serviceHook HookFunc
	routeHook   HookFunc
}

// NewApp gets a new application
func NewApp() *App {
	return nil
}

// NewLogger in App.loggers
func (app *App) NewLogger() {

}

// Logger of name
func (app *App) Logger(name string) *logrus.Logger {

}

// InitConfig in format of toml
func (app *App) initConfig() {
	// viper resolve config.toml
	// hook
}

func (app *App) initLogger() {
	// new default Logger
	// hook
}

func (app *App) initService() {
	// hoook
}

func (app *App) initRoute() {
	// hook
}

// Start running the application
func (app *App) Start() {
	app.initConfig()
	app.initLogger()
	app.initService()
	app.initRoute()

	app.route.Start(fmt.Sprintf(":%d", app.config.Port))
}
