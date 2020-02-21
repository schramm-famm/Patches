package handlers

import (
	"errors"
	"os"
	"patches/models"
	"patches/websockets"
	"strconv"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	gorillaws "github.com/gorilla/websocket"
)

type Env struct {
	DB       models.Datastore
	RC       *http.Client
	WSBroker *websockets.Broker
}

func NewEnv(db models.Datastore, rc *http.Client) *Env {
	return &Env{
		DB:       db,
		RC:       rc,
		WSBroker: websockets.NewBroker(db, rc),
	}
}

const (
	conversationContentRoute = "/ether/v1/converstions/%d/content"
)

var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var etherHost = os.Getenv("ETHER_HOST")

type Filter struct {
	Conversation []string  `schema:"convoID"`
	User         []string  `schema:"userID"`
	Type         []string  `schema:"type"`
	MaxTime      time.Time `schema:"maxTime"` //Newest
	MinTime      time.Time `schema:"minTime"` //Oldest
}

// Read all patches from database
func (env *Env) GetPatchesHandler(w http.ResponseWriter, r *http.Request) {
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
	if !filter.MaxTime.IsZero() {
		filterString = fmt.Sprintf("%s%s%s", filterString, " AND time <= ", filter.MaxTime)
	}
	if !filter.MinTime.IsZero() {
		filterString = fmt.Sprintf("%s%s%s", filterString, " AND time >= ", filter.MinTime)
	}

	// Get patches with filters
	if filterString != "" {
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

func (env *Env) getConversationContent(userID, conversationID int64) (string, error) {
	res, err := env.RC.Get("http://" + etherHost + fmt.Sprintf(conversationContentRoute, conversationID))
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusOK {
		return "", errors.New("Failed to get conversation content")
	} else if res.StatusCode == http.StatusNotFound {
		return "", errors.New("Conversation not found")
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (env *Env) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID, err := strconv.ParseInt(vars["conversation_id"], 10, 64)
	if err != nil {
		errMsg := "Invalid conversation ID"
		log.Println(errMsg + ": " + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade to WebSocket failed: ", err)
		return
	}

	fmt.Println("here?")
	client := env.WSBroker.CreateClient(conversationID, c)
	go client.Run()
}
