package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/rk1165/pse/ui"
)

type templateData struct {
	Flash string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			page,
		}
		// Parse the base template file into a template set.
		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
