package gohttp

import (
	"database/sql"
	"net/http"
)

type ControllerData struct {
	Action string
	Params []string
	Store  *sql.DB
}

var controllers map[string]http.HandlerFunc

func init() {
	controllers = make(map[string]http.HandlerFunc)
}

func RegisterController(name string, h http.HandlerFunc) {
	controllers[name] = h
}

func GetController(name string) http.HandlerFunc {
	return controllers[name]
}
