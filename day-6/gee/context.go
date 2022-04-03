package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type HandlerFunc func(ctx *Context)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	handlers   []HandlerFunc
	StatusCode int
	Params     map[string]string
	Method     string
	Path       string
	index      int
	engine *Engine
}

func newContext(req *http.Request, resp http.ResponseWriter) *Context {
	return &Context{
		Req:    req,
		Resp:   resp,
		index:  -1,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (ctx *Context) next() {
	ctx.index++
	s := len(ctx.handlers)
	for ; ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.setHeader("Content-Type", "text/html")
	c.setStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Resp, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (ctx *Context) setStatus(status int) {
	ctx.StatusCode = status
	ctx.Resp.WriteHeader(status)
}
func (ctx *Context) setHeader(key, value string) {
	ctx.Resp.Header().Set(key, value)
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.setStatus(code)
	encoder := json.NewEncoder(c.Resp)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Resp, err.Error(), 500)
	}
}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.setStatus(code)
	c.Resp.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}