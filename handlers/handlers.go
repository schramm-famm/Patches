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
		errMsg := "Error reading request:" + err.Error()
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	filter := new(models.Filter)
	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		errMsg := "Error reading request:" + err.Error()
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Get patches with filter
	patches, err := env.DB.GetPatches(filter)

	if err != nil {
		errMsg := "Error getting rows:" + err.Error()
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	} else if patches == nil {
		errMsg := "Error getting rows"
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Return responses
	log.Printf("%d patches returned", len(patches))
	json.NewEncoder(w).Encode(patches)
}
