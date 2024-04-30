package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"myuto.net/snippetbox/internal/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// home
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	//for _, snippets := range snippets {
	//	fmt.Fprintf(w, "%+v\n", snippets)
	//}
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

}

// snippetView 根据id查
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)

	//fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) sinppetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// snippetCreate 创建博文
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//	return
	//}
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fileErrors := make(map[string]string)
	// title 不能为空和超过100个字符
	if strings.TrimSpace(title) == "" {
		fileErrors["title"] = "This field can't be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fileErrors["title"] = "Title can't be more than 100 characters"
	}

	// content 不能为空
	if strings.TrimSpace(content) == "" {
		fileErrors["content"] = "This field can't be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fileErrors["expires"] = "This field is invalid"
	}

	if len(fileErrors) > 0 {
		fmt.Fprint(w, fileErrors)
		return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
