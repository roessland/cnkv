package main

import (
	"log"
	"net/http"
)

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, net/http!"))
}

func main() {
	http.HandleFunc("/", helloGoHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}