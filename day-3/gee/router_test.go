package gee

import "testing"

func newTestRouter() *Router {
	r := NewRouter()
	//r.addRoute("GET", "/", nil)
	//r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	//r.addRoute("GET", "/hi/:name", nil)
	//r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestGetRoute(t *testing.T) {
	newTestRouter()
}