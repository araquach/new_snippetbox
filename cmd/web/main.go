package main

import (
	"contra-design.com/new_snippetbox/pkg/models/mysql"
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	session *sessions.Session
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil { errorLog.Fatal(err)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed // before the main() function exits.
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:	errorLog,
		infoLog: 	infoLog,
		session: 	session,
		snippets:	&mysql.SnippetModel{DB:db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:		*addr,
		ErrorLog: 	errorLog,
		Handler:	app.routes(),
	}

	infoLog.Printf("Starting Server On %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err!= nil {
		return nil, err
	}
	return db, nil
}