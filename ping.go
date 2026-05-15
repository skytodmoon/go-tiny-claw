package main

import (
	"fmt"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}