package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"patches/handlers"
)

func main() {
	httpMux := mux.NewRouter()

	httpMux.HandleFunc("/api/patches", handlers.PostPatchesHandler).Methods("POST")
	httpMux.HandleFunc("/api/patches", handlers.GetPatchesHandler).Methods("GET")
	httpMux.HandleFunc("/api/patches", handlers.PutPatchesHandler).Methods("PATCH")
	httpMux.HandleFunc("/api/patches", handlers.DeletePatchesHandler).Methods("DELETE")

	httpSrv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      httpMux,
	}

	log.Fatal(httpSrv.ListenAndServe())
}