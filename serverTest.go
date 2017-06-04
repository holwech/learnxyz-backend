package main

import (
	"encoding/json"
	"fmt"
	"github.com/holwech/learnxyz-backend/models"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	fmt.Println("Loading page " + filename)
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

func main() {
	models.InitDB()
	models.DeleteAllAndPopulateWithTopics()
	//models.PrintTopics()
	//HTTP router example
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/topics/search", searchTopics)
	http.HandleFunc("/topics/add", addTopic)
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
