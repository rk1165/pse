package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/pse/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	search      models.SearchModelInterface
	request     models.RequestModelInterface
	formDecoder *form.Decoder
	db          *sql.DB
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dbName := flag.String("db", "engine.db", "SQLite Datasource name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlite3", *dbName)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		search:      &models.SearchModel{DB: db},
		request:     &models.RequestModel{DB: db},
		formDecoder: form.NewDecoder(),
		db:          db,
	}

	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)

}
