package main

import (
	"fmt"
	"lyx-web/day-1/base3/gee"
	"net/http"
)

func main() {
	e := gee.NewEngine()
	e.Get("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "%s", "hello lyx")
	})
	e.Get("/head", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			fmt.Fprintf(writer, "head[%s]=%s\n", k, v)
		}
	})
	e.RUN(":1111")
}
