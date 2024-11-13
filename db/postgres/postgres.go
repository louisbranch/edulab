package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(connection string) (*DB, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS experiments (
			id SERIAL PRIMARY KEY,
			public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
			name TEXT NOT NULL CHECK(name <> ''),
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS experiments_public_id ON experiments(public_id);
		`,
		`
		CREATE TABLE IF NOT EXISTS assessments (
			id SERIAL PRIMARY KEY,
			experiment_id INTEGER NOT NULL,
			public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
			type TEXT CHECK(type IN ('pre', 'post')),
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS assessments_public_id ON assessments(public_id);
		`,
		`
		CREATE TABLE IF NOT EXISTS questions (
			id SERIAL PRIMARY KEY,
			assessment_id INTEGER NOT NULL,
			text TEXT NOT NULL CHECK(text <> ''),
			type TEXT CHECK(type IN ('multiple', 'single', 'text')),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (assessment_id) REFERENCES assessments(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS question_choices (
			id SERIAL PRIMARY KEY,
			question_id INTEGER NOT NULL,
			text TEXT NOT NULL CHECK(text <> ''),
			is_correct BOOLEAN NOT NULL DEFAULT FALSE,
			FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS cohorts (
			id SERIAL PRIMARY KEY,
			experiment_id INTEGER NOT NULL,
			public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
			name TEXT NOT NULL CHECK(name <> ''),
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS cohorts_public_id ON cohorts(public_id);
		`,
		`
		CREATE TABLE IF NOT EXISTS demographics (
			id SERIAL PRIMARY KEY,
			experiment_id INTEGER NOT NULL,
			text TEXT NOT NULL CHECK(text <> ''),
			type TEXT CHECK(type IN ('multiple', 'single', 'text')),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS demographic_options (
			id SERIAL PRIMARY KEY,
			demographic_id INTEGER NOT NULL,
			text TEXT NOT NULL CHECK(text <> ''),
			FOREIGN KEY (demographic_id) REFERENCES demographics(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS participants (
			id SERIAL PRIMARY KEY,
			public_id TEXT NOT NULL UNIQUE CHECK(public_id <> ''),
			experiment_id INTEGER NOT NULL,
			cohort_id INTEGER NOT NULL,
			access_token TEXT NOT NULL UNIQUE CHECK(access_token <> ''),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE,
			FOREIGN KEY (cohort_id) REFERENCES cohorts(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS participations (
			experiment_id INTEGER NOT NULL,
			assessment_id INTEGER NOT NULL,
			participant_id INTEGER NOT NULL,
			answers TEXT,
			demographics TEXT,
			FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE,
			FOREIGN KEY (assessment_id) REFERENCES assessments(id) ON DELETE CASCADE,
			FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE
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
