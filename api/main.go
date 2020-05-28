package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func publicationsHandler(p *publication.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// case http.MethodGet:
		// 	publications, err := publication.GetMany(db)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	bytes, err := json.Marshal(publications)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Write(bytes)
		// case http.MethodPost:
		// 	publication, err := publication.PostOne(db, r.Body)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	bytes, err := json.Marshal(publication)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusCreated)
		// 	w.Write(bytes)
		// case http.MethodDelete:
		// 	if err := publication.DeleteMany(db, r.Body); err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	w.WriteHeader(http.StatusNoContent)
		}
	})
}

func publicationHandler(p *publication.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), status.UnprocessableEntity)
			return
		}

		switch r.Method {
		// case http.MethodGet:
		// 	publication, err := publication.GetOne(db, ID)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	bytes, err := json.Marshal(publication)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Write(bytes)
		// case http.MethodPatch:
		// 	if err := publication.PatchOne(db, r.Body, ID); err != nil {
		// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		// 		return
		// 	}
		// 	w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			if s := p.DeletePublication(id); s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}
			w.WriteHeader(status.NoContent)
		}
	})
}

func main() {
	conf, err := config.Load("config/.env")
	if err != nil {
		panic(err)
	}

	db := postgres.Setup(conf.ConnectionString, "internal/sql/")
	defer db.Disconnect()

	p := &publication.Service{DB: *db}

	r := mux.NewRouter()
	API := r.PathPrefix("/api").Subrouter()

	API.HandleFunc("/publications", publicationsHandler(p)).
		Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	API.HandleFunc("/publications/{id:[0-9]+}", publicationHandler(p)).
		Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)

	srv := &http.Server{
		Handler:      handlers.CORS()(r),
		Addr:         conf.Address,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
