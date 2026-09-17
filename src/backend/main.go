package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"log"
	"net/http"

	"time"

	"sync"

	"fmt"

	"reflect"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dataBase := newDataBase()
	defer dataBase.Close()

	server := newServer(dataBase)

	http.ListenAndServe(":8080", server.Mux)
}

func (server *Server) router() {
	server.Mux.HandleFunc("POST /api/login", server.apiLogin)
	server.Mux.HandleFunc("GET /api/test", server.apiTestDB)

	fileServer := http.FileServer(http.Dir("./templates"))

	server.Mux.Handle("/templates/", http.StripPrefix("/templates/", fileServer))

	server.Mux.Handle("/", fileServer)
}

func (server *Server) apiTestDB(w http.ResponseWriter, r *http.Request) {
	userList := server.queryDB(reflect.TypeOf(User{}), "SELECT * FROM users")

	if len(userList) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("database connected"))
		return
	}

	firstUser := userList[0].(User)
	println("Test succes, found user:", firstUser.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

type Server struct {
	Mux           *http.ServeMux
	DB            *sql.DB
	Sessions      map[string]int
	SessionsMutex sync.RWMutex
}

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}

func newServer(dataBase *sql.DB) *Server {
	server := &Server{
		DB:       dataBase,
		Mux:      http.NewServeMux(),
		Sessions: make(map[string]int),
	}

	server.router()
	return server
}

func newDataBase() *sql.DB {
	db, err := sql.Open("sqlite3", "../whoknows.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to SQLite database")

	return db
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (server *Server) apiLogin(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	userList := server.queryDB(reflect.TypeOf(User{}), "SELECT * FROM users WHERE username = ?", r.FormValue("username"))
	if len(userList) == 0 {
		log.Println("func: apiLogin, queryDB returned empty userlist when seaching for username: " + r.FormValue("username"))
	}

	foundUser := userList[0].(User)

	if verifyPassword(foundUser.Password, r.FormValue("password")) == false {
		log.Println("func: apiLogin, user verifacation password missmatch")
	} else {
		token := generateSessionToken()

		server.SessionsMutex.Lock()
		server.Sessions[token] = foundUser.Id
		server.SessionsMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
		})

		println("User", foundUser.Username, "is currently logged in with session")
		w.Write([]byte("Login succesfull"))
	}
}

func (server *Server) queryDB(interchangeableStruct reflect.Type, query string, args ...any) []any {
	var structArray []any
	rows, error := server.DB.Query(query, args...)

	if error != nil {
		log.Fatal(error)
	}
	defer rows.Close()

	for rows.Next() {
		newStructPointer := reflect.New(interchangeableStruct)
		structElement := newStructPointer.Elem()

		numCols := structElement.NumField()
		columns := make([]any, numCols)

		for i := range numCols {
			field := structElement.Field(i)
			columns[i] = field.Addr().Interface()
		}

		error := rows.Scan(columns...)
		if error != nil {
			log.Fatal(error)
		}

		structArray = append(structArray, structElement.Interface())
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return structArray
}

func hashPassword(password string) string {
	passwordBytes := []byte(password)
	passwordHash := md5.Sum(passwordBytes)
	passwordHashString := hex.EncodeToString(passwordHash[:])

	return passwordHashString
}

func verifyPassword(storedHash string, password string) bool {
	passwordHash := hashPassword(password)
	return storedHash == passwordHash
}
