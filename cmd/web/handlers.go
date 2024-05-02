package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"myuto.net/snippetbox/internal/models"
	"myuto.net/snippetbox/internal/validator"
	"net/http"
	"strconv"
)

type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	//FiledErrors map[string]string
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	//FiledErrors         map[string]string
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

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
	//flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.newTemplateData(r)
	data.Snippet = snippet
	//data.Flash = flash
	app.render(w, http.StatusOK, "view.tmpl", data)

	//fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
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
	//err := r.ParseForm()
	//if err != nil {
	//	app.clientError(w, http.StatusBadRequest)
	//	return
	//}

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	//err = app.formDecoder.Decode(&form, r.PostForm)
	//if err != nil {
	//	app.clientError(w, http.StatusBadRequest)
	//	return
	//}

	//form := &snippetCreateForm{
	//	Title:   r.PostForm.Get("title"),
	//	Content: r.PostForm.Get("content"),
	//	Expires: expires,
	//	//FiledErrors: map[string]string{},
	//}

	form.CheckFiled(validator.NotBlank(form.Title), "title", "This field can't be blank")
	form.CheckFiled(validator.MaxChars(form.Title, 100), "title", "This field can't be more than 100 characters")
	form.CheckFiled(validator.NotBlank(form.Content), "content", "This field can't be blank")
	form.CheckFiled(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	/*// title 不能为空和超过100个字符
	if strings.TrimSpace(form.Title) == "" {
		form.FiledErrors["title"] = "This field can't be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FiledErrors["title"] = "Title can't be more than 100 characters"
	}

	// content 不能为空
	if strings.TrimSpace(form.Content) == "" {
		form.FiledErrors["content"] = "This field can't be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FiledErrors["expires"] = "This field is invalid"
	}

	if len(form.FiledErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}*/
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "成功创建")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckFiled(validator.NotBlank(form.Name), "name", "This field can't be blank")
	form.CheckFiled(validator.NotBlank(form.Email), "email", "This field can't be blank")
	form.CheckFiled(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckFiled(validator.NotBlank(form.Password), "password", "This field can't be blank")
	form.CheckFiled(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFiledError("email", "Address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckFiled(validator.NotBlank(form.Email), "email", "This field can't be blank")
	form.CheckFiled(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckFiled(validator.NotBlank(form.Password), "password", "This field can't be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) isAuthenticatedUser(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
