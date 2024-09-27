package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type RequestModelInterface interface {
	Insert(request *Request) error
}

type Request struct {
	Url     string
	Title   string
	Links   string
	Content string
}

type RequestModel struct {
	DB *sql.DB
}

func (r *RequestModel) Insert(request *Request) error {
	log.Printf("inserting data for request %+v", request)
	stmt := `INSERT INTO requests(url, title, links, content) VALUES(?, ?, ?, ?)`
	_, err := r.DB.Exec(stmt, request.Url, request.Title, request.Links, request.Content)
	if err != nil {
		return err
	}
	log.Println("request Data inserted successfully")
	return nil
}
