package models

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var Db *sql.DB

// How does this work?
// Run psql in the terminal to login into postgres.
// username should be visible as the first name in the terminal.
// Run \l to show all databases, run \dt to show all tables for
// current database, run \du to show all users.
// Do NOT login as the user called postgres.
// If the database does NOT ask for password when you run psql, your
// settings are wrong and unsecure.
func InitDB() {
	fmt.Println("Setting up DB")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	dbUser, _ := reader.ReadString('\n')
	fmt.Print("Enter password: ")
	dbPassword, _ := reader.ReadString('\n')
	fmt.Print("Enter DB name: ")
	dbName, _ := reader.ReadString('\n')
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)

	var err error
	Db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Panic(err)
	}

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}

	fmt.Println("Set up success!")
	initTopicsDb()
	initResourceDb()
	fmt.Println("DB ready!")
}
