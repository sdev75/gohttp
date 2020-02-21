package gohttp

import (
	"database/sql"
	"net/http"
)

type ControllerFunc func(http.ResponseWriter, *http.Request)

type ControllerData struct {
	Action string
	Params []string
	Store  *sql.DB
}
