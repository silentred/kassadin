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

	"encoding/json"

	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/silentred/kassadin/util"
	"github.com/silentred/kassadin/util/container"
	"github.com/silentred/kassadin/util/rotator"
	"github.com/silentred/kassadin/util/strings"
	"github.com/spf13/viper"
)

var (
	// AppMode is App's running envirenment. Valid values are dev and prod
	AppMode    string
	ConfigFile string
	LogPath    string
)

func init() {
	flag.StringVar(&AppMode, "mode", "", "RunMode of the application: dev or prod")
	flag.StringVar(&ConfigFile, "cfg", "", "absolute path of config file")
	flag.StringVar(&LogPath, "logPath", ".", "logPath is where log file will be")
}

// HookFunc when app starting and tearing down
type HookFunc func(*App) error

// App represents the application
type App struct {
	Store    *container.Map
	Injector container.Injector
	Route    *echo.Echo

	loggers map[string]*logrus.Logger
	Config  AppConfig

	configHook   HookFunc
	loggerHook   HookFunc
	serviceHook  HookFunc
	routeHook    HookFunc
	shutdownHook HookFunc
}

// NewApp gets a new application
func NewApp() *App {
	app := &App{
		Store:    &container.Map{},
		Injector: container.NewInjector(),
		Route:    echo.New(),
		loggers:  make(map[string]*logrus.Logger),
	}
	// register App itself
	app.Set("app", app, nil)
	return app
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

// DefaultLogger gets default logger
func (app *App) DefaultLogger() *logrus.Logger {
	return app.Logger("")
}

// Set object into app.Store and Map it into app.Injector
func (app *App) Set(key string, object interface{}, ifacePtr interface{}) {
	app.Store.Set(key, object)
	if ifacePtr != nil {
		app.Injector.MapTo(object, ifacePtr)
	} else {
		app.Injector.Map(object)
	}
}

// Get object from app.Store
func (app *App) Get(key string) interface{} {
	return app.Store.Get(key)
}

// Inject dependencies to the object. Please MAKE SURE that the dependencies should be stored at app.Injector
// before this method is called. Please use app.Set() to make this happen.
func (app *App) Inject(object interface{}) error {
	return app.Injector.Apply(object)
}

// InitConfig in format of toml
func (app *App) InitConfig() {
	// use viper to resolve config.toml
	if ConfigFile == "" {
		var configName = app.getConfigFile()
		viper.AddConfigPath(".")
		viper.AddConfigPath(util.SelfDir())
		viper.SetConfigName(configName)
	} else {
		viper.SetConfigFile(ConfigFile)
	}

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
	config.Log = l

	// TODO: session config
	// mysql config
	mysql := MysqlConfig{}
	mysqlConfig := viper.Get("mysql")
	mysqlConfigBytes, err := json.Marshal(mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(mysqlConfigBytes, &mysql.Instances)
	if err != nil {
		log.Fatal(err)
	}
	config.Mysql = mysql

	// redis config
	redis := RedisInstance{}
	redis.Host = viper.GetString("redis.host")
	redis.Port = viper.GetInt("redis.port")
	redis.Db = viper.GetInt("redis.db")
	redis.Pwd = viper.GetString("redis.password")
	config.Redis = redis

	app.Config = config

	// hook
	if app.configHook != nil {
		err = app.configHook(app)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (app *App) getConfigFile() string {
	var configName = "config"
	if AppMode != "" {
		configName = fmt.Sprintf("%s.%s", "config", AppMode)
	}
	return configName
}

func (app *App) InitLogger() {
	// new default Logger
	var writer io.Writer
	var spliter rotator.Spliter
	var err error

	logConfig := app.Config.Log

	switch logConfig.Providor {
	case ProvidorFile:
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

			writer = rotator.NewFileRotator(logConfig.LogPath, app.Config.Name, "log", spliter)
		} else {
			writer, err = os.Open(filepath.Join(logConfig.LogPath, app.Config.Name+".log"))
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		writer = os.Stdout
	}

	defaultLogger := logrus.New()
	defaultLogger.Formatter = &logrus.JSONFormatter{}
	defaultLogger.Out = writer
	switch app.Config.Mode {
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

func (app *App) Init() {
	app.InitConfig()
	app.InitLogger()
	app.initService()
	app.initRoute()
}

// Start running the application
func (app *App) Start() {
	app.Init()
	//app.route.Start(fmt.Sprintf(":%d", app.config.Port))
	app.graceStart()
	app.shutdown()
}

func (app *App) graceStart() error {
	// Start server
	go func() {
		if err := app.Route.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
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
