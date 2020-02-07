package handlers

import (
	"patches/models"

	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"time"

	"github.com/gorilla/schema"
)

type Env struct {
	DB models.Datastore
}

type Filter struct {
    Conversation   	[]string 		`schema:"convoID"`
    User 			[]string 		`schema:"userID"`
    Type  			[]string 		`schema:"type"`
	MaxTime     	time.Time 		`schema:"maxTime"`	//Newest
	MinTime     	time.Time 		`schema:"minTime"`	//Oldest
}

// Read all patches from database
func (env *Env)GetPatchesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
        // Handle error
    }

    filter := new(Filter)
    if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
        // Handle error
    }
	
	// Create filter string
	filterString := ""
	for _, convoID := range filter.Conversation {
        filterString = fmt.Sprintf("%s%s%s", filterString, " AND convo_id = ", convoID) 
	}
	for _, userID := range filter.User {
        filterString = fmt.Sprintf("%s%s%s", filterString, " AND user_id = ", userID) 
	}
	for _, patchType := range filter.Type {
        filterString = fmt.Sprintf("%s%s%s", filterString, " AND type = ", patchType) 
	}
	if (!filter.MaxTime.IsZero()) {
		filterString = fmt.Sprintf("%s%s%s", filterString, " AND time <= ", filter.MaxTime)
	}
	if (!filter.MinTime.IsZero()) {
		filterString = fmt.Sprintf("%s%s%s", filterString, " AND time >= ", filter.MinTime)
	}

	// Get patches with filters
	if(filterString != ""){
		filterString = filterString[4:]
	}
	filterString = fmt.Sprintf(" WHERE %s", filterString) 
	patches, err := env.DB.GetPatches(filterString)
	
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
	if err != nil {
		log.Print(err)
		json.NewEncoder(w).Encode(err)
	} else {
		response := "Patch created"
		log.Print(response)
		json.NewEncoder(w).Encode(response)
	}

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

	del, err := env.DB.DeletePatches(convoID.Convo_ID)

	// Return rows affected
	
	if del > 0 {
		response := fmt.Sprintf("%d rows deleted", del)
		json.NewEncoder(w).Encode(response)
	} else {
		response := "No rows deleted"
		json.NewEncoder(w).Encode(response)
	}
}