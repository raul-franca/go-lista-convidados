package main

import (
	"fmt"
	"html/template"
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

}
