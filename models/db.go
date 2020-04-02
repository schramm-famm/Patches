package models

import (
	"database/sql"
	"io/ioutil"
	"strings"

	"github.com/lib/pq"
)

// Datastore defines the CRUD operations of patches in the database
type Datastore interface {
	CreatePatch(patch *Patch) error
	GetPatches(filter *Filter) ([]Patch, error)
	DeletePatches(convo_id int64) (int64, error)
}

// DB represents an SQL database connection
type DB struct {
	*sql.DB
}

// DBConnect initializes a new DB
func DBConnect(connectionString string) (*DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec("CREATE DATABASE patches"); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code != pq.ErrorCode("42P04") {
			return nil, err
		}
	}
	db.Close()

	db, err = sql.Open("postgres", connectionString+" dbname=patches")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := setupTable(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// setupTable creates the necessary "patches" table if it doesn't already exist.
func setupTable(db *sql.DB) error {
	var tmp string

	// Check if the "patches" table exists by querying for the table. If it does
	// exist, just return.
	queryTable := "SELECT table_name FROM information_schema.tables WHERE table_name='patches'"
	err := db.QueryRow(queryTable).Scan(&tmp)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != sql.ErrNoRows {
		return nil
	}

	sqlScript, err := ioutil.ReadFile("dbSchema.sql")
	if err != nil {
		return err
	}

	statements := strings.Split(string(sqlScript), ";")
	if len(statements) > 0 {
		statements = statements[:len(statements)-1]
	}

	for _, statement := range statements {
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
	}
	return nil
}
