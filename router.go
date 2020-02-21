package gohttp

import (
	"context"
	"net/http"
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
	_, cancel := context.WithCancel(Ctx)
	defer cancel()
	//log.Printf("=====================================\n")
	//log.Printf("[%p] Request from %q - %q '%q'\n", r, r.RemoteAddr,r.Method, r.RequestURI)
	m := self.Match(w, r)
	//log.Printf("[%p] Match -- H: '%v' - A: '%q' - P: '%v'\n", r, m.Handler, m.Action, m.Params)
	data := ControllerData{
		Action: m.Action,
		Params: m.Params,
		Store:  DbPool,
	}

	r = r.WithContext(context.WithValue(r.Context(), "data", &data))
	m.Handler(w, r)
	//log.Printf("[%p] Request completed! ctx: %p req.ctx: %p (took %s)\n", r, ctx, r.Context(), time.Since(start))
}
