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
        public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
		type TEXT CHECK(type IN ('pre_assessment', 'post_assessment')),
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
	);
	`,
		`
    CREATE UNIQUE INDEX IF NOT EXISTS assessments_public_id ON
        assessments(public_id);
    `,
		`
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY,
		assessment_id INTEGER NOT NULL,
		prompt TEXT NOT NULL CHECK(prompt <> ''),
		type TEXT CHECK(type IN ('multiple_choice', 'single_choice', 'text')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (assessment_id) REFERENCES assessments(id) ON DELETE CASCADE
	);`,
		`
	CREATE TABLE IF NOT EXISTS question_choices (
		id INTEGER PRIMARY KEY,
		question_id INTEGER NOT NULL,
		text TEXT NOT NULL CHECK(text <> ''),
		is_correct BOOLEAN NOT NULL DEFAULT 0,
		FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
	);`,
		`
	CREATE TABLE IF NOT EXISTS cohorts (
		id INTEGER PRIMARY KEY,
		experiment_id INTEGER NOT NULL,
        public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
		name TEXT NOT NULL CHECK(name <> ''),
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
	);
	`,
		`
    CREATE UNIQUE INDEX IF NOT EXISTS cohorts_public_id ON
        cohorts(public_id);
    `,
		`
	CREATE TABLE IF NOT EXISTS demographics (
		id INTEGER PRIMARY KEY,
		experiment_id INTEGER NOT NULL,
		i18n_key TEXT,  -- Localization key (e.g., "age", "gender")
		text TEXT,  -- Direct text for non-localized options
		type TEXT CHECK(type IN ('multiple_choice', 'single_choice', 'text')),
		FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE,
		CHECK (i18n_key IS NOT NULL OR text IS NOT NULL)  -- Ensures at least one is present
	);`,
		`
	CREATE TABLE IF NOT EXISTS demographic_options (
		id INTEGER PRIMARY KEY,
		demographic_id INTEGER NOT NULL,
		i18n_key TEXT,  -- Localization key (e.g., "gender_male")
		text TEXT,  -- Direct text for non-localized options
		FOREIGN KEY (demographic_id) REFERENCES demographics(id) ON DELETE CASCADE,
		CHECK (i18n_key IS NOT NULL OR text IS NOT NULL)  -- Ensures at least one is present
	);`,
	}

	for _, q := range queries {
		_, err = db.Exec(q)

		if err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}
