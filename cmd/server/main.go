package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/db/postgres"
	"github.com/louisbranch/edulab/db/sqlite"
	"github.com/louisbranch/edulab/web/html"
	"github.com/louisbranch/edulab/web/server"
	"github.com/louisbranch/edulab/wizard"
)

func main() {

	dev := true
	if os.Getenv("APP_ENV") == "production" {
		dev = false
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	files := os.Getenv("FILES_PATH")
	if files == "" {
		files = "web"
	}

	var db edulab.Database
	var err error

	dburl := os.Getenv("DATABASE_URL")
	dbuser := os.Getenv("POSTGRES_USER")

	if dburl != "" {
		log.Println("using database url")
		db, err = postgres.New(dburl)
	} else if dbuser == "" {
		log.Println("using sqlite database")
		db, err = sqlite.New("edulab.db")
	} else {
		log.Println("using postgres database")
		pswd := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOSTNAME")
		dbname := os.Getenv("POSTGRES_DB")

		sslmode := "verify-full"
		if dev {
			sslmode = "disable"
		}

		connection := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
			dbuser, pswd, host, dbname, sslmode)
		log.Printf("connection: %s\n", connection)
		db, err = postgres.New(connection)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = wizard.Experiment(db)
	if err != nil {
		log.Fatal(err)
	}

	srv := &server.Server{
		DB:       db,
		Template: html.New(filepath.Join(files, "templates")),
		Assets:   http.FileServer(http.Dir(filepath.Join(files, "assets"))),
		Random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	mux := srv.NewServeMux()

	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
