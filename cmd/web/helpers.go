package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data interface{}) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base.tmpl", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)

}

func (app *application) renderPartial(w http.ResponseWriter, status int, page string, data interface{}) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, page, data)
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

func (app *application) newTemplateData(w http.ResponseWriter, r *http.Request) *templateData {
	session, _ := app.session.Get(r, "flash-session")
	flash, ok := session.Values["flash"].(string)
	app.infoLog.Printf("flash=%s", flash)
	if ok {
		delete(session.Values, "flash")
		session.Save(r, w)
	}
	return &templateData{
		Flash: flash,
	}
}
