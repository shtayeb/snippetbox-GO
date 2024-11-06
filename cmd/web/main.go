package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"github.com/shtayeb/snippetbox/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

// Closures for dependency injection
// https://gist.github.com/alexedwards/5cd712192b4831058b21

func main() {
	// Define a command line glag with name 'addr', a default value of ':4000'
	// The value returned from the flag.String() function is a pointer to the
	addr := flag.String("addr", ":4000", "HTTP network address")
	// database dsn
	dsn := flag.String("dsn", "postgres://go_user:go_1234@localhost/snippetbox", "Database source")
	flag.Parse()

	// Logging
	// Create a logger for logging infomation messages
	// distination, prefix, and other info
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)

	// Create a logger for error logging messages
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a decoder instance
	formDecoder := form.NewDecoder()

	// Initialize a new instance of our application struct. containing the dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// TO use our custom logger across our application we need to create a custom http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on port %s", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	err = srv.ListenAndServe()

	// you should avoid using the Panic() and Fatal() variations outside of your main() function
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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
