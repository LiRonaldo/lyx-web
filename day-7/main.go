package main

import (
	"fmt"
	"html/template"
	"log"
	"lyx-web/day-7/gee"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type student struct {
	Name string
	Age  int8
}
func main() {
	r := gee.New()
	r.Use(gee.Logger(),Recover())


	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("/Users/qt/go/src/lyx-web/day-6/templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.Get("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.Get("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.Get("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})
	r.Get("/panic", func(ctx *gee.Context) {
		names:=[]string{"liyx"}
		fmt.Println(names[3])
	})

	r.Run(":9999")
}

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}


func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func Recover()gee.HandlerFunc  {
	return func(ctx *gee.Context) {
		defer func() {
			if err:=recover();err!=nil{
				message:=fmt.Sprintf("%s",err)
				log.Printf(trace(message))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}