package main

import (
	"patches/models"
	
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"patches/handlers"
)

func main() {
	db, err := models.DBConnect()
	if err != nil {
		log.Fatal(err)
		return
	}

	env := &handlers.Env{db}

	httpMux := mux.NewRouter()

	httpMux.HandleFunc("/patches/v1/patches", env.PostPatchesHandler).Methods("POST")
	httpMux.HandleFunc("/patches/v1/patches", env.GetPatchesHandler).Methods("GET")
	httpMux.HandleFunc("/patches/v1/patches", env.DeletePatchesHandler).Methods("DELETE")

	httpSrv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      httpMux,
	}

	log.Fatal(httpSrv.ListenAndServe())
}