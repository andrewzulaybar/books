package main

import (
	"log"
	"net/http"
	"time"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/handlers"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/andrewzulaybar/books/api/test/data"
	h "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	conf, err := config.Load(".env")
	if err != nil {
		panic(err)
	}

	db, dc := postgres.Setup(conf.ConnectionString)
	defer dc()

	l := &location.Service{DB: *db}
	a := &author.Service{DB: *db, LocationService: *l}
	w := &work.Service{DB: *db, AuthorService: *a}
	p := &publication.Service{DB: *db, WorkService: *w}
	data.LoadPublications(p)

	r := mux.NewRouter()
	API := r.PathPrefix("/api").Subrouter()

	API.HandleFunc("/publication", handlers.Publications(p)).
		Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	API.HandleFunc("/publication/{id:[0-9]+}", handlers.Publication(p)).
		Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
	API.HandleFunc("/work", handlers.Works(w)).
		Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	API.HandleFunc("/work/{id:[0-9]+}", handlers.Work(w)).
		Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
	API.HandleFunc("/author", handlers.Authors(a)).
		Methods(http.MethodGet, http.MethodPost, http.MethodDelete)
	API.HandleFunc("/author/{id:[0-9]+}", handlers.Author(a)).
		Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)

	srv := &http.Server{
		Handler:      h.CORS()(r),
		Addr:         conf.Address,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
