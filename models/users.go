package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	Id           int    `json:"id"`
	Username     string `json:'username'`
	Email        string `json:'email'`
	PasswordHash string `json:password`
}

func initUsersDb() {
	// Init topics table if not exists
	stmt, err := Db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
			username text NOT NULL,
			email text NOT NULL,
			password_hash text NOT NULL,
			created_at date
		)
	`)
	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func CreateUser(username string, email string, password string) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
		os.Exit(1)
	}
	stmt, err := Db.Prepare(`
		INSERT INTO topics (
			username, email, password_hash, created_at
		) VALUES (
			$1, $2, $3, $4
		)
	`,
	)
	if err != nil {
		log.Panic(err)
	}
	stmt.Exec(
		username,
		email,
		string(hash),
		time.Now(),
	)

	fmt.Println("Hash to store:", string(hash))
	// Store this "hash" somewhere, e.g. in your database

}

//func verifyUser(username string, password string) bool {
//	// After a while, the user wants to log in and you need to check the password he entered
//	hashFromDatabase := hash
//
//	// Comparing the password with the hash
//	if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(password)); err != nil {
//		// TODO: Properly handle error
//		log.Fatal(err)
//	}
//
//	fmt.Println("Password was correct!")
//}
