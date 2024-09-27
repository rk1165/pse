package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SearchModelInterface interface {
	Insert(title, url, content string) error
	Find(term string, offset int) ([]SearchResult, error)
}

type SearchResult struct {
	Title   string
	Url     string
	Content string
}

type SearchModel struct {
	DB *sql.DB
}

func (s *SearchModel) Insert(title, url, content string) error {
	stmt := `INSERT INTO posts(title, url, content) VALUES (?, ?, ?)`
	_, err := s.DB.Exec(stmt, title, url, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *SearchModel) Find(term string, offset int) ([]SearchResult, error) {

	stmt := `SELECT title, url, content FROM posts WHERE posts match ? order by rank limit 10 offset ?`
	rows, err := s.DB.Query(stmt, term, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		searchResult := SearchResult{}
		err = rows.Scan(&searchResult.Title, &searchResult.Url, &searchResult.Content)
		if err != nil {
			return nil, err
		}
		searchResult.Content = searchResult.Content[:200] // just show the first 200 characters
		results = append(results, searchResult)
	}
	log.Printf("found results: %v", len(results))
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil

}
