package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Config
type Config struct {
	Port int
}

// Application
type Application struct {
	Config   Config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// main is the main entry point of app
func main() {
	var config Config
	config.Port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		Config:   config,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	err := app.Serve()
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
