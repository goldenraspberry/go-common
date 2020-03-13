package config

import (
	"github.com/dustin/go-broadcast"
	"github.com/goldenraspberry/go-common/utils"
)

type ReloadListener func()

const (
	ConfigReloadSignal = 1
)

var (
	bcService = broadcast.NewBroadcaster(10)
)

func publishReloadSignal() {
	bcService.Submit(ConfigReloadSignal)
}

func AddReloadListener(listener ReloadListener) {
	c := make(chan interface{}, 5)

	utils.Go(func() {
		for range c {
			utils.Go(listener)
		}
	})

	RegisterListenerChannel(c)
}

func RegisterListenerChannel(c chan interface{}) {
	bcService.Register(c)
}
