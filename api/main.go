package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/publication"
)

func publicationsHandler(db *postgres.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			publications := publication.Get(db)
			bytes, err := json.Marshal(publications)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		case http.MethodPost:
			publication, err := publication.Post(db, r.Body)
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
		case http.MethodPatch:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodDelete:
			err := publication.Delete(db, r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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

	http.HandleFunc("/api/publications", publicationsHandler(db))

	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
