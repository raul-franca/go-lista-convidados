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

type formData struct {
	*Rsvp
	Errors []string
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

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}
func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

func formHandler(writer http.ResponseWriter, request *http.Request) {

	// se request = method GET responde com o form vazio
	// que esta neste caso está vindo de welcome.html
	if request.Method == http.MethodGet {
		err := templates["form"].Execute(writer, formData{
			Rsvp:   &Rsvp{},
			Errors: []string{},
		})
		if err != nil {
			panic(err)
		}
		// se request = method POST popula a struct Rsvp com os dados vindos do form
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		// valida os campos caso algum nao tenha sido preenchido retorna com o erro para form.hml
		var errors []string
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}
		if len(errors) > 0 {
			templates["form"].Execute(writer, formData{
				Rsvp:   &responseData,
				Errors: errors,
			})
		} else {

			responses = append(responses, &responseData)

			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
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
