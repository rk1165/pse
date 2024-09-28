package main

import (
	"github.com/rk1165/pse/ui"
	"io/fs"
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	staticFiles, _ := fs.Sub(ui.Files, "static")
	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /index", app.index)
	mux.HandleFunc("POST /search", app.lookup)
	mux.HandleFunc("POST /submit", app.submit)
	return mux
}
