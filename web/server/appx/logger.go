package appx

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"simplegame.com/simplegame/common/netx"
)

type loggerConfig struct {
	FileDir    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      int
	Stack      bool
}

func newLogger(conf loggerConfig) *zap.Logger {
	ip, err := netx.GetLocalIP()
	if err != nil {
		panic(fmt.Sprintf("init logger failed, err: %v", err))
	}
	filename := fmt.Sprintf("%s/web-%s.log", conf.FileDir, ip)
	encoder := getEncoder()
	writerSyncer := getWriterSyncer(filename, conf.MaxSize,
		conf.MaxBackups, conf.MaxAge)
	level := getLevel(conf.Level)

	core := zapcore.NewCore(encoder, writerSyncer, level)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	return logger
}

func getLevel(level int) zapcore.LevelEnabler {
	return zapcore.Level(level)
}

func getEncoder() zapcore.Encoder {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConf.TimeKey = "time"
	encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConf.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConf.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encoderConf)
}

func getWriterSyncer(fileName string, size, backups, age int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    size,
		MaxBackups: backups,
		MaxAge:     age,
	}

	return zapcore.AddSync(lumberJackLogger)
}
