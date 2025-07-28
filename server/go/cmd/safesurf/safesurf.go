package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abhizaik/SafeSurf/internal/service/checks"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "SafeSurf")
		fmt.Fprintln(w, checks.GetStatusCode("https://google.com"))

	})
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
