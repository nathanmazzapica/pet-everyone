package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl := template.Must(template.ParseFiles(
		"app/templates/layout.html",
		"app/templates/header.html",
		fmt.Sprintf("app/templates/%s.html", name),
	))

	return tmpl.ExecuteTemplate(w, "layout", data)
}
