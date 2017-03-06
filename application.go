package kassadin

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/silentred/kassadin/util"
	"github.com/silentred/kassadin/util/container"
	"github.com/silentred/kassadin/util/rotator"
	"github.com/silentred/kassadin/util/strings"
	"github.com/spf13/viper"
	"github.com/silentred/kassadin/db"
	"github.com/silentred/kassadin/redis"
	"github.com/golang/go/src/pkg/path"
)

var (
	AppMode string
)

func init() {
	flag.StringVar(&AppMode, "mode", "dev", "RunMode of the application: dev or prod")
}

// Map stores objects
type Map map[string]interface{}

// HookFunc when app starting and tearing down
type HookFunc func(*App) error

// App represents the application
type App struct {
	Store        *Map
	Injector     container.Injector
	Route        *echo.Echo
	loggers      map[string]*logrus.Logger
	config       AppConfig
	configHook   HookFunc
	loggerHook   HookFunc
	serviceHook  HookFunc
	routeHook    HookFunc
	shutdownHook HookFunc
}

// NewApp gets a new application
func NewApp() *App {
	return &App{
		Store:    &Map{},
		Injector: container.NewInjector(),
		Route:    echo.New(),
		loggers:  make(map[string]*logrus.Logger),
	}
}

// Logger of name
func (app *App) Logger(name string) *logrus.Logger {
	if name == "" {
		return app.loggers["default"]
	}
	if l, ok := app.loggers[name]; ok {
		return l
	}

	return nil
}

// InitConfig in format of toml
func (app *App) initConfig() {
	// use viper to resolve config.toml
	var configName = "config"
	if AppMode != "" {
		configName = fmt.Sprintf("%s.%s", "config", AppMode)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath(util.SelfDir())
	viper.SetConfigName(configName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	// make AppConfig; set data from viper
	config := AppConfig{}
	config.Name = viper.GetString("app.name")
	config.Mode = viper.GetString("app.runMode")
	config.Port = viper.GetInt("app.port")

	// log config
	l := LogConfig{}
	l.Name = "default"
	l.LogPath = viper.GetString("app.logPath")
	l.Providor = viper.GetString("app.logProvider")
	l.RotateEnable = viper.GetBool("app.logRotate")
	l.RotateMode = viper.GetString("app.logRotateType")
	l.RotateLimit = viper.GetString("app.logLimit")

	// TODO: session config
	// TODO: mysql config

	currpath, _ := os.Getwd()
	confpath := path.Join(currpath, AppMode, ".toml")

	dbmap, err := db.InitDB(confpath)
	if err != nil {
		log.Fatal(err)
	}
	app.Store["mysql"] = dbmap
	// TODO: redis config
	redis := redis.New(confpath)
	if redis == nil {
		log.Fatal("redis instance is nil")
	}
	app.Store["redis"] = dbmap

	config.Log = l
	app.config = config

	// hook
	if app.configHook != nil {
		err = app.configHook(app)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) initLogger() {
	// new default Logger
	var writer io.Writer
	var spliter rotator.Spliter

	logConfig := app.config.Log

	if logConfig.RotateEnable {
		switch logConfig.RotateMode {
		case RotateByDay:
			spliter = rotator.NewDaySpliter()
		case RotateBySize:
			limitSize, err := strings.ParseByteSize(logConfig.RotateLimit) // 100 MB
			if err != nil {
				log.Fatal(err)
			}
			spliter = rotator.NewSizeSpliter(uint64(limitSize))
		default:
			log.Fatalf("invalid RotateMode: %s", logConfig.RotateMode)
		}

		writer = rotator.NewFileRotator(logConfig.LogPath, app.config.Name, "log", spliter)
	}

	if writer == nil {
		writer = os.Stdout
	}

	defaultLogger := logrus.New()
	defaultLogger.Formatter = &logrus.JSONFormatter{}
	defaultLogger.Out = writer
	switch app.config.Mode {
	case ModeDev:
		defaultLogger.Level = logrus.DebugLevel
	case ModeProd:
		defaultLogger.Level = logrus.ErrorLevel
	default:
		defaultLogger.Level = logrus.DebugLevel
	}

	// set logger
	app.loggers["default"] = defaultLogger

	// hook
	if app.loggerHook != nil {
		err := app.loggerHook(app)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) initService() {
	// hoook
	if app.serviceHook != nil {
		err := app.serviceHook(app)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) initRoute() {
	// hook
	if app.routeHook != nil {
		err := app.routeHook(app)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) shutdown() {
	// hook
	if app.shutdownHook != nil {
		err := app.shutdownHook(app)
		if err != nil {
			log.Fatal(err)
		}
	}
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
		if err := app.Route.Start(fmt.Sprintf(":%d", app.config.Port)); err != nil {
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
	if err := app.Route.Shutdown(ctx); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
