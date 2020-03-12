package config

import "sync"

type ReloadListener func()

var (
	listListener []ReloadListener
	listenerLock sync.Mutex
)

func publishReloadSignal() {
	listenerLock.Lock()
	listenerLock.Unlock()
	if len(listListener) == 0 {
		return
	}
	for _, listener := range listListener {
		listener()
	}
}

func AddReloadListener(listener ReloadListener) {
	listenerLock.Lock()
	listListener = append(listListener, listener)
	listenerLock.Unlock()
}
