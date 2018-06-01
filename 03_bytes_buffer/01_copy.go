package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func curl1() {
	r, err := http.Get("http://google.com")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, r.Body)
	if err = r.Body.Close(); err != nil {
		log.Fatalln(err)
	}
}
