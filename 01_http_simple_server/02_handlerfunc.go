package main

import (
	"fmt"
	"log"
	"net/http"
)

func server2() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
	})

	err := http.ListenAndServe(":9999", h)
	log.Fatal(err)
}
