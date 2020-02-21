package gohttp

import (
	"database/sql"
	"net/http"
)

const version = 100

type ControllerFunc func(http.ResponseWriter, *http.Request)

type ControllerData struct {
	Action string
	Params []string
	Writer *ResponseWriter
	Store  *sql.DB
}
