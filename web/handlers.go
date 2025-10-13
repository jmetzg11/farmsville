package main

import (
	"log"
	"net/http"
	"strconv"
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

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) claimProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse and validate product ID
	productIDStr := r.FormValue("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid product ID"))
		return
	}

	// Parse and validate quantity
	qtyStr := r.FormValue("qty")
	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Quantity must be a positive number"))
		return
	}

	// Validate name
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name is required"))
		return
	}

	// Get current remaining from database (don't trust client data)
	remaining, err := app.getProductRemaining(productID)
	if err != nil {
		log.Printf("Error fetching product remaining: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Validate quantity against actual remaining
	if qty > remaining {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Quantity must be at most " + strconv.Itoa(remaining)))
		return
	}

	// Create the claim (this will update the database)
	notes := r.FormValue("notes")
	err = app.createProductClaim(productID, qty, name, notes)
	if err != nil {
		log.Printf("Error creating product claim: %v", err)
		http.Error(w, "Failed to create claim", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Claim submitted successfully"))
}
