package main

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

type Status struct {
	Response    string
	Description string
}

const DELAY_ON = true

func simulateDelay(delay int) {
	if DELAY_ON {
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// TODO: Ensure unique entries based on topic
// TODO: Fix case sensitivity
// TODO: Ensure to filter out unserious entries/weird characters etc
func addTopic(w http.ResponseWriter, r *http.Request) {
	simulateDelay(2)
	fmt.Println("URL accessed: ", r.URL)
	fmt.Println("Queries are: ", r.URL.Query())
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
func searchTopics(w http.ResponseWriter, r *http.Request) {
	simulateDelay(2)
	fmt.Println("URL accessed: ", r.URL)
	fmt.Println("Queries are: ", r.URL.Query())

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

func getResources(w http.ResponseWriter, r *http.Request) {
	simulateDelay(2)
	fmt.Println("URL accessed: ", r.URL)
	fmt.Println("Queries are: ", r.URL.Query())

	queries := r.URL.Query()

	// Sets the search key
	topicIdArr := queries["id"]
	topicId := strings.Join(topicIdArr, "")
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
			LIMIT $2
		`, topicId, limit)
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

func addResource(w http.ResponseWriter, r *http.Request) {
	simulateDelay(2)
	fmt.Println("URL accessed: ", r.URL)
	fmt.Println("Queries are: ", r.URL.Query())
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

func main() {
	models.InitDB()
	//models.PrintTopics()
	//HTTP router example
	http.HandleFunc("/topics/add", addTopic)
	http.HandleFunc("/topics/search", searchTopics)
	http.HandleFunc("/resources/add", addResource)
	http.HandleFunc("/resources/get", getResources)
	//http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/save/", saveHandler)
	//http.HandleFunc("/changeBody/", changeBodyHandler)
	//http.HandleFunc("/goodListener/", goodListenerHandler)
	//http.HandleFunc("/fruits/", fruitsHandler)
	//http.HandleFunc("/getUser/", getUserHandler)
	//http.HandleFunc("/gettopic/", getTopicHandler)
	if DELAY_ON {
		fmt.Println("Simulated latency is turned ON")
	}
	fmt.Println("Serving at http://localhost:8091/ + /search, /add")
	http.ListenAndServe(":8091", nil)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}
