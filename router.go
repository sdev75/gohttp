package gohttp

import (
	"context"
	"net/http"
	"time"
)

type SimpleRouterMatch struct {
	Handler func(http.ResponseWriter, *http.Request)
	Action  string
	Params  []string
}

type SimpleRouter struct {
	Match func(http.ResponseWriter, *http.Request) *SimpleRouterMatch
}

func NewSimpleRouter() *SimpleRouter {
	return new(SimpleRouter)
}

func (self *SimpleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//start := time.Now()
	//log.Printf("=====================================\n")
	//log.Printf("[%p] Request from %q - %q '%q'\n", r, r.RemoteAddr,r.Method, r.RequestURI)

	ctx, cancel := context.WithCancel(Ctx)
	defer cancel()

	m := self.Match(w, r)
	if m == nil {
		cancel()
		return
	}

	//log.Printf("[%p] Match -- H: '%v' - A: '%q' - P: '%v'\n", r, m.Handler, m.Action, m.Params)

	data := ControllerData{
		Action: m.Action,
		Params: m.Params,
		Store:  DbPool,
		Writer: &ResponseWriter{
			w: w,
			r: r,
		},
	}
	r = r.WithContext(context.WithValue(r.Context(), "data", &data))

	go func() {
		//log.Printf("[%p] go: Invoking handler() with ctx: %p / req_ctx: %p...\n",	r, ctx, r.Context())

		go func() {
			defer cancel()
			m.Handler(w, r)
		}()

		select {
		case <-time.After(time.Second * 30):
			//log.Printf("[%p] [go] context timeout: %p\n", r, ctx)
			cancel()
		case <-r.Context().Done():
			//log.Printf("[%p] [go] req context cancelled: %p\n", r, r.Context())
			cancel()
		case <-ctx.Done():
			//log.Printf("[%p] [go] context cancelled: %p\n", r, ctx)
			cancel()
		}

		//log.Printf("[%p] [go] done: %p\n", r, ctx)
	}()

	select {
	case <-r.Context().Done():
		//log.Printf("[%p] req context cancelled: %p\n", r, r.Context())
		cancel()
	case <-ctx.Done():
		//log.Printf("[%p] ctx cancelled: %p\n", r, ctx)
	}

	//log.Printf("[%p] Request completed! ctx: %p req.ctx: %p (took %s)\n", r, ctx, r.Context(), time.Since(start))
}
