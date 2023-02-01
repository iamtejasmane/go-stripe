package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

const version = "1.0.0"
const ccsVersion = "1"

type config struct {
	port int
	env string
	api string
	db struct {
		dsn string
	}
	stripe struct {
		secret string
		key string
	}
}

type application struct {
	config config
	infolog *log.Logger
	errorlog *log.Logger
	templateCache map[string]*template.Template
	version string
}
func (app *application) serve() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: 30 * time.Second,
		ReadTimeout: 10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	app.infolog.Printf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port)

	return srv.ListenAndServe()
}
func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {developement|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc := make(map[string]*template.Template)

	app := &application{
		config: cfg,
		infolog: infoLog,
		errorlog: errorLog,
		templateCache: tc,
		version: version,
	}

	err := app.serve()
	if err != nil {
		app.errorlog.Println(err)
		log.Fatal(err)
	}

}