package mysql

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type config struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	ParseTime    bool
	Charset      string
	Timeout      string
	MaxOpenConns int
	MaxIdleConns int
}

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func NewDB() *gorm.DB {
	dbOnce.Do(func() {
		conf := readConfig()
		var err error
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%"+
			"s&parseTime=%t&loc=Local&timeout=%s",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Database,
			conf.Charset,
			conf.ParseTime,
			conf.Timeout,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("read mysql config error, err: %v", err))
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Errorf("init mysql pool error, err: %v", err))
		}
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	})

	return db
}

func readConfig() config {
	conf := config{}
	viper.SetConfigFile("../../infra/mysql/config/config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config error, err: %v", err))
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("unmarshal config error, err: %v", err))
	}

	return conf
}
