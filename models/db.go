package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Datastore interface {
	CreatePatch(patch *Patch) error
	GetPatches(filter *Filter) ([]Patch, error)
	DeletePatches(convo_id int64) (int64, error)
}

// DB represents an SQL database connection
type DB struct {
	*sql.DB
}

// Connect to database
func DBConnect(connectionString string) (*DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &DB{db}, nil
}
