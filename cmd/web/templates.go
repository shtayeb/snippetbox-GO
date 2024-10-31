package main

import (
	"html/template"
	"path/filepath"

	"github.com/shtayeb/snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
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

		// Parse the base template file into a template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
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
