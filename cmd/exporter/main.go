package main

import (
	"fmt"
	"log"
	"os"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/db/postgres"
	"github.com/louisbranch/edulab/db/sqlite"
	"github.com/louisbranch/edulab/result"
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

	res, err := result.New(db, "1")
	if err != nil {
		log.Fatal(err)
	}

	aq := []result.AssessmentQuestions{
		{AssessmentID: "1", QuestionID: "2"},
		{AssessmentID: "2", QuestionID: "7"},
	}

	cmp, err := result.NewComparison(res, aq, []string{"1", "2"})
	if err != nil {
		log.Fatal(err)
	}

	err = cmp.ToCSV("comparison.csv")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Comparison data exported to comparison.csv")
}
