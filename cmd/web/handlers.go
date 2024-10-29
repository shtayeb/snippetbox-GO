package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	// file in the slice
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Serve Error 1", 500)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal Serve Error 2", 500)
	}

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)

		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
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

	w.Write([]byte("Create a new snippet"))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// the ServeFile does not sanitize file path. if you are constrcting it from user data sanitize it first. filepath.clean()
	http.ServeFile(w, r, "./ui/static/file.zip")
}
