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

func Init(handler http.Handler, cfg Config, pool *sql.DB) *http.Server {
	Ctx, CtxCancel = context.WithCancel(context.Background())
	Srv = &http.Server{
		Addr:              cfg.Addr,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		Handler:           handler,
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	Srv.Shutdown(ctx)
	CtxCancel()
	if DbPool != nil {
		DbPool.Close()
	}
}
