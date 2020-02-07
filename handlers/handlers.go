package handlers

import (
	"patches/models"

	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
)

type Env struct {
	DB models.Datastore
}

// Read all patches from database
func (env *Env)GetPatchesHandler(w http.ResponseWriter, r *http.Request) {
	patches, err := env.DB.GetPatches()
	
	if err != nil {
		log.Print("Error getting rows")
		log.Print(err)
	}
	// Return responses
	json.NewEncoder(w).Encode(patches)
}

// Create new patch in database
func (env *Env) PostPatchesHandler(w http.ResponseWriter, r *http.Request) {

	// Read http request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
	}
	patch := &models.Patch{}
	err = json.Unmarshal(body, &patch)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
		return
	}

	err = env.DB.CreatePatch(patch)

	// Send  response
	response := "Patch created"
	log.Print(response)
	json.NewEncoder(w).Encode(response)

}

// Delete patch from database
func (env *Env) DeletePatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Read http request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
	}
	convoID := &models.Convo_ID{}
	err = json.Unmarshal(body, &convoID)
	if err != nil {
		log.Print("Error reading request")
		log.Print(err)
		return
	}

	del, err := env.DB.DeletePatch(convoID.Convo_ID)

	// Return rows affected
	
	if del > 0 {
		response := fmt.Sprintf("%d rows deleted", del)
		json.NewEncoder(w).Encode(response)
	} else {
		response := "No rows deleted"
		json.NewEncoder(w).Encode(response)
	}
}