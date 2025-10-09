package main

import (
	"log"
	"net/http"
)

type homeData struct {
	Products []Product
	Claims   []ProductClaimed
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	products, claims, err := app.getFutureProducts()
	if err != nil {
		log.Printf("Error fetching future products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := homeData{
		Products: products,
		Claims:   claims,
	}

	app.render(w, http.StatusOK, "home.tmpl", data)
}
