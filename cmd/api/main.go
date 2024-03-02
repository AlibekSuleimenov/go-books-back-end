package main

import (
	"fmt"
	"github.com/alibeksuleimenov/go-books-back-end/internal/data"
	"github.com/alibeksuleimenov/go-books-back-end/internal/driver"
	"log"
	"net/http"
	"os"
)

// Config is the type for app configuration
type Config struct {
	Port int
}

// Application is the type for sharing data across the app
type Application struct {
	Config      Config
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	Models      data.Models
	Environment string
}

// main is the main entry point of app
func main() {
	var config Config
	config.Port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")

	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()

	app := &Application{
		Config:      config,
		InfoLog:     infoLog,
		ErrorLog:    errorLog,
		Models:      data.New(db.SQL),
		Environment: environment,
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

// Serve starts the web server
func (app *Application) Serve() error {
	app.InfoLog.Println("API listening on port", app.Config.Port)

	serve := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: app.routes(),
	}

	return serve.ListenAndServe()
}
