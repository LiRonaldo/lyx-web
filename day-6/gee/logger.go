package gee

import (
	"fmt"
	"time"
)

func Logger() (handler HandlerFunc) {
	return func(ctx *Context) {
		t := time.Now()
		ctx.next()
		fmt.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
