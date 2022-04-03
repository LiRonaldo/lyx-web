package gee

import (
	"net/http"
	"strings"
)

type (
	RouterGroup struct {
		engine      *Engine
		prefix      string
		parent      *RouterGroup
		middlewares []HandlerFunc
	}
	Engine struct {
		*RouterGroup
		groups []*RouterGroup
		Router *Router
	}
)

func New() *Engine {
	engine := &Engine{}
	engine.Router=newRouter()
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(cmp string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + cmp,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) Get(pattern string, handler HandlerFunc) {
	group.AddRouter(pattern, "GET", handler)
}

func (group *RouterGroup) AddRouter(pattern, method string, handler HandlerFunc) {
	path := group.prefix + pattern
	group.engine.Router.addRouter(path, method, handler)
}

func (e *Engine) Run(prot string) {
	http.ListenAndServe(prot, e)
}
func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx := newContext(req, resp)
	ctx.handlers = middlewares
	e.Router.Handler(ctx)

}
