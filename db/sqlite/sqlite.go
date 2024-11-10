package sqlite

import (
	"database/sql"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func init() {
	sql.Register("sqlite3_with_fk",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				_, err := conn.Exec("PRAGMA foreign_keys = ON", nil)
				return err
			},
		})
}

func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite3_with_fk", path)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`
    CREATE TABLE IF NOT EXISTS experiments(
        id INTEGER PRIMARY KEY,
        public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
        name TEXT NOT NULL CHECK(name <> ''),
        description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `,
		`
    CREATE UNIQUE INDEX IF NOT EXISTS experiments_public_id ON
        experiments(public_id);
    `,
		`
	CREATE TABLE IF NOT EXISTS assessments (
		id INTEGER PRIMARY KEY,
		experiment_id INTEGER NOT NULL,
		name TEXT NOT NULL CHECK(name <> ''),
		description TEXT,
		is_pre BOOLEAN NOT NULL DEFAULT 0,  -- Indicates if it's the pre-assessment
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
	);
	`,
	}

	for _, q := range queries {
		_, err = db.Exec(q)

		if err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}
