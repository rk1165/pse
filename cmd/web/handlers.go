package main

import (
	"fmt"
	"github.com/rk1165/pse/internal/models"
	"net/http"
	"strconv"
)

type indexingForm struct {
	Url     string `form:"url"`
	Title   string `form:"title"`
	Links   string `form:"links"`
	Content string `form:"content"`
}

type PaginatedResult struct {
	Results     []models.Post
	Query       string
	NextPage    int
	PrevPage    int
	CurrentPage int
	HasNextPage bool
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "index.tmpl", nil)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("error parsing form %v", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	searchTerm := r.Form.Get("q")
	pageNo := r.Form.Get("page")

	currPage, err := strconv.Atoi(pageNo)
	if err != nil || currPage < 1 {
		currPage = 1
	}
	offset := (currPage - 1) * ItemsPerPage
	prevPage := max(1, currPage-1)
	nextPage := currPage + 1

	app.infoLog.Printf("searchTerm: %s", searchTerm)

	searchResults, err := app.post.Find(searchTerm, offset)
	if err != nil {
		err = fmt.Errorf("searching returned error %v", err)
		app.serverError(w, err)
		return
	}

	hasNextPage := len(searchResults) == ItemsPerPage
	paginatedResults := &PaginatedResult{
		Results:     searchResults,
		CurrentPage: currPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		Query:       searchTerm,
		HasNextPage: hasNextPage,
	}

	app.renderPartial(w, http.StatusOK, "results.tmpl", paginatedResults)

}

func (app *application) submit(w http.ResponseWriter, r *http.Request) {
	var form indexingForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.infoLog.Printf("form: %+v", form)
	request := &models.Request{Url: form.Url, Title: form.Title, Links: form.Links, Content: form.Content}

	err = app.request.Insert(request)
	if err != nil {
		err = fmt.Errorf("error while inserting %v", err)
		app.serverError(w, err)
		return
	}

	go func() {
		ch := make(chan int, 1)
		Index(request, app, ch)
		status := <-ch
		if status != http.StatusOK {
			err = fmt.Errorf("error while indexing %v", err)
			app.serverError(w, err)
			return
		}
	}()

	session, _ := app.session.Get(r, "flash-session")
	session.Values["flash"] = fmt.Sprintf("Started Indexing %s", form.Url)
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// TODO: page showing the currently indexed pages, also to delete them
