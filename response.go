package gohttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ResponseWriter struct {
	w   http.ResponseWriter
	r   *http.Request
	buf strings.Builder
}

func (r *ResponseWriter) Success(data interface{}) {

	val, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-r.r.Context().Done():
		return
	default:
		r.w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		r.w.WriteHeader(http.StatusOK)
		res := fmt.Sprintf("{\"data\": %s}", val)
		r.w.Write([]byte(res))
	}

}

func (r *ResponseWriter) Error(code int, params ...string) {
	var desc string
	if len(params) == 0 {
		desc = http.StatusText(code)
	} else {
		desc = params[0]
	}
	select {
	case <-r.r.Context().Done():
		return
	default:
		r.w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		r.w.WriteHeader(code)
		res := fmt.Sprintf("{\"error\": true, \"code\": %d, \"desc\": \"%s\"}",
			code, strconv.Quote(desc))
		r.w.Write([]byte(res))
	}

}
