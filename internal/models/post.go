package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type PostModelInterface interface {
	Insert(post Post) error
	Find(term string, offset int) ([]Post, error)
}

type Post struct {
	Title   string
	Url     string
	Content string
}

type PostModel struct {
	DB *sql.DB
}

func (s *PostModel) Insert(post Post) error {
	stmt := `INSERT INTO posts(title, url, content) VALUES (?, ?, ?)`
	_, err := s.DB.Exec(stmt, post.Title, post.Url, post.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostModel) Find(term string, offset int) ([]Post, error) {

	stmt := `SELECT title, url, content FROM posts WHERE posts match ? order by rank limit 10 offset ?`
	rows, err := s.DB.Query(stmt, term, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Post

	for rows.Next() {
		searchResult := Post{}
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
