package handlers

import (
	"html/template"
	"net/http"
)

// GlavHandler - обработчик главной страницы
func GlavHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
