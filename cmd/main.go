package main

import (
	"log"
	"net/http"

	json2protobuff "github.com/anmit007/JSON2PROTOBUFF"
)

func main() {
	http.HandleFunc("/convert", json2protobuff.WebHandler)
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
