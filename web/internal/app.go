package internal

import (
	"fmt"
	"sync"
	"time"

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
	Addr   string
	Port   uint
	Logger loggerConfig
}

func NewApplication() *application {
	appOnce.Do(func() {
		app = new(application)
		app.readConfig()
		app.initLogger()
	})

	return app
}

func (a *application) readConfig() {
	var conf config
	viper.SetConfigName("conf")
	viper.AddConfigPath("./internal/conf")
	if err := viper.ReadInConfig(); err != nil {
		panic("read config error")
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic("unmarshal config error")
	}

	fmt.Printf("init config succeed, log will print to: %s",
		conf.Logger.FileName)
}

func (a *application) initLogger() {
	a.logger = newLogger(a.config.Logger)
}

func (a *application) initEngine() {
	a.engine = gin.Default()
}

func (a *application) ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		a.logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

func (a *application) ginRecover(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
	}
}

func (a *application) Run() {
	a.engine.Run(fmt.Sprintf("%s:%d", a.config.Addr, a.config.Port))
}
