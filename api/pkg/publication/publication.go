package publication

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/andrewzulaybar/books/api/internal/postgres"
)

// Publication represents a specific edition of a work.
type Publication struct {
	ID             int    `json:"id"`
	Author         string `json:"author"`
	EditionPubDate string `json:"edition_pub_date"`
	Format         string `json:"format"`
	ImageURL       string `json:"image_url"`
	ISBN           string `json:"isbn"`
	ISBN13         string `json:"isbn13"`
	Language       string `json:"language"`
	NumPages       int    `json:"num_pages"`
	Publisher      string `json:"publisher"`
	Title          string `json:"title"`
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
		var publication Publication
		if err := rows.Scan(
			&publication.ID,
			&publication.Author,
			&publication.EditionPubDate,
			&publication.Format,
			&publication.ImageURL,
			&publication.ISBN,
			&publication.ISBN13,
			&publication.Language,
			&publication.NumPages,
			&publication.Publisher,
			&publication.Title,
		); err != nil {
			log.Printf(err.Error())
			continue
		}
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
                        (
                                author,
                                edition_pub_date,
                                format,
                                image_url,
                                isbn,
                                isbn13,
                                language,
                                num_pages,
                                publisher,
                                title
                        )
                VALUES
                        ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                RETURNING id`,
		publication.Author,
		publication.EditionPubDate,
		publication.Format,
		publication.ImageURL,
		publication.ISBN,
		publication.ISBN13,
		publication.Language,
		publication.NumPages,
		publication.Publisher,
		publication.Title,
	).Scan(&(publication.ID))
	if err != nil {
		message := http.StatusText(http.StatusUnprocessableEntity)
		return nil, errors.New(message)
	}
	return &publication, nil
}
