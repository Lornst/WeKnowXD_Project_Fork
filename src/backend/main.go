package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/search.html", "templates/register.html"))


func main() {
	mux := http.NewServeMux()
	router(mux)
	http.ListenAndServe(":8080", mux)
}

type PageData struct {
	Query   string
	Results []struct {
		Title       string
		URL         string
		Description string
	}
}

type RegisterData struct {
	Error    string
	Username string
	Email    string
}

func router(mux *http.ServeMux) {
	mux.HandleFunc("GET /register", registerHandler)
	mux.HandleFunc("GET /search", searchHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}


func searchHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Query: r.URL.Query().Get("q")}
	templates.ExecuteTemplate(w, "search.html", data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", RegisterData{})
}