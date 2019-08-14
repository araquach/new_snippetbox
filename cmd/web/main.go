package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "host=localhost  port=5432  user=adam-macbook  password=  dbname=postgres  sslmode=disable", "")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil { errorLog.Fatal(err)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed // before the main() function exits.
	defer db.Close()


	app := &application{
		errorLog:	errorLog,
		infoLog: 	infoLog,
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
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err!= nil {
		return nil, err
	}
	return db, nil
}