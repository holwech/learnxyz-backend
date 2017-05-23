package models

import (
	"fmt"
	"log"
)

type Topic struct {
	Name       string
	RelatedUrl []string
}

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

	//bks := make([]*Book, 0)
	//for rows.Next() {
	//bk := new(Book)
	//err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
	//if err != nil {
	//return nil, err
	//}
	//bks = append(bks, bk)
	//}
	//if err = rows.Err(); err != nil {
	//return nil, err
	//}
	//return bks, nil
}

func printRow(tp Topic) {
	fmt.Println("Name: ", tp.Name, ", RelatedUrl: ", tp.RelatedUrl)
}
