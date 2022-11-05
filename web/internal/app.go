package internal

import (
	"fmt"
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
	engine *gin.Engine
	config config
	logger *zap.Logger
}

type config struct {
	Addr    string
	Port    uint
	LogFile string
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

	fmt.Printf("init config succeed, log will print to: %s", conf.LogFile)
}

func (a *application) initLogger() {}

func (a *application) Run() {
	a.engine.Run(fmt.Sprintf("%s:%d", a.config.Addr, a.config.Port))
}
