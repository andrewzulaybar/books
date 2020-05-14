package publication

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/andrewzulaybar/books/api/internal/postgres"
)

// Publication represents a specific edition of a work.
type Publication struct {
	ID       int    `json:"id"`
	Author   string `json:"author"`
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
}

// Publications represents a list of publications.
type Publications []Publication

// Delete removes any publication from the database whose ID
// matches any of the IDs given in the request body.
func Delete(db *postgres.DB, body io.Reader) error {
	var identifiers struct {
		IDs []int `json:"ids"`
	}
	err := json.NewDecoder(body).Decode(&identifiers)
	if err != nil {
		message := http.StatusText(http.StatusUnprocessableEntity)
		return errors.New(message)
	}

	for _, id := range identifiers.IDs {
		_, err := db.Exec("DELETE FROM publication WHERE id = $1", id)
		if err != nil {
			message := http.StatusText(http.StatusUnprocessableEntity)
			return errors.New(message)
		}
	}
	return nil
}

// Get retrieves the entire list of publications from the database.
func Get(db *postgres.DB) Publications {
	rows, err := db.Query("SELECT * FROM publication")
	if err != nil {
		return Publications{}
	}

	var publications Publications = Publications{}
	for rows.Next() {
		var id int
		var author, title, imageURL string
		err := rows.Scan(&id, &author, &title, &imageURL)
		if err != nil {
			return Publications{}
		}
		publication := Publication{id, author, title, imageURL}
		publications = append(publications, publication)
	}
	return publications
}

// Post creates a publication from the properties in the request body.
func Post(db *postgres.DB, body io.Reader) (*Publication, error) {
	var publication Publication
	err := json.NewDecoder(body).Decode(&publication)
	if err != nil {
		message := http.StatusText(http.StatusUnprocessableEntity)
		return nil, errors.New(message)
	}

	err = db.QueryRow(
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
