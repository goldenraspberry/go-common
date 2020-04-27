package log

import (
	"context"
	"fmt"
)

type Logger interface{}

type log interface{}

func Logf(level string, category string, ctx context.Context, format string, v ...interface{}) {
	Log(level, category, ctx, fmt.Sprintf(format, v...))
}

func Log(level string,category string, ctx context.Context, format string) {
	// dolog
}
