package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/head", head)
	log.Fatalln(http.ListenAndServe(":8888", nil))
}

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "%q\n", req.URL.Path)
}
func head(resp http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(resp, "headr[%q]=%q\n", k, v)
	}
}
