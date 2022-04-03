package gee

import (
	"fmt"
	"net/http"
)

type HandelrFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandelrFunc
}

func NewEngine() *Engine {
	return &Engine{router: make(map[string]HandelrFunc)}
}
func (e *Engine) AddRouter(method, pattern string, handler HandelrFunc) {
	path := method + "-" + pattern
	e.router[path] = handler
}

func (e *Engine) Get(pattern string, handler HandelrFunc) {
	e.AddRouter("GET", pattern, handler)
}
func (e *Engine) Post(pattern string, handler HandelrFunc) {
	e.AddRouter("POST", pattern, handler)
}
func (e *Engine) RUN(port string) {
	http.ListenAndServe(port, e)
}
func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	path := req.Method + "-" + req.URL.Path
	if v, ok := e.router[path]; ok {
		v(resp, req)
	} else {
		fmt.Fprintf(resp, "404 not found %s", path)
	}
}
