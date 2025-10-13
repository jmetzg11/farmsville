package main

import (
	"io/fs"
	"net/http"
)

func (app *application) routes(prod bool) http.Handler {
	mux := http.NewServeMux()

	staticFS, _ := fs.Sub(Files, "ui/static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Serve photos locally in development only
	if !prod {
		photosFS := http.Dir("../data/photos")
		mux.Handle("GET /data/photos/", http.StripPrefix("/data/photos/", http.FileServer(photosFS)))
	}

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("POST /claim", app.claimProduct)

	return rateLimit(mux)
}
