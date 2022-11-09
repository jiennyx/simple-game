package errors

const (
	ServerError  = 10000
	Unauthorized = 10001
)

type GinError struct {
	Code int
	Msg  string
	Data interface{}
}

func (err *GinError) Error() string {
	return err.Msg
}

func Error(code int, msg string) *GinError {
	return &GinError{
		Code: code,
		Msg:  msg,
	}
}
