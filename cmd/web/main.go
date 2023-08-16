package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"rey.com/snippetbox/internal/models"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger

	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "rey:rey@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	// Initialize template catch
	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := application{
		infoLog: infoLog,
		errLog:  errLog,

		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	svr := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: app.errLog,
	}

	app.infoLog.Printf("Starting server on %s", *addr)
	err = svr.ListenAndServe()
	app.errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
