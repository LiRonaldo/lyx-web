package gee

import "net/http"

type HandlerFunc func(ctx *Context)

type Gee struct {
	router *Router
}

func NewGee() *Gee {
	return &Gee{
		router: NewRouter(),
	}
}
func (g *Gee) Get(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}
func (g *Gee) Post(pattern string, handler HandlerFunc) {
	g.router.addRoute("POST", pattern, handler)
}
func (g *Gee) RUN(port string) {
	http.ListenAndServe(port, g)
}
func (g *Gee) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, writer)
	g.router.Handler(ctx)
}

func (g *Gee) addRoute(method string, pattern string, handler HandlerFunc) {
	g.router.addRoute(method, pattern, handler)
}
