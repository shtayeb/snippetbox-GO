package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shtayeb/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// file in the slice
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }
	//
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	http.Error(w, "Internal Serve Error 1", 500)
	//
	// 	return
	// }
	//
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	http.Error(w, "Internal Serve Error 2", 500)
	// }

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper

		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snai Climb Mount Fuji"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// the ServeFile does not sanitize file path. if you are constrcting it from user data sanitize it first. filepath.clean()
	http.ServeFile(w, r, "./ui/static/file.zip")
}
