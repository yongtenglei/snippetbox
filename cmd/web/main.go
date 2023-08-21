package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"rey.com/snippetbox/internal/models"
)

type application struct {
	debug bool

	infoLog *log.Logger
	errLog  *log.Logger

	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "rey:rey@/snippetbox?parseTime=true", "MySQL data source name")
	debug := flag.Bool("debug", false, "Enable debug mode")
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := application{
		debug: *debug,

		infoLog: infoLog,
		errLog:  errLog,

		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		// if you allow https request and set minimum version to TLS 1.3
		// Then you do not need any other additional mitigation against
		// CSRF attacks (like justinas/nosurf) besides alexedwards/scs
		// MinVersion:       tls.VersionTLS13,
	}

	svr := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  app.errLog,
		TLSConfig: tlsConfig,

		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		// It's sensible to set WriteTimeout Greater than ReadTimeout
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %s", *addr)
	err = svr.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
