package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/holwech/learnxyz-backend/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Status struct {
	Response    string
	Description string
}

/*
func cleanQueries(queries url.Values) {
	for key, query := range queries {
		for i, val := range query {
			val[i] = strings.TrimSpace(strings.ToLower(val[i]))
		}
	}
}
*/

const DELAY_ON = true

func simulateDelay(delay int) {
	if DELAY_ON {
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// TODO: Ensure unique entries based on topic
// TODO: Fix case sensitivity
// TODO: Ensure to filter out unserious entries/weird characters etc
func AddTopic(w http.ResponseWriter, r *http.Request) {
	simulateDelay(1)
	queries := r.URL.Query()
	stmt, err := models.Db.Prepare(`
		INSERT INTO topics (
			topic, discipline, sub_discipline, field, description, install_date
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`,
	)
	if err != nil {
		log.Panic(err)
	}
	_, err = stmt.Exec(
		strings.Join(queries["topic"], ""),
		strings.Join(queries["discipline"], ""),
		strings.Join(queries["subDiscipline"], ""),
		strings.Join(queries["field"], ""),
		strings.Join(queries["description"], ""),
		time.Now(),
	)
	checkErr(err)
	status := Status{"Success", "Topic inserted"}
	if err != nil {
		status = Status{"Failed", "Insertion failed"}
	}
	responseUJson, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}

// TODO: Fix case sensitivity
func GetTopics(w http.ResponseWriter, r *http.Request) {
	simulateDelay(1)

	queries := r.URL.Query()
	fmt.Println(queries)

	// Sets the search key
	search := strings.Join(queries["search"], "")

	// Number of results returned
	limit := 9
	sLimit, ok := queries["limit"]
	if ok {
		limit, err := strconv.ParseInt(strings.Join(sLimit, ""), 10, 0)
		checkErr(err)
		if limit > 100 {
			limit = 100
		}
	}
	topics := models.GetTopics(search, queries["subDisciplineFilter[]"], limit)
	responseUJson, _ := json.Marshal(topics)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}

// TODO: Fix case sensitivity
func SearchTopics(w http.ResponseWriter, r *http.Request) {
	simulateDelay(1)

	queries := r.URL.Query()

	// Sets the search key
	sSearch := queries["search"]
	search := strings.Join(sSearch, "")
	// Number of results returned
	limit := 20
	sLimit, ok := queries["limit"]
	if ok {
		limit, err := strconv.ParseInt(strings.Join(sLimit, ""), 10, 0)
		checkErr(err)
		if limit > 100 {
			limit = 100
		}
	}

	rows, err := models.Db.Query("SELECT * FROM topics WHERE topic LIKE '%' || $1 || '%' LIMIT $2", search, limit)
	checkErr(err)
	topics := []models.Topic{}
	for rows.Next() {
		var topic models.Topic
		err = rows.Scan(
			&topic.Id,
			&topic.Topic,
			&topic.Discipline,
			&topic.SubDiscipline,
			&topic.Field,
			&topic.Description,
			&topic.Date,
		)
		topics = append(topics, topic)
		checkErr(err)
	}
	responseUJson, _ := json.Marshal(topics)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}
