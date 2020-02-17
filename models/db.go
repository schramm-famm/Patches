package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "vagrant"
	dbname   = "patches"
)

type Datastore interface {
	CreatePatch(patch *Patch) error
	GetPatches(filterString string) ([]Patch, error)
	DeletePatches(convo_id int64) (int64, error)
}

// DB represents an SQL database connection
type DB struct {
	*sql.DB
}

// Connect to database
func DBConnect() (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
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
