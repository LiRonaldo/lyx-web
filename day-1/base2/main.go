package main

import (
	"fmt"
	"net/http"
)

type Engine struct {

}

func main()  {
	engine:=Engine{}
	http.ListenAndServe(":9999",&engine)
}
func (e *Engine)ServeHTTP(resp http.ResponseWriter, req *http.Request){
	switch req.URL.Path {
	case "/head":
		for k, v := range req.Header {
			fmt.Fprintf(resp, "headr[%q]=%q\n", k, v)
		}
	case "hello":
		fmt.Fprintf(resp,"%s",req.URL.Path)
	default:
		fmt.Fprintf(resp,"404 not found %s",req.URL.Path)
	}

}
