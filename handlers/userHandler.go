package handlers

import (
	"encoding/json"
	"github.com/holwech/learnxyz-backend/models"
	"net/http"
	"strings"
)

/*
func cleanQueries(queries url.Values) {
	for key, query := range queries {
		for i, val := range query {
			val[i] = strings.TrimSpace(strings.ToLower(val[i]))
		}
	}
}
*/

func CreateUser(w http.ResponseWriter, r *http.Request) {
	simulateDelay(1)
	queries := r.URL.Query()
	username := strings.Join(queries["username"], "")
	email := strings.Join(queries["email"], "")
	password := strings.Join(queries["password"], "")
	models.CreateUser(username, email, password)

	status := Status{"Success", "Topic inserted"}
	responseUJson, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}
