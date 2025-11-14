package main

import (
	"fmt"

	"log"
	"net/http"

)

func main() {
	http.HandleFunc("/", homePage)
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}
