package main

import (
	"fmt"
	"lyx-web/day-3/gee"
	"net/http"
)

func main() {
	engine := gee.NewGee()
	engine.Get("/", func(ctx *gee.Context) {
		ctx.Html(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	engine.Get("/hello", func(ctx *gee.Context) {
		ctx.Html(http.StatusOK, fmt.Sprintf("%s", ctx.Query("name")))
	})
	engine.Get("/hello/:name", func(ctx *gee.Context) {
		ctx.Html(http.StatusOK, fmt.Sprintf("%s", ctx.Param("name")))
	})
	engine.Get("/assets/*filepath", func(ctx *gee.Context) {
		ctx.Html(http.StatusOK,fmt.Sprintf("%s",ctx.Param("filepath")))
	})
	engine.RUN(":3333")
}
