package logx

type Logger interface {
	Info(msg string, kvs ...any)
	Debug(msg string, kvs ...any)
	Error(msg string, kvs ...any)
}
