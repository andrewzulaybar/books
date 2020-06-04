package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/gorilla/mux"
)

// Publication handles requests made to /api/publication/{id:[0-9]+}
func Publication(p *publication.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), status.UnprocessableEntity)
			return
		}

		switch r.Method {
		case http.MethodGet:
			s, pub := p.GetPublication(id)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(*pub)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodPatch:
			var pub publication.Publication
			pub.ID = id
			if err := json.NewDecoder(r.Body).Decode(&pub); err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}

			s, updated := p.PatchPublication(&pub)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(*updated)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodDelete:
			s := p.DeletePublication(id)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}
			w.WriteHeader(s.Code())
		}
	})
}

// Publications handles requests made to /api/publication
func Publications(p *publication.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s, pubs := p.GetPublications()
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(pubs)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodPost:
			var pub publication.Publication
			if err := json.NewDecoder(r.Body).Decode(&pub); err != nil {
				http.Error(w, err.Error(), status.UnprocessableEntity)
				return
			}

			s, new := p.PostPublication(&pub)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(*new)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodDelete:
			type identifiers struct {
				IDs []int `json:"ids"`
			}

			var toDelete identifiers
			if err := json.NewDecoder(r.Body).Decode(&toDelete); err != nil {
				http.Error(w, err.Error(), status.UnprocessableEntity)
				return
			}

			s, ids := p.DeletePublications(toDelete.IDs)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}
			if ids != nil {
				notFound := identifiers{IDs: ids}
				bytes, err := json.Marshal(notFound)
				if err != nil {
					http.Error(w, err.Error(), status.InternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(bytes)
				return
			}
			w.WriteHeader(s.Code())
		}
	})
}
