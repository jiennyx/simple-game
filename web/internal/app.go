package internal

import (
	"fmt"
	"simplegame.com/simplegame/web/internal/middleware"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	app     *application
	appOnce sync.Once
)

type application struct {
	config config
	engine *gin.Engine
	logger *zap.Logger
}

type config struct {
	Addr    string
	Port    uint
	Proxies []string
	Logger  loggerConfig
}

func NewApplication() *application {
	appOnce.Do(func() {
		app = new(application)
		app.readConfig()
		app.initLogger()
		app.initEngine()
	})

	return app
}

func (a *application) readConfig() {
	viper.SetConfigFile("../internal/conf/config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config error, err: %v", err))
	}
	if err := viper.Unmarshal(&a.config); err != nil {
		panic(fmt.Errorf("unmarshal config error, err: %v", err))
	}

	fmt.Printf("init config succeed, log will print to: %s\n",
		a.config.Logger.FileName)
}

func (a *application) initLogger() {
	a.logger = newLogger(a.config.Logger)
}

func (a *application) initEngine() {
	a.engine = gin.New()
	a.engine.SetTrustedProxies(a.config.Proxies)
	a.engine.Use(middleware.LoggerHandler(),
		middleware.RecoveryHandler(a.config.Logger.Stack))
	a.engine.Use(middleware.ErrorHandler())
	a.engine.Use(middleware.JWTAuthHandler())
	initRoute(a.engine)
}

func (a *application) Run() {
	a.engine.Run(fmt.Sprintf("%s:%d", a.config.Addr, a.config.Port))
}
