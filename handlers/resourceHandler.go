package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/holwech/learnxyz-backend/models"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func GetResources(w http.ResponseWriter, r *http.Request) {
	simulateDelay(1)

	queries := r.URL.Query()

	// Sets the search key
	topicId := strings.Join(queries["topicId"], "")
	label := strings.Join(queries["type"], "")
	fmt.Println("label is", label)

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

	rows, err := models.Db.Query(`
			SELECT * FROM resources
			WHERE $1 = ANY (related_topic_ids)
			AND type = $2
			LIMIT $3
		`, topicId, label, limit)
	checkErr(err)
	resources := []models.Resource{}
	for rows.Next() {
		var resource models.Resource
		err = rows.Scan(
			&resource.Id,
			&resource.Title,
			&resource.Url,
			&resource.Description,
			&resource.Type,
			pq.Array(&resource.Tags),
			pq.Array(&resource.RelatedTopicId),
			&resource.Date,
		)
		resources = append(resources, resource)
		checkErr(err)
	}
	responseUJson, _ := json.Marshal(resources)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}

func AddResource(w http.ResponseWriter, r *http.Request) {
	simulateDelay(1)
	queries := r.URL.Query()
	stmt, err := models.Db.Prepare(`
		INSERT INTO resources (
			title, url, description, type, tags, related_topic_ids, install_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`,
	)
	if err != nil {
		log.Panic(err)
	}
	relatedTopicId, _ := strconv.ParseInt(strings.Join(queries["relatedTopicId"], ""), 10, 64)
	relatedTopicIdArr := [1]int64{relatedTopicId}
	fmt.Println(queries["tags"])
	fmt.Println("Related topic id is: ", relatedTopicId)
	_, err = stmt.Exec(
		strings.Join(queries["title"], ""),
		strings.Join(queries["url"], ""),
		strings.Join(queries["description"], ""),
		strings.Join(queries["type"], ""),
		pq.StringArray(queries["tags"]),
		pq.Array(relatedTopicIdArr),
		time.Now(),
	)
	checkErr(err)
	status := Status{"Success", "Resource inserted"}
	if err != nil {
		status = Status{"Failed", "Insertion failed"}
	}
	responseUJson, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}
