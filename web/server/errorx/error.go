package errorx

import "google.golang.org/grpc/codes"

const (
	OK           = 0
	ServerError  = 10000
	Unauthorized = 10001
	ParamError   = 10002
)

var (
	statusMsg = map[int]string{
		OK:           "ok",
		ServerError:  "server error",
		Unauthorized: "unauthorized",
		ParamError:   "param error",
	}
)

func GetCodeMsg(code int) string {
	return statusMsg[code]
}

type GinError struct {
	Code int
	Msg  string
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

func ErrorFromStatus(code codes.Code, msg string) *GinError {
	return Error(int(code), msg)
}
