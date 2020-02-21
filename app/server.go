package main

import (
	"patches/models"

	"log"
	"net/http"
	"patches/handlers"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var db models.Datastore

	env := handlers.NewEnv(db, &http.Client{Timeout: time.Second * 10})

	httpMux := mux.NewRouter()

	httpMux.HandleFunc("/patches/v1/patches", env.PostPatchesHandler).Methods("POST")
	httpMux.HandleFunc("/patches/v1/patches", env.GetPatchesHandler).Methods("GET")
	httpMux.HandleFunc("/patches/v1/patches", env.DeletePatchesHandler).Methods("DELETE")
	httpMux.HandleFunc("/patches/v1/connect/{conversation_id:[0-9]+}", env.ConnectHandler).Methods("GET")

	httpSrv := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      httpMux,
	}

	log.Fatal(httpSrv.ListenAndServe())
}
