package kassadin

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/silentred/kassadin/util"
	"github.com/silentred/kassadin/util/container"
	"github.com/spf13/viper"
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
	config  AppConfig

	configHook   HookFunc
	loggerHook   HookFunc
	serviceHook  HookFunc
	routeHook    HookFunc
	shutdownHook HookFunc
}

// NewApp gets a new application
func NewApp() *App {
	return nil
}

// NewLogger in App.loggers
func (app *App) NewLogger(config LogConfig) {

}

// Logger of name
func (app *App) Logger(name string) *logrus.Logger {
	return nil
}

// InitConfig in format of toml
func (app *App) initConfig() {
	// viper resolve config.toml
	viper.AddConfigPath(".")
	viper.AddConfigPath(util.SelfDir())
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

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

func (app *App) shutdown() {
	// hook
}

// RegisterConfigHook at initConfig
func (app *App) RegisterConfigHook(hook HookFunc) {
	app.configHook = hook
}

func (app *App) RegisterLoggerHook(hook HookFunc) {
	app.loggerHook = hook
}

func (app *App) RegisterServiceHook(hook HookFunc) {
	app.serviceHook = hook
}

func (app *App) RegisterRouteHook(hook HookFunc) {
	app.routeHook = hook
}

func (app *App) RegisterShutdownHook(hook HookFunc) {
	app.shutdownHook = hook
}

// Start running the application
func (app *App) Start() {
	app.initConfig()
	app.initLogger()
	app.initService()
	app.initRoute()

	//app.route.Start(fmt.Sprintf(":%d", app.config.Port))
	app.graceStart()

	app.shutdown()
}

func (app *App) graceStart() error {
	// Start server
	go func() {
		if err := app.route.Start(fmt.Sprintf(":%d", app.config.Port)); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.route.Shutdown(ctx); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
