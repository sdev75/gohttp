package gohttp

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Init(cfg Config, pool *sql.DB) *http.Server {
	Ctx, CtxCancel = context.WithCancel(context.Background())
	Srv = &http.Server{
		Addr:              cfg.Addr,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 30,
		Handler:           cfg.Handler,
	}

	DbPool = pool

	appSignal := make(chan os.Signal, 1)
	signal.Notify(appSignal,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT)

	go func() {
		select {
		case sig := <-appSignal:
			log.Printf("Received signal: %s\n", sig.String())
			CtxCancel()
		}
	}()

	return Srv
}

func Start() {
	if err := Srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func Stop() {
	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	Srv.Shutdown(ctx)

	log.Println("Cleaning up resources...")
	CtxCancel()

	log.Println("Closing database pool...")
	DbPool.Close()

	log.Println("Goodbye!")
}
