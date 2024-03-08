package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			OK      bool   `json:"ok"`
			Message string `json:"message"`
		}
		payload.OK = true
		payload.Message = "Hello, umon!!!"

		out, err := json.MarshalIndent(payload, "", "\t")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	})

	app.infoLog.Print("API lintening on port: ", app.config.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.port), nil)
}
