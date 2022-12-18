package zaplog

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"simplegame.com/simplegame/common/logx"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger       *zap.SugaredLogger
	writerSyncer *lumberjack.Logger
}

func (l *logger) Info(msg string, kvs ...any) {
	l.logger.Infow(msg, kvs...)
}

func (l *logger) Debug(msg string, kvs ...any) {
	l.logger.Debugw(msg, kvs...)
}
func (l *logger) Error(msg string, kvs ...any) {
	l.logger.Errorw(msg, kvs...)
}

type Config struct {
	Level      int
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func NewZapLogger(
	level Level,
	filename string,
	options ...Option,
) (logx.Logger, func() error) {
	res := new(logger)
	res.writerSyncer = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	for _, opt := range options {
		opt(res)
	}
	encoder := getEncoder()
	core := zapcore.NewCore(encoder,
		zapcore.AddSync(res.writerSyncer),
		zapLevel[level])

	res.logger = zap.New(core, zap.AddCaller()).Sugar()

	return res, res.logger.Sync
}

type Option func(l *logger)

func MaxSize(size int) Option {
	return func(l *logger) {
		l.writerSyncer.MaxSize = size
	}
}

func MaxBackups(backups int) Option {
	return func(l *logger) {
		l.writerSyncer.MaxBackups = backups
	}
}

func MaxAge(age int) Option {
	return func(l *logger) {
		l.writerSyncer.MaxAge = age
	}
}

func Compress(compress bool) Option {
	return func(l *logger) {
		l.writerSyncer.Compress = compress
	}
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
