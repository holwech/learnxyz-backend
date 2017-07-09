package main

import (
	"encoding/json"
	"fmt"
	"github.com/holwech/learnxyz-backend/models"
	"github.com/holwech/learnxyz-backend/router"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var cred map[string]string
	file, e := ioutil.ReadFile("./cred/cred.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		fmt.Println("Create a folder '/cred/cred.json' with required login credentials")
		fmt.Println("cred.json should contain values 'psqlUsername', 'psqlPassword' and 'psqlDbName'")
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	DELAY_ON := true
	models.InitDB(cred["psqlUsername"], cred["psqlPassword"], cred["psqlDbName"])
	router := router.NewRouter()

	if DELAY_ON {
		fmt.Println("Simulated latency is turned ON")
	}

	fmt.Println("Serving at http://localhost:8091/ + /search, /add")
	log.Fatal(http.ListenAndServe(":8091", router))
}
