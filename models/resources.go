package models

import (
	"log"
	"time"
)

// Topic
type Resource struct {
	Id             int       `json:"id"`
	Title          string    `json:"title"`
	Url            string    `json:"url"`
	Description    string    `json:"description"`
	Type           string    `json"type"`
	Tags           []string  `json:"tags"`
	RelatedTopicId []int64   `json:"relatedTopicId"`
	Date           time.Time `json:"install_date"`
}

func initResourceDb() {
	// Init topics table if not exists
	stmt, err := Db.Prepare(`
		CREATE TABLE IF NOT EXISTS resources (
			id serial PRIMARY KEY,
			title text NOT NULL,
			url text NOT NULL,
			description text,
			type text NOT NULL,
			tags text[],
			related_topic_ids integer[] NOT NULL,
			install_date date
		)
	`)
	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func InsertResource(resource Resource) {
	stmt, err := Db.Prepare(`
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
	stmt.Exec(
		resource.Title,
		resource.Url,
		resource.Description,
		resource.RelatedTopicId,
		time.Now(),
	)
}
