package main

import (
	"html/template"
	"io/fs"
	"myuto.net/snippetbox/internal/models"
	"myuto.net/snippetbox/ui"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
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
			"html/partials/*.tmpl",
			page,
		}

		// Create a new template set by parsing the 'base.tmpl' template.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		// Add the 'layout' template to the template set.
		//ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		//if err != nil {
		//	return nil, err
		//}
		//// Parse the 'page' template itself.
		//ts, err = ts.ParseFiles(page)
		//if err != nil {
		//	return nil, err
		//}

		cache[name] = ts
	}
	return cache, nil
}
