package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/pse/pkg/logger"
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
	logger.InfoLog.Printf("inserting data for request %+v", request)
	stmt := `INSERT INTO requests(url, title, links, content) VALUES(?, ?, ?, ?)`
	_, err := r.DB.Exec(stmt, request.Url, request.Title, request.Links, request.Content)
	if err != nil {
		return err
	}
	logger.InfoLog.Println("request Data inserted successfully")
	return nil
}
