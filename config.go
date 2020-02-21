package gohttp

import (
	"net/http"
	"time"
)

type Config struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Handler           http.Handler
}
