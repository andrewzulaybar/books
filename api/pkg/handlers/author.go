package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/gorilla/mux"
)

// Author handles requests made to /api/author/{id:[0-9]+}
func Author(a *author.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), status.UnprocessableEntity)
			return
		}

		switch r.Method {
		case http.MethodGet:
			s, author := a.GetAuthor(id)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(*author)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodPatch:
			var author author.Author
			author.ID = id
			if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}

			s, updated := a.PatchAuthor(&author)
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
			s := a.DeleteAuthor(id)
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}
			w.WriteHeader(s.Code())
		}
	})
}

// Authors handles requests made to /api/author
func Authors(a *author.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s, authors := a.GetAuthors()
			if s.Err() != nil {
				http.Error(w, s.Message(), s.Code())
				return
			}

			bytes, err := json.Marshal(authors)
			if err != nil {
				http.Error(w, err.Error(), status.InternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(s.Code())
			w.Write(bytes)
		case http.MethodPost:
			var author author.Author
			if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
				http.Error(w, err.Error(), status.UnprocessableEntity)
				return
			}

			s, new := a.PostAuthor(&author)
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

			s, ids := a.DeleteAuthors(toDelete.IDs)
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
