package utils

import (
	"log"
	"runtime/debug"
)

//项目中封装的GO，带有panic的 recover
func Go(goroutine func()) {
	goWithRecover(goroutine)
}

//项目中封装的GO，带有panic的 recover
func GoWithArgs(goroutine func(args ...interface{}), args ...interface{}) {
	Go(func() { goroutine(args...) })
}

//项目中封装的GO，带有panic的 recover
func goWithRecover(goroutine func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatalf("Error in Go routine: %s, Stack : %s", err, string(debug.Stack()))
			}
		}()
		goroutine()
	}()
}
