package logger

import (
	"context"
	"fmt"
	"runtime"
)

// ErrorKV логгирует ошибки с уровнем error
func ErrorKV(ctx context.Context, data Data) {
	FromContext(ctx).Errorw(data.Msg, data.toParameters()...)
}

// InfoKV логгирует ошибки с уровнем info
func InfoKV(ctx context.Context, data Data) {
	FromContext(ctx).Infow(data.Msg, data.toParameters()...)
}

// WarnKV логгирует ошибки с уровнем warn
func WarnKV(ctx context.Context, data Data) {
	FromContext(ctx).Warnw(data.Msg, data.toParameters()...)
}

// FatalKV логгирует ошибки с уровнем fatal
func FatalKV(ctx context.Context, data Data) {
	FromContext(ctx).Fatalw(data.Msg, data.toParameters()...)
}

// Data данные для обертки над логгированием.
type Data struct {
	Msg    string
	Error  error
	Panic  any
	Detail any
}

func (d *Data) toParameters() []any {
	var ans []any

	if d.Error != nil {
		ans = append(ans, "error", d.Error)
	}

	if d.Panic != nil {
		ans = append(ans, "panic", d.Panic)
	}

	if d.Detail != nil {
		ans = append(ans, "detail", fmt.Sprintf("%+v", d.Detail))
	}

	_, filename, lineNum, _ := runtime.Caller(2)
	ans = append(ans, "file_name", fmt.Sprintf("%+v", filename))
	ans = append(ans, "line_num", fmt.Sprintf("%+v", lineNum))

	return ans
}
