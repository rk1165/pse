package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"html/template"
	"net/http"
	"runtime/debug"
)

const BasePath = "ui/html/pages/base.tmpl"

func (app *application) render(w http.ResponseWriter, status int, page string, data interface{}) {

	pagePath := "ui/html/pages/" + page
	t, err := template.ParseFiles(BasePath, pagePath)
	if err != nil {
		err := fmt.Errorf("cannot parse %s template", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)

}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}
