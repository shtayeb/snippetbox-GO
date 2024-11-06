package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	// mux is http.Handler and it has ServerHttp() method so it satisfies the interface
	// the mux takes request and passes it to the necesssary handler based on route
	// You can think of a Go web application as a chain of ServeHTTP() methods being called one after another.
	//
	// Requests are handled concurrently - all http requests are served on their own go routines
	// This helps with speed but also creates `race condition` when accessing shared resources from handles
	// mux := http.NewServeMux()

	router := httprouter.New()

	// Custom error handler methods
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// relative to the project directory
	// To disable directory listing of fileserver.
	// 1- create an empty index.html file in the directory so that it can be fetched when a directory is requested
	// 2- create a custom implementation of http.FileSystem and have it return an os.ErrNotExist error for directories
	// http.Dir("./ui/static/")
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// /static/ - is subtree path. subtree paths end with /
	// /test - is redirected to /test/. if a subtree is registered
	// router.Handle("/static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// because secureHeaders is just a function, and the function returns a http.Handler
	// return secureHeaders(mux)
	// wraps the existing chain with the logRequest middleware
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// Alice middleware
	// Create a middleware chain containing our 'standard' middleware which will be used for every request our appliction receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// return the 'standard' middleware chain followed by the servemux
	return standard.Then(router)
}
