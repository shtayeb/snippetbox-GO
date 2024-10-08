package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
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
		log.Print(err.Error())
		http.Error(w, "Internal Serve Error 1", 500)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Serve Error 2", 500)
	}

}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)

		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// w.WriteHeader can be called once per request and after the status code is set it cant be changed
		// if you dont call w.WriteHeader explicitly, then the first call to w.Write will automatically send a 200 ok
		// so if you want to send a non-200 code you must call w.WriteHeader before any call to w.Write

		// This
		w.WriteHeader(405)
		w.Header().Set("Allow", "POST")
		w.Write([]byte("Method not allowed"))

		// Or helper

		// http.StatusMethodNotAllowed = 405
		http.Error(w, "Method not allowed", 405)

		// to delete the system generate headers
		w.Header()["Date"] = nil

		return
	}

	w.Write([]byte("Create a new snippet"))
}
