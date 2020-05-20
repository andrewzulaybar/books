package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func publicationsHandler(db *postgres.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			publications, err := publication.GetMany(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			bytes, err := json.Marshal(publications)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		case http.MethodPost:
			publication, err := publication.PostOne(db, r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			bytes, err := json.Marshal(publication)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(bytes)
		case http.MethodDelete:
			if err := publication.DeleteMany(db, r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	})
}

func publicationHandler(db *postgres.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		switch r.Method {
		case http.MethodGet:
			publication, err := publication.GetOne(db, ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			bytes, err := json.Marshal(publication)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		case http.MethodPatch:
			if err := publication.PatchOne(db, r.Body, ID); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			if err := publication.DeleteOne(db, ID); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	})
}

func main() {
	conf, err := config.Load("config/.env")
	if err != nil {
		panic(err)
	}

	db, err := postgres.Connect(conf.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer postgres.Disconnect(db)

	if err := postgres.Init(db); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	API := r.PathPrefix("/api").Subrouter()

	API.HandleFunc("/publications", publicationsHandler(db)).
		Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	API.HandleFunc("/publications/{id:[0-9]+}", publicationHandler(db)).
		Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)

	srv := &http.Server{
		Handler:      handlers.CORS()(r),
		Addr:         conf.Address,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
