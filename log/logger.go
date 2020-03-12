package log

import (
	"context"
	"fmt"
)

type Logger interface{}

type log interface{}

func Logf(level string, ctx context.Context, format string, v ...interface{}) {
	Log(level, ctx, fmt.Sprintf(format, v...))
}

func Log(level string, ctx context.Context, format string) {

}
