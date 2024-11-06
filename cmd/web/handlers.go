package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/shtayeb/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use newTemplateData helper
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// use the render helper
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.

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
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// w.WriteHeader can be called once per request and after the status code is set it cant be changed
	// if you dont call w.WriteHeader explicitly, then the first call to w.Write will automatically send a 200 ok
	// so if you want to send a non-200 code you must call w.WriteHeader before any call to w.Write

	// This
	// w.WriteHeader(405)
	// w.Header().Set("Allow", http.MethodPost)
	// w.Write([]byte("Method not allowed"))

	// Or helper

	// http.StatusMethodNotAllowed = 405
	// app.clientError(w, http.StatusMethodNotAllowed)

	// to delete the system generate headers
	// w.Header()["Date"] = nil

	// w.Header().Set("Allow", http.MethodPost)
	// app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper
	//
	// return
	// }

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.

	// Case ParseForm() which add any data in POST,PATCH,PUT request bodies to the r.PostForm map.

	// Limit the request body size to 4096 bytes
	// r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// The r.PostForm.Get() method always returns the form data as a string
	// However wer are expecting our 'expires' value to be a number, so we need to manually convert it
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("Content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// Initizlize a map to hold any validation errors for the form fields

	// Note: When we check the length of the title field, we’re using
	// the utf8.RuneCountInString() function — not Go’s len()
	// function. This is because we want to count the number of
	// characters in the title rather than the number of bytes. To
	// illustrate the difference, the string "Zoë" has 3 characters but a
	// length of 4 bytes because of the umlauted ë character.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// the ServeFile does not sanitize file path. if you are constrcting it from user data sanitize it first. filepath.clean()
	http.ServeFile(w, r, "./ui/static/file.zip")
}
