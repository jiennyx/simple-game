package appx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"simplegame.com/simplegame/web/server/middleware"

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
	server http.Server
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
		app.initServer()
	})

	return app
}

func (a *application) readConfig() {
	viper.SetConfigFile("../../config/config.toml")
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

func (a *application) initServer() {
	a.server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.config.Addr, a.config.Port),
		Handler: a.engine,
	}
}

func (a *application) Run() {
	// a.engine.Run(fmt.Sprintf("%s:%d", a.config.Addr, a.config.Port))
	go func() {
		if err := a.server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("run server failed, err: %v", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	zap.L().Info("Shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		zap.L().Fatal("server shutdown failed", zap.Error(err))
	}

	zap.L().Info("server exiting...")
}
