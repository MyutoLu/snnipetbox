package main

import (
	"errors"
	"fmt"
	"html/template"
	"myuto.net/snippetbox/internal/models"
	"net/http"
	"strconv"
)

// home
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		//log.Print(err.Error())
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		//log.Print(err.Error())
		//app.errorLog.Print(err.Error())
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, err)
		return
	}
}

// snippetView 根据id查
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}
	return
	fmt.Fprintf(w, "%+v", snippet)
}

// snippetCreate 创建博文
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O sanil"
	content := "O sanil\nClimb Mount Fuji,\nBut slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//w.Write([]byte("Create a new snippet..."))

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
