package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
	w.Write(jsonFruits)
}

func main() {
	//HTTP router example
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/changeBody/", changeBodyHandler)
	//http.HandleFunc("/goodListener/", goodListenerHandler)
	http.HandleFunc("/fruits/", fruitsHandler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
