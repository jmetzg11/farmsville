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
}

type ProductClaimed struct {
	ID          int
	DateTime    time.Time
	User        string
	ProductName string
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
			COALESCE(p.photo_url, '') as photo_url,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id', pc.id,
				'datetime', pc.datetime,
				'user', pc.user,
				'product_name', pn.name,
				'qty', pc.qty,
				'notes', COALESCE(pc.notes, '')
			)) FILTER (WHERE pc.id IS NOT NULL), '[]') as claims
		FROM farmsville_product p
		JOIN farmsville_event e ON p.event_id = e.id
		JOIN farmsville_productname pn ON p.product_name_id = pn.id
		LEFT JOIN farmsville_productclaimed pc ON p.id = pc.product_id
		WHERE e.date > $1
		GROUP BY p.id, p.event_id, pn.name, p.qty, p.remaining, p.notes, p.photo_url, e.date
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
			&claimsJSON,
		)
		if err != nil {
			return nil, nil, err
		}

		if p.PhotoURL != "" {
			p.PhotoURL = os.Getenv("PHOTOS_URL") + p.PhotoURL
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
