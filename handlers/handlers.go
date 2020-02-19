package handlers

import (
	"patches/models"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

type Env struct {
	DB models.Datastore
}

// Get patches from the database with filtering
func (env *Env) GetPatchesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print("Error reading request")
		log.Print(err)
		return
	}

	filter := new(models.Filter)
	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		log.Print("Error reading request")
		log.Print(err)
		return
	}

	// Get patches with filter
	patches, err := env.DB.GetPatches(filter)

	if err != nil {
		log.Print("Error getting rows")
		log.Print(err)
		return
	}
	// Return responses
	json.NewEncoder(w).Encode(patches)
}
