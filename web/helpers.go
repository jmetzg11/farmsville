package main

import (
	"bytes"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
)

//go:embed ui/html ui/static
var Files embed.FS

func connectDB(prod bool) (*sql.DB, error) {
	var dsn string
	if prod {
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			return nil, fmt.Errorf("DATABASE_URL environment variable not set")
		}
		fmt.Println("Connecting to production database")
	} else {
		dsn = "postgresql://admin:admin@localhost:5432/farmsville?sslmode=disable"
		fmt.Println("Connection to local development database")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Database connection established")

	return db, nil
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pacific, _ := time.LoadLocation("America/Los_Angeles")
	funcMap := template.FuncMap{
		"percentage": func(remaining, total int) float64 {
			if total == 0 {
				return 0
			}
			return (float64(remaining) / float64(total)) * 100
		},
		"formatPacific": func(t time.Time) string {
			return t.In(pacific).Format("2006-01-02 15:04")
		},
	}

	pages, err := fs.Glob(Files, "ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"ui/html/base.html",
			"ui/html/claim_modal.html",
			page,
		}

		ts, err := template.New(name).Funcs(funcMap).ParseFS(Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func (app *application) render(w http.ResponseWriter, status int, page string, data any) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

var limiter = rate.NewLimiter(4, 10)

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for static assets
		if isStaticAsset(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isStaticAsset(path string) bool {
	staticPrefixes := []string{"/static/", "/data/photos/"}
	for _, prefix := range staticPrefixes {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

type ClaimNotification struct {
	ID          int
	ProductName string
	Qty         int
	UserName    string
	Notes       string
}

func (app *application) sendClaimEmail(claim ClaimNotification) error {
	gmailUser := os.Getenv("GMAIL_USER")
	gmailPass := os.Getenv("GMAIL_PASS")
	adminEmails := os.Getenv("ADMIN_EMAILS")

	if gmailUser == "" || gmailPass == "" {
		return fmt.Errorf("email credentials not configured")
	}

	if adminEmails == "" {
		return fmt.Errorf("admin emails not configured")
	}

	recipients := strings.Split(adminEmails, ",")
	for i := range recipients {
		recipients[i] = strings.TrimSpace(recipients[i])
	}

	subject := fmt.Sprintf("New Product Claim: %s", claim.ProductName)
	body := fmt.Sprintf(`A new product has been claimed:

Product: %s
Quantity: %d
Claimed by: %s
Notes: %s

`, claim.ProductName, claim.Qty, claim.UserName, claim.Notes)

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", gmailUser, strings.Join(recipients, ","), subject, body))

	auth := smtp.PlainAuth("", gmailUser, gmailPass, "smtp.gmail.com")

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		gmailUser,
		recipients,
		message,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
