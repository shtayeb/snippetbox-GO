package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

// Closures for dependency injection
// https://gist.github.com/alexedwards/5cd712192b4831058b21

func main() {
	// Define a command line glag with name 'addr', a default value of ':4000'
	// The value returned from the flag.String() function is a pointer to the
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Logging
	// Create a logger for logging infomation messages
	// distination, prefix, and other info
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)

	// Create a logger for error logging messages
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of our application struct. containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// mux is http.Handler and it has ServerHttp() method so it satisfies the interface
	// the mux takes request and passes it to the necesssary handler based on route
	// You can think of a Go web application as a chain of ServeHTTP() methods being called one after another.
	//
	// Requests are handled concurrently - all http requests are served on their own go routines
	// This helps with speed but also creates `race condition` when accessing shared resources from handles
	mux := http.NewServeMux()

	// relative to the project directory
	// To disable directory listing of fileserver.
	// 1- create an empty index.html file in the directory so that it can be fetched when a directory is requested
	// 2- create a custom implementation of http.FileSystem and have it return an os.ErrNotExist error for directories
	// http.Dir("./ui/static/")
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())

	// /static/ - is subtree path. subtree paths end with /
	// /test - is redirected to /test/. if a subtree is registered
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// TO use our custom logger across our application we need to create a custom http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on port %s", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()

	// you should avoid using the Panic() and Fatal() variations outside of your main() function
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
