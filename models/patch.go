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
	ConvoID   int64     `json:"convo_id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"type"`
}

type Filter struct {
	Conversation int64     `schema:"convo_id"`
	User         []int64   `schema:"user_id"`
	Type         []string  `schema:"type"`
	EndTime      time.Time `schema:"end_time"`
	StartTime    time.Time `schema:"start_time"`
}

// GetPatches gets patch rows from the database using filters
func (db *DB) GetPatches(filter *Filter) ([]Patch, error) {

	// Create filter string
	filterString := new(strings.Builder)
	if filter.Conversation != 0 {
		fmt.Fprintf(filterString, "convo_id = %d", filter.Conversation)
	} else {
		log.Print("Error with conversation id")
		return nil, nil
	}

	if len(filter.User) >= 1 {
		fmt.Fprintf(filterString, " AND (")

		for i, userID := range filter.User {
			if i == 0 {
				fmt.Fprintf(filterString, "user_id = %d", userID)
			} else {
				fmt.Fprintf(filterString, " OR user_id = %d", userID)
			}
		}

		fmt.Fprintf(filterString, ")")
	}

	if len(filter.Type) >= 1 {
		fmt.Fprintf(filterString, " AND (")

		for i, patchType := range filter.Type {
			if i == 0 {
				fmt.Fprintf(filterString, "type = '%s'", patchType)
			} else {
				fmt.Fprintf(filterString, " OR type = '%s'", patchType)
			}
		}

		fmt.Fprintf(filterString, ")")
	}

	if !filter.EndTime.IsZero() {
		t := filter.EndTime.Format("2006-01-02 15:04:05")
		fmt.Fprintf(filterString, " AND time <= '%s'", t)
	}
	if !filter.StartTime.IsZero() {
		t := filter.StartTime.Format("2006-01-02 15:04:05")
		fmt.Fprintf(filterString, " AND time >= '%s'", t)
	}

	// Create query string with filters
	queryString := new(strings.Builder)
	fmt.Fprintf(queryString, "SELECT * FROM patches WHERE %s", filterString)
	rows, err := db.Query(queryString.String())
	if err != nil {
		log.Print("Error getting rows")
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

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

// CreatePatch adds a new patch to the database
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

// DeletePatches deletes patches from the database by conversation
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
