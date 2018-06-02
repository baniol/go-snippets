package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const port = 8080

type User struct {
	Name  string `json:"username"`
	Email string `json:"email"`
}

func server1() {
	http.HandleFunc("/", handlerFunction)
	http.HandleFunc("/stream", streamHandlerFunction)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/user-stream", userHandlerStream)
	log.Printf("Listening on port :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func handlerFunction(rw http.ResponseWriter, r *http.Request) {
	user := User{"john.doe", "jon.doe@example.com"}
	str, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(rw, string(str))
}

func streamHandlerFunction(w http.ResponseWriter, r *http.Request) {
	user := User{"john.doe", "jon.doe@example.com"}
	encoder := json.NewEncoder(w)
	encoder.Encode(&user)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Println(user)
	fmt.Fprintf(w, "Received")
}

func userHandlerStream(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Println(user)
	fmt.Fprintf(w, "Received")
}
