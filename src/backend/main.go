package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	htmlPages := http.FileServer(http.Dir("./templates"))

	mux.Handle("/", htmlPages)
	mux.HandleFunc("GET /login", login)

	http.ListenAndServe(":8080", mux)
}

func login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/login.html")
}
