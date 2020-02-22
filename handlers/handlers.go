package handlers

import (
	"patches/models"
	"patches/websockets"
	"strconv"

	"encoding/json"
	"log"
	"net/http"

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

var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

	go env.WSBroker.StartClient(conversationID, c)
}
