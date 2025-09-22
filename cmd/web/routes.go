package main

import (
	"io/fs"
	"net/http"

	"github.com/justinas/alice"
	"github.com/rk1165/pse/ui"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	staticFiles, _ := fs.Sub(ui.Files, "static")
	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /index", app.index)
	mux.HandleFunc("POST /search", app.search)
	mux.HandleFunc("POST /submit", app.submit)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
