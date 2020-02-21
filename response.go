package gohttp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func WriteSuccess(ctx context.Context, w http.ResponseWriter, data interface{}) {
	if ctx.Err() != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	val, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{"))
	w.Write(val)
	w.Write([]byte("}"))
}

func WriteError(ctx context.Context, w http.ResponseWriter, statusCode int, params ...string) {
	if ctx.Err() != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var desc string
	if len(params) == 0 {
		desc = http.StatusText(statusCode)
	} else {
		desc = params[0]
	}

	w.WriteHeader(statusCode)
	w.Write([]byte("{\"error\":true,\"desc\":\""))
	w.Write([]byte(desc))
	w.Write([]byte("\"}"))
}
