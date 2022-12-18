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

	"simplegame.com/simplegame/common/clients"
	"simplegame.com/simplegame/common/logx"
	"simplegame.com/simplegame/common/logx/zaplog"
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
	config   config
	engine   *gin.Engine
	server   http.Server
	logger   logx.Logger
	logFlush func() error
}

type config struct {
	Web    webConfig
	Etcd   clients.EtcdConfig
	Logger zaplog.Config
}

type webConfig struct {
	Addr        string
	ServiceName string
	IP          string
	Port        uint
	Proxies     []string
	Services    []string
	Stack       bool
}

func NewApplication() *application {
	appOnce.Do(func() {
		app = new(application)
		app.readConfig()
		app.initLogger()
		app.initEngine()
		app.initServicePool()
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
}

func (a *application) initLogger() {
	a.config.Logger.Filename = fmt.Sprintf(
		a.config.Logger.Filename,
		a.config.Web.ServiceName,
		a.config.Web.IP,
	)
	a.logger, a.logFlush = zaplog.NewZapLogger(
		zaplog.Level(a.config.Logger.Level),
		a.config.Logger.Filename,
		zaplog.MaxSize(a.config.Logger.MaxSize),
		zaplog.MaxAge(a.config.Logger.MaxAge),
		zaplog.MaxBackups(a.config.Logger.MaxBackups),
		zaplog.Compress(a.config.Logger.Compress),
	)
}

func (a *application) initEngine() {
	a.engine = gin.New()
	a.engine.SetTrustedProxies(a.config.Web.Proxies)
	a.engine.Use(middleware.LoggerHandler(),
		middleware.RecoveryHandler(a.config.Web.Stack))
	a.engine.Use(middleware.ErrorHandler())
	a.engine.Use(middleware.JWTAuthHandler())
	initRoute(a.engine)
}

func (a *application) initServicePool() {
	err := clients.DiscoverService(a.config.Web.Services, a.config.Etcd)
	if err != nil {
		panic(fmt.Sprintf("init service pool failed, err: %v", err))
	}
}

func (a *application) initServer() {
	a.server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.config.Web.Addr, a.config.Web.Port),
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
