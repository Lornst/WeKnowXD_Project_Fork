package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	//fmt.Println("Hello, world!")
	router()
}

func router() {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/search", getSearch)
	router.HandleFunc("POST /api/register", postRegister)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}

func getSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	_ = q // Use the query parameter 'q' as needed
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}

	response := map[string]any{
		"search_results": []string{},
		//Data suppose to be written here
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Register endpoint hit")
}
