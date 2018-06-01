package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// curl -H "Content-Type: application/json" -X POST -d '{"name": "Katowice", "area": 3000000}' localhost:8080/city

type city struct {
	Name string `json:"name"`
	Area uint64 `json:"area"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var inputCity city
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&inputCity)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Printf("City: %s; Area: %d", inputCity.Name, inputCity.Area)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}
}

func main() {
	http.HandleFunc("/city", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
