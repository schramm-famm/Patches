package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"patches/handlers"
	"patches/models"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		os.Getenv("PATCHES_DB_HOST"),
		os.Getenv("PATCHES_DB_PORT"),
		os.Getenv("PATCHES_DB_USERNAME"),
		os.Getenv("PATCHES_DB_PASSWORD"),
		os.Getenv("PATCHES_DB_DATABASE"))
	db, err := models.DBConnect(connectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	env := handlers.NewEnv(db, &http.Client{Timeout: time.Second * 10})

	httpMux := mux.NewRouter()

	httpMux.HandleFunc("/patches/v1/patches", env.GetPatchesHandler).Methods("GET")
	httpMux.HandleFunc("/patches/v1/connect/{conversation_id:[0-9]+}", env.ConnectHandler).Methods("GET")

	httpSrv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      httpMux,
	}

	log.Fatal(httpSrv.ListenAndServe())
}
