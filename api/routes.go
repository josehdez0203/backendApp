package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.Route("/api", func(mux chi.Router) {
		mux.Get("/", app.Hello)
		mux.Post("/user/", app.newUser)
		mux.Post("/authenticate", app.authenticate)
		mux.Post("/refresh", app.refreshToken)
		mux.Get("/logout", app.logout)
		// mux.Get("/movies", app.AllMovies)
		// mux.Get("/movies/{id}", app.GetMovie)
		// mux.Get("/genres", app.AllGenres)
		// mux.Get("/movies/genres/{id}", app.AllMoviesByGenre)
		mux.Route("/admin", func(r chi.Router) {
			r.Use(app.authRequired)
			// r.Get("/movies", app.MovieCatalog)
			// r.Get("/movies/{id}", app.MovieForEdit)
			// r.Put("/movies/0", app.InsertMovie)
			// r.Patch("/movies/{id}", app.UpdateMovie)
			// r.Delete("/movies/{id}", app.DeleteMovie)
		})
	})

	return mux
}
