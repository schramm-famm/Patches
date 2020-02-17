package models

import (
	"fmt"
	"log"
	"time"
)

type Patch struct {
	Timestamp time.Time `json:"timestamp"`
	Patch     string    `json:"patch"`
	Convo_ID  string    `json:"convo_id"`
	User_ID   string    `json:"user_id"`
	Type      string    `json:"type"`
}

type Convo_ID struct {
	Convo_ID int64 `json:"convo_id,string"`
}

// Get patch rows from database using filters
func (db *DB) GetPatches(filterString string) ([]Patch, error) {

	// Create query string with filters
	queryString := fmt.Sprintf("%s%s", "SELECT * FROM patches", filterString)
	rows, err := db.Query(queryString)
	if err != nil {
		log.Print("Error getting rows")
		log.Print(err)
		return nil, err
	}

	p := Patch{}
	patches := make([]Patch, 0)
	// Create patches object
	for rows.Next() {
		err := rows.Scan(&p.Timestamp, &p.Patch, &p.Convo_ID, &p.User_ID, &p.Type)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		patches = append(patches, p)
	}

	return patches, err
}

// Add new patch to the database
func (db *DB) CreatePatch(patch *Patch) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Insert patch datas into database
	_, err = tx.Exec("INSERT INTO patches(time,patch,convo_id,user_id,type) VALUES ($1, $2, $3, $4, $5) ", patch.Timestamp.Format(time.RFC3339), patch.Patch, patch.Convo_ID, patch.User_ID, patch.Type)
	if err != nil {
		log.Print("Error inserting")
		log.Print(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Delete patches from the database by conversation
func (db *DB) DeletePatches(convo_id int64) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// Delete patches from db
	del, err := db.Exec("DELETE FROM patches WHERE convo_id = $1", convo_id)
	if err != nil {
		log.Print("Error deleting")
		log.Print(err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	log.Print("Rows deleted")

	deleted, err := del.RowsAffected()
	if err != nil {
		log.Print("Error deleting rows")
		log.Print(err)
	}

	return deleted, err
}
