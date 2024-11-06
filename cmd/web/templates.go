package main

import (
	"html/template"
	"path/filepath"

	"github.com/shtayeb/snippetbox/internal/models"
	"time"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

// custom template functions (like
// our humanDate() function) can accept as many parameters as they
// need to, but they must return one value only. The only exception to
// this is if you want to return an error as the second value, in which
// case that’s OK too.
func humanDate(t time.Time) string {
	return t.Format("03 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// filepath.Glob() function to get a slice of all filepaths that match the pattern "./ui/html/pages/*.tmpl
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// loop through pages
	for _, page := range pages {
		// extract the filename (like 'home.tmpl') from the full path
		name := filepath.Base(page)

		// the template template.FuncMap must be registered with the template set

		// Parse the base template file into a template set
		// ts, err := template.ParseFiles("./ui/html/base.tmpl")
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// parse the files into a template set.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page (like 'home.tmpl') as the key
		cache[name] = ts
	}

	// return the cache map
	return cache, nil
}
