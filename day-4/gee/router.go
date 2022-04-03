package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*Node
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*Node, 0),
		handlers: make(map[string]HandlerFunc, 0),
	}
}
func (r *Router) addRouter(pattern, method string, handlerFunc HandlerFunc) {
	key := method + "-" + pattern
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &Node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handlerFunc
}

func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	vs := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			vs = append(vs, part)
			if vs[0] == "*" {
				break
			}
		}
	}
	return vs
}

func (r *Router) getRoute(method, pattern string) (n *Node, param map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	parts := parsePattern(pattern)
	n = root.search(parts, 0)
	if n != nil {
		param := make(map[string]string, 0)
		partsParam := parsePattern(n.partner)
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
	n, p := r.getRoute(ctx.method, ctx.Path)
	if n != nil {
		key := ctx.method + "-" + n.partner
		ctx.param = p
		r.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 not found%s\n", ctx.Path)
	}
}
