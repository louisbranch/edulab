package sqlite

import (
	"log"
	"os"
	"testing"

	"github.com/louisbranch/edulab"
)

func TestDBInterface(t *testing.T) {
	var _ edulab.Database = &DB{}

}

func testDB() (*DB, string) {
	tmpfile, err := os.CreateTemp("", "edulab.db")
	if err != nil {
		log.Fatal(err)
	}
	name := tmpfile.Name()
	db, err := New(name)
	if err != nil {
		log.Fatal(err)
	}
	return db, name
}