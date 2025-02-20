package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/josehdez0203/backendApp/internal/repository"
	"github.com/josehdez0203/backendApp/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	Domain       string
	DNS          string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	// aplication config
	var app application
	// read from command line
	flag.StringVar(&app.DNS, "dns", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api-key", "883bdf20a6a816f337a8b3c57aeabf17", "api key")
	flag.Parse()
	// Connect to the database
	conn, err := app.connectDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()
	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Hour * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Running ⭐ ", port)
	// Start a web server

	err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
