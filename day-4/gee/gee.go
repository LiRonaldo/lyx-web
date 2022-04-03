package gee

import "net/http"

type (
	RouterGroup struct {
		engine      *Engine
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
	}

	Engine struct {
		*RouterGroup
		router *Router
		groups []*RouterGroup
	}
)

func NewGee() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.AddRouter(pattern, "GET", handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRouter("POST", pattern, handler)
}

func (group *RouterGroup) AddRouter(comp, method string, hander HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRouter(pattern, method, hander)
}

func (e *Engine) RUN(port string) error {
	return http.ListenAndServe(port, e)
}
func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, resp)
	e.router.Handler(ctx)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}
