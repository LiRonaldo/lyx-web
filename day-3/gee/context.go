package gee

import (
	"fmt"
	"net/http"
)

type Context struct {
	req    *http.Request
	writer http.ResponseWriter
	path   string
	method string
	param  map[string]string
}

func NewContext(req *http.Request, writer http.ResponseWriter) *Context {
	return &Context{
		req:    req,
		writer: writer,
		path:   req.URL.Path,
		method: req.Method,
	}
}
func (ctx *Context) Param(key string) string {
	return ctx.param[key]
}

func (ctx *Context) Query(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *Context) String(status int, format string, values ...interface{}) {
	ctx.SetStatus(status)
	ctx.writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (ctx *Context) SetStatus(status int) {
	ctx.writer.WriteHeader(status)
}
func (ctx *Context) SetHeader(key, value string) {
	ctx.writer.Header().Set(key, value)
}
func (ctx *Context) Html(status int, html string) {
	ctx.SetStatus(status)
	ctx.SetHeader("Content-Type", "text/html")
	ctx.writer.Write([]byte(html))
}
