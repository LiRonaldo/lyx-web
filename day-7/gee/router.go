package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (router *Router) addRouter(pattern, method string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-"+pattern
	root, ok := router.roots[method]
	if !ok {
		root = &node{}
	}
	root.insert(parts, pattern, 0)
	router.roots[method] = root
	router.handlers[key] = handler
}

func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	vs := strings.Split(pattern, "/")
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (router *Router) Handler(ctx *Context) {
	root, params := router.getRouter(ctx.Method, ctx.Path)
	if root != nil {
		key := ctx.Method + "-" + ctx.Path
		ctx.Params = params
		ctx.handlers = append(ctx.handlers, router.handlers[key])

	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		})
	}
	ctx.Next()
}

func (router *Router) getRouter(method, reqUrl string) (*node, map[string]string) {
	root, ok := router.roots[method]
	if !ok {
		return nil, nil
	}
	searchPaths := parsePattern(reqUrl)
	params := make(map[string]string, 0)

	n := root.search(searchPaths, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPaths[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchPaths[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil

}
