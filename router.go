package gohttp

import (
	"context"
	"net/http"
	"time"
)

type Router struct {
	MatchedRoute MatchedRoute
	MatchHandler func(*http.Request)
}

func NewRouter() *Router {
	return new(Router)
}

type MatchedRoute struct {
	Controller func(http.ResponseWriter, *http.Request)
	Action     string
	Params     []string
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//start := time.Now()
	//log.Printf("=====================================\n")
	//log.Printf("[%p] Request from %q - %q '%q'\n", req, req.RemoteAddr,
	//	req.Method, req.RequestURI)

	r.MatchHandler(req)
	//log.Printf("[%p] Matched Route -- Controller: '%v' - Action: '%q' - Params: '%v'\n",
	//	r.MatchedRoute.Controller, r.MatchedRoute.Action, r.MatchedRoute.Params)

	ctx, cancel := context.WithCancel(Ctx)
	data := ControllerData{
		Action: r.MatchedRoute.Action,
		Params: r.MatchedRoute.Params,
		Store:  DbPool,
		Writer: &ResponseWriter{
			w: w,
			r: req,
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), "data", &data))

	go func() {
		//log.Printf("[%p] go: Invoking handler() with ctx: %p / req_ctx: %p...\n",
		//	req, ctx, req.Context())

		go func() {
			defer cancel()
			r.MatchedRoute.Controller(w, req)
		}()

		select {
		case <-time.After(time.Second * 30):
			//log.Printf("[%p] [go] context timeout: %p\n", req, ctx)
			cancel()
		case <-req.Context().Done():
			//log.Printf("[%p] [go] request context cancelled: %p\n", req, req.Context())
			cancel()
		case <-ctx.Done():
			//log.Printf("[%p] [go] context cancelled: %p\n", req, ctx)
			cancel()
		}

		//log.Printf("[%p] [go] done: %p\n", req, ctx)
	}()

	select {
	case <-req.Context().Done():
		//log.Printf("[%p] req context cancelled: %p\n", req, req.Context())
		cancel()
	case <-ctx.Done():
		//log.Printf("[%p] ctx cancelled: %p\n", req, ctx)
	}

	//log.Printf("[%p] Request completed! ctx: %p req.ctx: %p (took %s)\n", req, ctx, req.Context(), time.Since(start))
}
