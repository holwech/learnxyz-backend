package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	Db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}
}
