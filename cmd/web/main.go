package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/pse/internal/models"
	"github.com/rk1165/pse/pkg/logger"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	post          models.PostModelInterface
	request       models.RequestModelInterface
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	db            *sql.DB
	session       *sessions.CookieStore
}

const ItemsPerPage = 10

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dbName := flag.String("db", "engine.db", "SQLite Datasource name")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbName)
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

	app := &application{
		errorLog:      logger.ErrorLog,
		infoLog:       logger.InfoLog,
		post:          &models.PostModel{DB: db},
		request:       &models.RequestModel{DB: db},
		templateCache: templateCache,
		formDecoder:   form.NewDecoder(),
		db:            db,
		session:       sessions.NewCookieStore([]byte("secret")),
	}

	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     logger.ErrorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.InfoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	logger.ErrorLog.Fatal(err)

}
