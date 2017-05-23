package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"learnxyz-backend/models"
	"log"
	"net/http"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

type TestStruct struct {
	Message string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func changeBodyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	decoder := json.NewDecoder(r.Body)
	var t TestStruct
	err := decoder.Decode(&t)
	if err != nil {
		//panic(err)
	}
	defer r.Body.Close()
	log.Println(t.Message)
	title := r.URL.Path[len("/changeBody/"):]
	p := &Page{Title: title, Body: []byte(t.Message)}
	p.save()
}

//func goodListenerHandler(w http.ResponseWriter, r *http.Request) {
//
//r.ParseForm()
//log.Println(r.Form)
//decoder := json.NewDecoder(r.Body)
//err := decoder.Decode(&t)
//if err != nil {
////panic(err)
//}
//defer r.Body.Close()
//}

type Objects struct {
	Fruits []string
}

func fruitsHandler(w http.ResponseWriter, r *http.Request) {
	f := Objects{[]string{"Lemon", "Peach", "Jordgubbe", "Knoblauch"}}
	jsonFruits, err := json.Marshal(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonFruits)
}

type User struct {
	Uid        string
	Username   string
	Departname string
	Created    string
}

type Topic struct {
	Name       string
	RelatedUrl []string
}

func getTopicHandler(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Path[len("/gettopic/"):]
	rows, err := models.Db.Query("SELECT * FROM topics WHERE name=$1", topicName)
	checkErr(err)
	responseTopic := Topic{}

	rows.Next()
	rows.Scan(&responseTopic.Name, &responseTopic.RelatedUrl)
}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	//r.ParseForm()
	//fmt.Println("PRINT FORM: ", r.Form)
	//Decode json from request

	//queryU := User{1, "99", "99", "99"}

	queryU := User{}
	err := json.NewDecoder(r.Body).Decode(&queryU)
	if err != nil {
		log.Println("Failed to decode: ", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println("QueryUser id: ", queryU.Uid)
	var responseU User

	fmt.Println("# Inserting values")

	rows, err := models.Db.Query("SELECT * FROM userinfo WHERE uid=$1", queryU.Uid)
	checkErr(err)
	rows.Next()
	rows.Scan(&responseU.Uid, &responseU.Username, &responseU.Departname, &responseU.Created)
	fmt.Println("Fitting user: ", responseU.Username)
	responseUJson, _ := json.Marshal(responseU)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(responseUJson)
}
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "myDatabaseName"
)

func main() {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	models.InitDB(dbinfo)

	fmt.Println("# Inserting values")

	models.Db.QueryRow("CREATE TABLE IF NOT EXISTS topics(name text NOT NULL,related_url text[])")
	var lastInsertId int
	models.Db.QueryRow("INSERT INTO topics(name, related_url) VALUES($1,$2);", "linear algebra", "{\"url1\", \"url2\", \"url3\"}")
	//err := row.Scan()
	//checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Updating")
	stmt, err := models.Db.Prepare("update userinfo set username=$1 where uid=$2")
	checkErr(err)

	res, err := stmt.Exec("astaxieupdate", lastInsertId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	fmt.Println("# Querying")
	rows, err := models.Db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println("uid | username | department | created ")
		fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
	}

	fmt.Println("# Deleting")
	stmt, err = models.Db.Prepare("delete from userinfo where uid=$1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")
	models.PrintTopics()
	//HTTP router example
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/changeBody/", changeBodyHandler)
	//http.HandleFunc("/goodListener/", goodListenerHandler)
	http.HandleFunc("/fruits/", fruitsHandler)
	http.HandleFunc("/getUser/", getUserHandler)
	http.HandleFunc("/gettopic/", getTopicHandler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}
