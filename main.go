package main

import (
	"fmt"
	"github.com/holwech/learnxyz-backend/models"
	"github.com/holwech/learnxyz-backend/router"
	"log"
	"net/http"
)

func main() {
	DELAY_ON := true
	models.InitDB()
	router := router.NewRouter()

	if DELAY_ON {
		fmt.Println("Simulated latency is turned ON")
	}

	fmt.Println("Serving at http://localhost:8091/ + /search, /add")
	log.Fatal(http.ListenAndServe(":8091", router))
}
