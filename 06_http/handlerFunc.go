package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func main() {
	h := http.HandlerFunc(handler)

	err := http.ListenAndServe(":9999", h)
	log.Fatal(err)
}
