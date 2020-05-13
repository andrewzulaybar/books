package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/andrewzulaybar/books/api/db"
)

var database *sql.DB

// Publication represents a specific edition of a work
type Publication struct {
	ID       int    `json:"id"`
	Author   string `json:"author"`
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
}

func main() {
	pool, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Disconnect(pool)

	err = db.Init(pool)
	if err != nil {
		panic(err)
	}
	database = pool

	http.HandleFunc("/api/publications", PublicationsHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// PublicationsHandler handles requests made to /api/publications
func PublicationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		publications := getPublications()
		bytes, err := json.Marshal(publications)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	case "POST":
		publication, err := createPublication(r.Body)
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
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func createPublication(body io.Reader) (*Publication, error) {
	var publication Publication
	err := json.NewDecoder(body).Decode(&publication)
	if err != nil {
		message := http.StatusText(http.StatusUnprocessableEntity)
		return nil, errors.New(message)
	}

	err = database.QueryRow(
		`INSERT INTO publication
                        (author, image_url, title)
                VALUES
                        ($1, $2, $3)
                RETURNING id`,
		publication.Author,
		publication.ImageURL,
		publication.Title,
	).Scan(&(publication.ID))
	if err != nil {
		message := http.StatusText(http.StatusUnprocessableEntity)
		return nil, errors.New(message)
	}
	return &publication, nil
}

func getPublications() []Publication {
	rows, err := database.Query("SELECT * FROM publication")
	if err != nil {
		return []Publication{}
	}

	var publications []Publication = []Publication{}
	for rows.Next() {
		var id int
		var author, title, imageURL string
		err := rows.Scan(&id, &author, &title, &imageURL)
		if err != nil {
			return []Publication{}
		}
		publication := Publication{id, author, title, imageURL}
		publications = append(publications, publication)
	}
	return publications
}
