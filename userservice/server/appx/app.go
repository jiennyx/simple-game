package appx

import (
	"sync"

	"gorm.io/gorm"
	"simplegame.com/simplegame/userservice/infra/mysql"
)

var (
	app     *application
	appOnce sync.Once
)

type application struct {
	config config
	db     *gorm.DB
}
type config struct {
	Addr string
	Port uint
}

func NewApplication() *application {
	appOnce.Do(func() {
		app = new(application)
		app.db = mysql.NewDB()
	})

	return app
}

func (a *application) injectService() {

}

func (a *application) Run() {

}
