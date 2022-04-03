package gee

import (
	"html/template"
	"net/http"
	"path"
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
		groups        []*RouterGroup
		Router        *Router
		htmlTemplates *template.Template
		funcMap       template.FuncMap
	}
)

func New() *Engine {
	engine := &Engine{}
	engine.Router = newRouter()
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

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.Get(urlPattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.setStatus(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Resp, c.Req)
	}
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
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
	ctx.engine = e
	ctx.handlers = middlewares
	e.Router.Handler(ctx)

}
