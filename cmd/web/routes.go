package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /index", app.index)
	mux.HandleFunc("POST /search", app.lookup)
	mux.HandleFunc("POST /submit", app.submit)
	return mux
}
