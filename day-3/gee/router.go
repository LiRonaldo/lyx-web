package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	node    map[string]*Node
	handler map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		node:    make(map[string]*Node),
		handler: make(map[string]HandlerFunc),
	}

}
func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.node[method]; !ok {
		r.node[method] = &Node{}
	}
	r.node[method].insert(pattern, parts, 0)
	r.handler[key] = handler
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
func (r *Router) getRoute(method, pattern string) (n *Node, param map[string]string) {
	root, ok := r.node[method]
	if !ok {
		return nil, nil
	}
	parts := parsePattern(pattern)
	n = root.search(parts, 0)
	if n != nil {
		param := make(map[string]string, 0)
		partsParam := parsePattern(n.pattern)
		for index, part := range partsParam {
			if part[0] == ':' {
				param[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				param[part[1:]] = strings.Join(parts[index:], "/")
				break
			}
		}
		return n, param
	}
	return nil, nil
}

func (r *Router) Handler(ctx *Context) {
	n, p := r.getRoute(ctx.method, ctx.path)
	if n != nil {
		k := ctx.method + "-" + n.pattern
		ctx.param = p
		r.handler[k](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.path)
	}
}
