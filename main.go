package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct { // RSVP responda por favor
	Name, Email, Phone string
	WillAttend         bool
}

// Vai funcionar como um banco de dados com as respostas
var responses = make([]*Rsvp, 0, 10)                   // []slice, make(type, len inicial, capacidade)
var templates = make(map[string]*template.Template, 3) //map[key]valor, tamanho
// template do tipo Template

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"} //lista dos nomes dos html/template
	for index, name := range templateNames {
		//"layout.html", name + ".html" concatena layout com os arquivos. show!!
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil { // se nao over error add a templates que um map [key] value
			templates[name] = t // map templates onde key eh name e o conteudo é um Template
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

func main() {
	loadTemplates()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}
func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}
func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		err := templates["form"].Execute(writer, formData{
			Rsvp:   &Rsvp{},
			Errors: []string{},
		})
		if err != nil {
			panic(err)
		}
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		responses = append(responses, &responseData)

		if responseData.WillAttend {
			templates["thanks"].Execute(writer, responseData.Name)
		} else {
			templates["sorry"].Execute(writer, responseData.Name)
		}

	}
}

type formData struct {
	*Rsvp
	Errors []string
}
