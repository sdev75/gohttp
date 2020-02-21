package gohttp

import (
	"context"
	"database/sql"
	"net/http"
	"os"
)

var (
	Ctx       context.Context
	CtxCancel context.CancelFunc
	Srv       *http.Server
	DbPool    *sql.DB
	IntSignal chan os.Signal
)
