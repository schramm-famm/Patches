package models

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Patch struct {
	Timestamp time.Time `json:"timestamp"`
	Patch     string    `json:"patch"`
	ConvoID   string    `json:"convo_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
}

type Filter struct {
	Conversation string    `schema:"convo_id"`
	User         []string  `schema:"user_id"`
	Type         []string  `schema:"type"`
	EndTime      time.Time `schema:"end_time"`
	StartTime    time.Time `schema:"start_time"`
}

// Get patch rows from database using filters
func (db *DB) GetPatches(filter *Filter) ([]Patch, error) {

	// Create filter string
	filterString := new(strings.Builder)
	if filter.Conversation != "" {
		fmt.Fprintf(filterString, "convo_id = %s", filter.Conversation)
	} else {
		log.Print("Error with conversation id")
		return nil, nil
	}

	if len(filter.User) >= 1 {
		fmt.Fprintf(filterString, " AND (")

		for i, userID := range filter.User {
			if i == 0 {
				fmt.Fprintf(filterString, "user_id = %s", userID)
			} else {
				fmt.Fprintf(filterString, " OR user_id = %s", userID)
			}
		}

		fmt.Fprintf(filterString, ")")
	}

	if len(filter.Type) >= 1 {
		fmt.Fprintf(filterString, " AND (")

		for i, patchType := range filter.Type {
			if i == 0 {
				fmt.Fprintf(filterString, "type = %s", patchType)
			} else {
				fmt.Fprintf(filterString, " OR type = %s", patchType)
			}
		}

		fmt.Fprintf(filterString, ")")
	}

	if !filter.EndTime.IsZero() {
		fmt.Fprintf(filterString, " AND time <= %s", filter.EndTime)
	}
	if !filter.StartTime.IsZero() {
		fmt.Fprintf(filterString, " AND time >= %s", filter.StartTime)
	}

	// Create query string with filters
	queryString := new(strings.Builder)
	fmt.Fprintf(queryString, "SELECT * FROM patches WHERE %s", filterString)
	rows, err := db.Query(queryString.String())
	defer rows.Close()
	if err != nil {
		log.Print("Error getting rows")
		log.Print(err)
		return nil, err
	}

	p := Patch{}
	patches := make([]Patch, 0)
	// Create patches object
	for rows.Next() {
		err := rows.Scan(&p.Timestamp, &p.Patch, &p.ConvoID, &p.UserID, &p.Type)
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

	// Insert patch into database
	_, err := db.Exec("INSERT INTO patches(time,patch,convo_id,user_id,type) VALUES ($1, $2, $3, $4, $5) ", patch.Timestamp.Format(time.RFC3339), patch.Patch, patch.ConvoID, patch.UserID, patch.Type)
	if err != nil {
		log.Print("Error inserting")
		log.Print(err)
		return err
	}

	return nil
}

// Delete patches from the database by conversation
func (db *DB) DeletePatches(convo_id int64) (int64, error) {

	// Delete patches from db
	del, err := db.Exec("DELETE FROM patches WHERE convo_id = $1", convo_id)
	if err != nil {
		log.Print("Error deleting")
		log.Print(err)
		return 0, err
	}

	log.Print("Rows deleted")

	deleted, err := del.RowsAffected()
	if err != nil {
		log.Print("Error deleting rows")
		log.Print(err)
		return 0, err
	}

	return deleted, err
}
