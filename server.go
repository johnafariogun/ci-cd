package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Go Web!, Ci-Cd is a go")
}

func main() {
	http.HandleFunc("/", handler) // Route all requests to the handler function
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server
}
