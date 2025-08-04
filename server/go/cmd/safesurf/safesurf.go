package main

import (
	"log"

	"github.com/abhizaik/SafeSurf/internal/handler"
	"github.com/abhizaik/SafeSurf/internal/service/rank"
)

func main() {
	r := handler.SetupRouter()

	err := rank.LoadDomainRanks()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
