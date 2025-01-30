package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates")
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
