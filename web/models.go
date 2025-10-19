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

type ContentBlock struct {
	BlockType    string `json:"block_type"`
	Order        int    `json:"order"`
	TextContent  string `json:"text_content"`
	PhotoURL     string `json:"photo_url"`
	PhotoCaption string `json:"photo_caption"`
	YouTubeURL   string `json:"youtube_url"`
}

type BlogPostDetail struct {
	Title         string
	CreatedAt     time.Time
	ContentBlocks []ContentBlock
}

func (app *application) getBlogPostDetail(id int) (*BlogPostDetail, error) {
	query := `
		SELECT
			bp.title,
			bp.created_at,
			COALESCE(json_agg(
				jsonb_build_object(
					'block_type', cb.block_type,
					'order', cb.order,
					'text_content', COALESCE(cb.text_content, ''),
					'photo_url', COALESCE(ph.filename, ''),
					'photo_caption', COALESCE(ph.caption, ''),
					'youtube_url', COALESCE(cb.youtube_url, '')
				) ORDER BY cb.order
			) FILTER (WHERE cb.block_type IS NOT NULL), '[]') AS content_blocks
		FROM farmsville_blogpost bp
		LEFT JOIN farmsville_contentblock cb ON bp.id = cb.blog_post_id
		LEFT JOIN farmsville_photo ph ON cb.photo_id = ph.id
		WHERE bp.id = $1
		GROUP BY bp.id, bp.title, bp.created_at
	`

	var post BlogPostDetail
	var blocksJSON string

	err := app.db.QueryRow(query, id).Scan(
		&post.Title,
		&post.CreatedAt,
		&blocksJSON,
	)
	if err != nil {
		return nil, err
	}

	var blocks []ContentBlock
	json.Unmarshal([]byte(blocksJSON), &blocks)

	for i := range blocks {
		if blocks[i].PhotoURL != "" {
			blocks[i].PhotoURL = os.Getenv("PHOTOS_URL") + "/blog/" + blocks[i].PhotoURL
		}
	}

	post.ContentBlocks = blocks
	return &post, nil
}
