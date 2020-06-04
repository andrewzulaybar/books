package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/gorilla/mux"
)

// Work handles requests made to /api/work/{id:[0-9]+}
func Work(ws *work.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), status.UnprocessableEntity)
			return
		}

		switch r.Method {
		case http.MethodGet:
			s, work := ws.GetWork(id)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(*work)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodPatch:
			var work work.Work
			work.ID = id
			if err := json.NewDecoder(r.Body).Decode(&work); err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}

			s, updated := ws.PatchWork(&work)
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
			s := ws.DeleteWork(id)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}
			w.WriteHeader(s.Code())
		}
	})
}

// Works handles requests made to /api/work
func Works(ws *work.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s, works := ws.GetWorks()
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(works)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodPost:
			var work work.Work
			if err := json.NewDecoder(r.Body).Decode(&work); err != nil {
				http.Error(w, err.Error(), status.UnprocessableEntity)
				return
			}

			s, new := ws.PostWork(&work)
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

			s, ids := ws.DeleteWorks(toDelete.IDs)
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
