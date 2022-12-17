package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load .env file for dsn string
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

// mocked data articles
var articles = Articles{
	{Title: "Test Title", Desc: "Test Description", Content: "Hello World"},
	{Title: "Test Title 2", Desc: "Test Description 2", Content: "Hello World 2"},
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func allArticles(w http.ResponseWriter, r *http.Request) {
	pl("Endpoint Hit: all Articles")
	// pretty print json with SetIndent
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(articles); err != nil {
		panic(err)
	}
}

func newArticle(w http.ResponseWriter, r *http.Request) {
	pl("Endpoint Hit: new Article")
	w.Header().Set("Content-Type", "application/json")
	var article Article
	_ = json.NewDecoder(r.Body).Decode(&article)
	articles = append(articles, article)
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := getDB()
	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}
	pl("Endpoint Hit: all Users")
	results, err := db.Query("SELECT id,name,email,role FROM users")
	if err != nil {
		panic(err.Error())
	}
	// results *sql.Rows
	for results.Next() {
		user := User{}
		users := []User{}
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
		enc := json.NewEncoder(w)
		// pretty print json with SetIndent
		enc.SetIndent("", "    ")
		if err := enc.Encode(users); err != nil {
			panic(err)
		}
	}
	defer db.Close()
}

// connect to database test connection and return db
func getDB() (*sql.DB, error) {
	LoadEnv()
	dsn := os.Getenv("DB_URL")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error verifying connection with db.Ping")
		panic(err.Error())
	}
	fmt.Println("Successfully connected to database")
	return db, err
}

var pl = fmt.Println
var ff = fmt.Fprintf

func homePage(w http.ResponseWriter, r *http.Request) {
	ff(w, "Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/articles", newArticle).Methods("POST")

	pl("Starting server on port 3000")
	if err := http.ListenAndServe(":3000", myRouter); err != nil {
		log.Fatal(err)
	}
}

func main() {
	handleRequests()
}
