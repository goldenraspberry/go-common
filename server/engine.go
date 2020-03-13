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

	var (
		server *http.Server
		// block engine
		wg   = &sync.WaitGroup{}
		lock = &sync.Mutex{}
	)

	server = newServer(listen, handler)
	wg.Add(1)

	config.AddReloadListener(func() {
		newListen := config.GetListen()
		if newListen != listen {
			newServer := newServer(newListen, handler)
			wg.Add(1)

			// add lock
			lock.Lock()

			// auto shutdown
			utils.GoWithArgs(func(args ...interface{}) {
				server := args[0].(*http.Server)
				listen := args[1].(string)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				log.Println("Shutdown listen[" + listen + "] Server ...")
				if err := server.Shutdown(ctx); err != nil {
					log.Fatalf("Server Shutdown: %v\n", err)
				}
				wg.Done()
				log.Println("Server listen[" + listen + "] exiting")
			}, server, listen)

			// replace and start new server
			server = newServer
			listen = newListen

			// free lock
			lock.Unlock()

			log.Println("Server listen[" + newListen + "] starting... ")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				wg.Done()
				log.Fatalf("listen: %s\n", err)
			}
		}
	})

	log.Println("Server listen[" + listen + "] starting... ")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		wg.Done()
		log.Fatalf("listen: %s\n", err)
	}
	wg.Wait()
}

func newServer(listen string, handler http.Handler) *http.Server {
	var srv *http.Server
	srv = &http.Server{
		Addr:    listen,
		Handler: handler,
	}

	return srv
}
