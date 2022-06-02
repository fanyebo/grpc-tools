package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Int32ToCode 将错误码转为codes.Code类型
func Int32ToCode(num int32) codes.Code {
	return codes.Code(uint32(num))
}

// CodeToError 将错误码转换为error
func CodeToError(code codes.Code, msgs ...string) error {
	msg := fmt.Sprintf("%d", code)
	if len(msgs) > 0 {
		msg = msgs[0]
	}

	st := status.New(code, msg)
	return st.Err()
}

// ErrorToCode 将错误转为错误码
func ErrorToCode(err error) codes.Code {
	st := status.Convert(err)
	return st.Code()
}
