package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/goldenraspberry/go-common/config"
	"github.com/goldenraspberry/go-common/utils"
)

func RunEngine(handler http.Handler) {
	listen := config.GetListen()

	var lock sync.Mutex

	server := newServer(listen, handler)

	config.AddReloadListener(func() {
		newListen := config.GetListen()
		if newListen != listen {
			newServer := newServer(newListen, handler)
			lock.Lock()
			oldServer := server
			server = newServer
			lock.Unlock()
			utils.Go(func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				log.Println("Shutdown Server ...")
				if err := oldServer.Shutdown(ctx); err != nil {
					log.Fatalf("Server Shutdown: %v\n", err)
				}
				log.Println("Server exiting")
			})
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func newServer(listen string, handler http.Handler) *http.Server {
	var srv *http.Server
	srv = &http.Server{
		Addr:    listen,
		Handler: handler,
	}

	return srv
}
