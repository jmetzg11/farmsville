package main

import (
	"bytes"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
)

//go:embed "ui/html ui/static"
var Files embed.FS 

func connectDB(prod bool) (*sql.DB, error) {
	var dsn string 
	if prod {
		dsn = os.Getenve("DATABASE_URL")
		if dsn == "" {
			return nil fmt.Errorf("kk")
		}
	}
}
