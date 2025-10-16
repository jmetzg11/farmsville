package main

import (
	"encoding/json"
	"os"
	"time"
)

type Product struct {
	ID          int
	EventID     int
	ProductName string
	Qty         int
	Remaining   int
	Notes       string
	PhotoURL    string
	PhotoType   string
}

type ProductClaimed struct {
	ID          int
	DateTime    time.Time
	User        string
	ProductName string `json:"product_name"`
	Qty         int
	Notes       string
}

func (app *application) getFutureProducts() ([]Product, []ProductClaimed, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT
			p.id,
			p.event_id,
			pn.name as product_name,
			p.qty,
			p.remaining,
			COALESCE(p.notes, '') as notes,
			COALESCE(ph.filename, '') as photo_url,
			COALESCE(ph.photo_type, '') as photo_type,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id', pc.id,
				'datetime', pc.datetime,
				'user', pc.user_name,
				'product_name', pn.name,
				'qty', pc.qty,
				'notes', COALESCE(pc.notes, '')
			)) FILTER (WHERE pc.id IS NOT NULL), '[]') as claims
		FROM farmsville_product p
		JOIN farmsville_event e ON p.event_id = e.id
		JOIN farmsville_productname pn ON p.product_name_id = pn.id
		LEFT JOIN farmsville_photo ph ON p.photo_id = ph.id
		LEFT JOIN farmsville_productclaimed pc ON p.id = pc.product_id
		WHERE e.date > $1
		GROUP BY p.id, p.event_id, pn.name, p.qty, p.remaining, p.notes, ph.filename, ph.photo_type, e.date
		ORDER BY e.date, pn.name
	`

	rows, err := app.db.Query(query, today)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var products []Product
	var allClaims []ProductClaimed

	for rows.Next() {
		var p Product
		var claimsJSON string

		err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.ProductName,
			&p.Qty,
			&p.Remaining,
			&p.Notes,
			&p.PhotoURL,
			&p.PhotoType,
			&claimsJSON,
		)
		if err != nil {
			return nil, nil, err
		}

		if p.PhotoURL != "" {
			subdir := "product"
			if p.PhotoType == "blog" {
				subdir = "blog"
			}
			p.PhotoURL = os.Getenv("PHOTOS_URL") + "/" + subdir + "/" + p.PhotoURL
		}

		products = append(products, p)

		var claims []ProductClaimed
		json.Unmarshal([]byte(claimsJSON), &claims)
		allClaims = append(allClaims, claims...)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return products, allClaims, nil
}

func (app *application) getProductRemaining(productID int) (int, error) {
	var remaining int
	query := `SELECT remaining FROM farmsville_product WHERE id = $1`
	err := app.db.QueryRow(query, productID).Scan(&remaining)
	if err != nil {
		return 0, err
	}
	return remaining, nil
}

func (app *application) createProductClaim(productID int, qty int, name string, notes string) error {
	tx, err := app.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get current remaining with row lock to prevent race conditions
	var remaining int
	query := `SELECT remaining FROM farmsville_product WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(query, productID).Scan(&remaining)
	if err != nil {
		return err
	}

	// Validate quantity
	if qty > remaining {
		return err
	}

	// Update remaining quantity
	updateQuery := `UPDATE farmsville_product SET remaining = remaining - $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, qty, productID)
	if err != nil {
		return err
	}

	// Insert claim record
	insertQuery := `
		INSERT INTO farmsville_productclaimed (product_id, qty, user_name, notes, datetime)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err = tx.Exec(insertQuery, productID, qty, name, notes)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type BlogPost struct {
	ID          int
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsPublished bool
}

func (app *application) getBlogPosts() ([]BlogPost, error) {
	query := `
		SELECT id, title, created_at
		FROM farmsville_blogpost
		WHERE is_published = true
		ORDER BY created_at DESC
	`

	rows, err := app.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
