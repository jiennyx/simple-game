package appx

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"simplegame.com/simplegame/common/api/user"
	"simplegame.com/simplegame/common/clients"
	"simplegame.com/simplegame/common/logx"
	"simplegame.com/simplegame/common/logx/zaplog"
	"simplegame.com/simplegame/common/netx"
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
	config   config
	db       *gorm.DB
	server   *grpc.Server
	pool     map[string]map[string]*grpc.ClientConn
	logger   logx.Logger
	logFlush func() error
}

type config struct {
	Network     string
	IP          string
	Port        int
	ServiceName string
	AllServices []string
	Etcd        clients.EtcdConfig
	Logger      zaplog.Config
}

func NewApplication() *application {
	appOnce.Do(func() {
		app = new(application)
		app.initConfig()
		app.initLogger()
		app.initDB()
		app.registerService()
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

	ip, err := netx.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("get local ip error, err: %v", err))
	}
	a.config.IP = ip
}

func (a *application) initDB() {
	a.db = mysql.NewDB()
}

func (a *application) initLogger() {
	a.config.Logger.Filename = fmt.Sprintf(
		a.config.Logger.Filename,
		a.config.ServiceName,
		a.config.IP,
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

func (a *application) registerService() {
	fmt.Println(a.config)
	err := clients.RegisterService(
		a.config.ServiceName,
		a.config.IP,
		a.config.Port,
		a.config.Etcd,
	)
	if err != nil {
		panic(fmt.Sprintf("register service to etcd failed, err: %v", err))
	}
}

func (a *application) cancelService() {
	err := clients.CancelService(
		a.config.ServiceName,
		a.config.IP,
		a.config.Port,
		a.config.Etcd,
	)
	if err != nil {
		panic(fmt.Sprintf("cancel service to etcd failed, err: %v", err))
	}
}

func (a *application) initServer() {
	userRepo := dao.NewUserRepository(a.db)
	registerDomainService, err := domainService.NewRegisterDomainService(
		a.logger,
		domainService.WithUserRepository(userRepo),
	)
	if err != nil {
		panic(fmt.Sprintf("create register damain service failed, err: %v", err))
	}
	registerApplicationService, err :=
		applicationService.NewRegisterApplicationService(
			a.logger,
			applicationService.WithRegisterDomainService(registerDomainService),
			applicationService.WithUserRepository(userRepo),
		)
	if err != nil {
		panic(fmt.Sprintf("create register application service failed, err: %v", err))
	}
	userServer, err := facade.NewUserServer(
		a.logger,
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

func (a *application) WaitShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT,
		syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigs:
		a.cancelService()
		a.logFlush()
		a.server.GracefulStop()
	}
}
