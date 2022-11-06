package internal

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
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
	a.engine = gin.Default()
	a.engine.Use(a.ginLogger(), a.ginRecovery(a.config.Logger.Stack))
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

func (a *application) ginRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						str := strings.ToLower(se.Error())
						if strings.Contains(str, "broken pipe") ||
							strings.Contains(str, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					a.logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					a.logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					a.logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func (a *application) Run() {
	a.engine.Run(fmt.Sprintf("%s:%d", a.config.Addr, a.config.Port))
}
