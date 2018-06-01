package main

import (
	"fmt"
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func server2() {
	h := http.HandlerFunc(mainHandler)

	err := http.ListenAndServe(":9999", h)
	log.Fatal(err)
}
