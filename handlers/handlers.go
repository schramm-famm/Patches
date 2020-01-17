package handlers

import (
	"time"
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"io/ioutil"

    _ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Patch struct {
	Timestamp   time.Time      		 `json:"timestamp,omitempty"`
	Patch 		string 				 `json:"patch,omitempty"`
	Convo_ID    string               `json:"convo_id,omitempty"`
	User_ID     string         		 `json:"user_id,omitempty"`
	Type 		string 				 `json:"type,omitempty"`
}

// Connect to database
func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../patches.db")
	if err != nil {
		log.Fatal(err)
		return db, err
	}
	// test connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return db, err
	}
	return db, nil
}

// Create new patch in database
func PostPatchesHandler(w http.ResponseWriter, r *http.Request) {

	// Read http request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
	}
	patch := Patch{}
	err = json.Unmarshal(body, &patch)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
		return
	}

	db, err := dbConnect()
	defer db.Close()
	if err != nil {
		log.Print("Error connecting to database")
		log.Print(err)
		return
	}
	// Insert patch datas into database
	_, err = db.Exec("INSERT INTO patches(timestamp, patch, convo_id, user_id, type) VALUES (?,?,?,?,?) ", patch.Timestamp, patch.Patch, patch.Convo_ID, patch.User_ID, patch.Type);
	if err != nil {
		log.Print("Error inserting")
		log.Print(err)
		return
	}
	// Send  response
	response := "Patch created"
	json.NewEncoder(w).Encode(response)

}

// Read all patches from database
func GetPatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to db
	db, err := dbConnect()
	defer db.Close()
	if err != nil {
		log.Print("Error connecting to database")
		log.Print(err)
		return
	}
	// Read table
	rows, err := db.Query("SELECT * FROM patches")
	if err != nil {
		log.Print("Error getting rows")
		log.Print(err)
		return
	}

	p := Patch{}
	patches := make([]Patch, 0)
	// Create patches object
	for rows.Next() {
		err := rows.Scan(&p.Timestamp, &p.Patch, &p.Convo_ID, &p.User_ID, &p.Type)
		if err != nil {
			log.Print(err)
		}
		patches = append(patches, p)
	}
	// Return responses
	json.NewEncoder(w).Encode(patches)
}

// Update patch in databsse
func PutPatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Read http request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
	}
	patch := Patch{}
	err = json.Unmarshal(body, &patch)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
		return
	}
	// Connect to db
	db, err := dbConnect()
	defer db.Close()
	if err != nil {
		log.Print("Error connecting to database")
		log.Print(err)
		return
	}
	// Update patch in db
	update, err := db.Exec("UPDATE patches SET patch = ?, convo_id = ?, user_id = ?, type = ? WHERE timestamp = ?;",  patch.Patch, patch.Convo_ID, patch.User_ID, patch.Type, patch.Timestamp)
	if err != nil {
		log.Print("Error updating")
		log.Print(err)
	}
	//db.Commit()
	// Return number of rows affected
	updated, err := update.RowsAffected()
	if err != nil {
		log.Print("Error updatin rows")
		log.Print(err)
	}
	if updated > 0 {
		response := "Rows updated"
		json.NewEncoder(w).Encode(response)
	} else {
		response := "No rows updated"
		json.NewEncoder(w).Encode(response)
	}

}

// Delete patch from database
func DeletePatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Read http request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
	}
	timestamp := time.Time{}
	err = json.Unmarshal(body, &timestamp)
	if err != nil {
		log.Print("Error converting request")
		log.Print(err)
		return
	}
	// Connect to db
	db, err := dbConnect()
	defer db.Close()
	if err != nil {
		log.Print("Error connecting to database")
		log.Print(err)
		return
	}
	// Delete patches from db
	del, err := db.Exec("DELETE FROM patches WHERE timestamp = ?", timestamp)
	if err != nil {
		log.Print("Error deleting")
		log.Print(err)
	}
	//err = db.Commit()
	// Return rows affected
	deleted, err := del.RowsAffected()
	if err != nil {
		log.Print("Error deleting rows")
		log.Print(err)
	}
	if deleted > 0 {
		response := "Rows deleted"
		json.NewEncoder(w).Encode(response)
	} else {
		response := "No rows deleted"
		json.NewEncoder(w).Encode(response)
	}
}