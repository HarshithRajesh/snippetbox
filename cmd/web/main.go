package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/HarshithRajesh/snippetbox/pkg/models/postgresql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	// db       *sql.DB
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *postgresql.SnippetModel
	users         *postgresql.UserModel
	session       *sessions.Session
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "Http network address")

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &postgresql.SnippetModel{DB: db},
		users:         &postgresql.UserModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Starting server on ", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

	}
	log.Println("Reached")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USERDB")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Check if the connection is established correctly
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
