package appx

import (
	"fmt"
	"net"
	"sync"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"simplegame.com/simplegame/common/api/user"
	applicationService "simplegame.com/simplegame/userservice/application/service"
	domainService "simplegame.com/simplegame/userservice/domain/service"
	"simplegame.com/simplegame/userservice/infra/mysql"
	"simplegame.com/simplegame/userservice/infra/mysql/dao"
	"simplegame.com/simplegame/userservice/interfaces/facade"
)

var (
	app     *application
	appOnce sync.Once
)

type application struct {
	config config
	db     *gorm.DB
	server *grpc.Server
}
type config struct {
	Network string
	Port    uint
}

func NewApplication() *application {
	appOnce.Do(func() {
		app = new(application)
		app.initConfig()
		app.initDB()
		app.initServer()
	})

	return app
}

func (a *application) initConfig() {
	viper.SetConfigFile("../../config/config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config error, err: %v", err))
	}
	if err := viper.Unmarshal(&a.config); err != nil {
		panic(fmt.Errorf("unmarshal config error, err: %v", err))
	}
}

func (a *application) initDB() {
	a.db = mysql.NewDB()
}

func (a *application) initServer() {
	userRepo := dao.NewUserRepository(a.db)
	registerDomainService, err := domainService.NewRegisterDomainService(
		domainService.WithUserRepository(userRepo),
	)
	if err != nil {
		panic(fmt.Sprintf("create register damain service failed, err: %v", err))
	}
	registerApplicationService, err :=
		applicationService.NewRegisterApplicationService(
			applicationService.WithRegisterDomainService(registerDomainService),
			applicationService.WithUserRepository(userRepo),
		)
	if err != nil {
		panic(fmt.Sprintf("create register application service failed, err: %v", err))
	}
	userServer, err := facade.NewUserServer(
		facade.WithRegisterApplicationService(registerApplicationService),
	)
	if err != nil {
		panic(fmt.Sprintf("create user server failed, err: %v", err))
	}

	a.server = grpc.NewServer()
	user.RegisterUserServer(a.server, &userServer)
}

func (a *application) Run() {
	listen, err := net.Listen(a.config.Network, fmt.Sprintf(":%d", a.config.Port))
	if err != nil {
		panic(fmt.Sprintf("run userservice server failed, err: %v", err))
	}
	a.server.Serve(listen)
}
