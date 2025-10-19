package main

import (
	"fmt"
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create claim"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Claim submitted successfully"))
}

type blogData struct {
	Posts []BlogPost
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	posts, err := app.getBlogPosts()
	if err != nil {
		log.Printf("Error fetching blog posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := blogData{
		Posts: posts,
	}

	app.render(w, http.StatusOK, "blog.html", data)
}

func (app *application) blogDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid blog post ID", http.StatusBadRequest)
		return
	}

	post, err := app.getBlogPostDetail(id)
	if err != nil {
		log.Printf("Error fetching blog post detail: %v", err)
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}
	fmt.Println(post)

	app.render(w, http.StatusOK, "blog_detail.html", post)
}
