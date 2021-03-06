package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

// Topic
type Topic struct {
	Id            int       `json:"id"`
	Topic         string    `json:"topic"`
	Discipline    string    `json:"discipline"`
	SubDiscipline string    `json:"subDiscipline"`
	Field         string    `json:"field"`
	Description   string    `json:"description"`
	Date          time.Time `json:"date"`
}

func initTopicsDb() {
	// Init topics table if not exists
	stmt, err := Db.Prepare(`
		CREATE TABLE IF NOT EXISTS topics (
			id serial PRIMARY KEY,
			topic text NOT NULL,
			discipline text NOT NULL,
			sub_discipline text,
			field text,
			description text,
			install_date date
		)
	`)
	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

// Loads the json data from topics.json and loads it into the
// postgres topics table.
func DeleteAllAndPopulateWithTopics() {
	//Deletes all rows in topics table
	_, err := Db.Exec("DELETE FROM topics")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file, err := ioutil.ReadFile("./topics.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var data []Topic
	json.Unmarshal(file, &data)
	fmt.Println("- - - - - - - - - -")
	fmt.Println("DB:Topics have been populated with the following topics")
	for _, element := range data {
		fmt.Println(": " + element.Topic)
		InsertTopic(element)
	}
}

func InsertTopic(topic Topic) {
	stmt, err := Db.Prepare(`
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
	stmt.Exec(
		topic.Topic,
		topic.Discipline,
		topic.SubDiscipline,
		topic.Field,
		topic.Description,
		time.Now(),
	)
}

func GetTopics(search string, subDisciplineFilter []string, limit int) []Topic {
	var rows *sql.Rows
	var err error
	if len(subDisciplineFilter) > 0 {
		arr := pq.StringArray(subDisciplineFilter)
		rows, err = Db.Query(`
			SELECT * FROM topics
			WHERE (length($1) = 0 OR topic LIKE '%' || $1 || '%')
			AND sub_discipline = ANY($2::text[])
			LIMIT $3`,
			search, arr, limit)
		checkErr(err)
	} else {
		rows, err = Db.Query(`
			SELECT * FROM topics
			WHERE (length($1) = 0 OR topic LIKE '%' || $1 || '%')
			LIMIT $2`,
			search, limit)
		checkErr(err)
	}
	topics := []Topic{}
	for rows.Next() {
		var topic Topic
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
	return topics
}

/*
func PrintTopics() {
	rows, err := Db.Query("SELECT * FROM topics")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	tp := Topic{}
	for rows.Next() {
		err := rows.Scan(&tp.Name, &tp.RelatedUrl)
		if err != nil {
			log.Println(err.Error())
		}
		printRow(tp)
	}

}

func printRow(tp Topic) {
	fmt.Println("Topic: ", tp.Topic, ", RelatedUrl: ", tp.RelatedUrl)
}
*/
