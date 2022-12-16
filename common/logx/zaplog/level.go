package zaplog

import (
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	InfoLevel Level = iota
	DebugLevel
	ErrorLevel
)

var (
	zapLevel = map[Level]zapcore.Level{
		InfoLevel:  zapcore.InfoLevel,
		DebugLevel: zapcore.DebugLevel,
		ErrorLevel: zapcore.ErrorLevel,
	}
)
