package config

import "sync"

var (
	lock sync.RWMutex
)

func init() {
	initEnv()
	initConfig()
}
