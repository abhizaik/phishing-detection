package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "SafeSurf")
	})
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
