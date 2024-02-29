package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port int
}

type Application struct {
	Config   Config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

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

func (app *Application) Serve() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			OK      bool   `json:"ok"`
			Message string `json:"message"`
		}
		payload.OK = true
		payload.Message = "Welcome"

		out, err := json.MarshalIndent(payload, "", "\t")
		if err != nil {
			app.ErrorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	})

	app.InfoLog.Println("API listening on port", app.Config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
}
