package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl", nil)
}
