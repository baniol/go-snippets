package main

import (
	"fmt"
	"log"
	"net/http"
)

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s", r.Host)
	fmt.Fprint(rw, "okidoki")
}

func main() {
	http.HandleFunc("/", mainHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
