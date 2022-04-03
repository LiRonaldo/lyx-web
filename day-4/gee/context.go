package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HandlerFunc func(ctx *Context)
type H map[string]interface{}

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	Path       string
	param      map[string]string
	StatusCode int
	method     string
}

func NewContext(req *http.Request, resp http.ResponseWriter) *Context {
	return &Context{
		Req:    req,
		Resp:   resp,
		Path:   req.URL.Path,
		method: req.Method,
	}
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Resp.Header().Set(key, value)
}
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Resp.WriteHeader(code)
}
func (ctx *Context) String(code int, format string, value ...interface{}) {
	ctx.Status(code)
	ctx.Resp.Write([]byte(fmt.Sprintf(format, value...)))
}

func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	ctx.Resp.Write([]byte(html))
}

func (ctx *Context) Param(key string) string {
	return ctx.param[key]
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Resp)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Resp, err.Error(), 500)
	}
}