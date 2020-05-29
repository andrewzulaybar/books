package publication

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

const columns string = `publication.id, edition_pub_date, format, image_url, isbn, isbn13, language, num_pages,
        publisher, work.id, author_id, description, initial_pub_date, original_language, title`

// Publication represents a specific edition of a work.
type Publication struct {
	ID             int       `json:"id"`
	EditionPubDate string    `json:"editionPubDate"`
	Format         string    `json:"format"`
	ImageURL       string    `json:"imageUrl"`
	ISBN           string    `json:"isbn"`
	ISBN13         string    `json:"isbn13"`
	Language       string    `json:"language"`
	NumPages       int       `json:"numPages"`
	Publisher      string    `json:"publisher"`
	Work           work.Work `json:"work"`
}

// Publications represents a list of publications.
type Publications []Publication

// Service wraps the database.
type Service struct {
	DB postgres.DB
}

// QueryMap returns a map of the SQL queries to be used within this service.
func (s Service) QueryMap() map[string]interface{} {
	return map[string]interface{}{
		"DeletePublication": `DELETE FROM publication WHERE id = $1`,
		"GetPublication": fmt.Sprintf(
			`SELECT %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id
                        WHERE publication.id = $1`,
			columns,
		),
		"GetPublications": fmt.Sprintf(
			`SELECT %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id`,
			columns,
		),
		"PatchPublication": func(p map[string]interface{}) string {
			var hasUpdate bool
			query := "UPDATE publication SET "
			for column, value := range p {
				switch value.(type) {
				case string:
					if value != "" {
						query += fmt.Sprintf("%s = '%s',", column, value)
						hasUpdate = true
					}
				case int:
					if value != 0 {
						query += fmt.Sprintf("%s = %d,", column, value)
						hasUpdate = true
					}
				}
			}
			if hasUpdate {
				return strings.TrimSuffix(query, ",") + " WHERE id = $1"
			}
			return ""
		},
		"PatchWork": func(w map[string]interface{}) string {
			var hasUpdate bool
			query := "UPDATE work SET "
			for column, value := range w {
				switch value.(type) {
				case string:
					if value != "" {
						query += fmt.Sprintf("%s = '%s',", column, value)
						hasUpdate = true
					}
				case int:
					if value != 0 {
						query += fmt.Sprintf("%s = %d,", column, value)
						hasUpdate = true
					}
				}
			}
			if hasUpdate {
				return strings.TrimSuffix(query, ",") +
					" WHERE id = (SELECT work_id FROM publication WHERE id = $1)"
			}
			return ""
		},
	}
}

// DeletePublication removes the entry in the publication table matching the given id.
func (s *Service) DeletePublication(id int) *status.Status {
	db := s.DB
	deletePublication := s.QueryMap()["DeletePublication"].(string)

	res, err := db.Exec(deletePublication, id)
	if err != nil {
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		return status.Newf(status.NotFound, "Publication with id = %d does not exist", id)
	}

	return status.New(status.NoContent, "")
}

// DeletePublications removes the entries in the publication table matching the given ids.
func (s *Service) DeletePublications(ids []int) (*status.Status, []int) {
	notFound := []int{}
	for _, id := range ids {
		if s := s.DeletePublication(id); s.Code() == status.NotFound {
			notFound = append(notFound, id)
		}
	}

	if len(notFound) > 0 {
		return status.Newf(status.OK, "The following publications could not be found: %v", notFound), notFound
	}
	return status.New(status.NoContent, ""), nil
}

// GetPublication retrieves the publication from the database matching the given ID.
func (s *Service) GetPublication(id int) (*status.Status, *Publication) {
	db := s.DB
	getPublication := s.QueryMap()["GetPublication"].(string)

	row := db.QueryRow(getPublication, id)
	pub, err := s.getPublication(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Newf(status.NotFound, "Publication with id = %d does not exist", id), nil
		}
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), pub
}

// GetPublications retrieves the entire list of publications from the database.
func (s *Service) GetPublications() (*status.Status, Publications) {
	db := s.DB
	getPublications := s.QueryMap()["GetPublications"].(string)

	rows, err := db.Query(getPublications)
	if err != nil {
		return status.New(status.InternalServerError, err.Error()), nil
	}

	publications := Publications{}
	for rows.Next() {
		pub, err := s.getPublication(rows)
		if err != nil {
			return status.New(status.InternalServerError, err.Error()), nil
		}
		publications = append(publications, *pub)
	}
	return status.New(status.OK, ""), publications
}

// PatchPublication updates the entry in the database matching the given ID
// with the attributes passed in the request body.
func (s *Service) PatchPublication(pub *Publication) (*status.Status, *Publication) {
	if err := s.updateWork(pub); err != nil {
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}

	if err := s.updatePublication(pub); err != nil {
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}

	return s.GetPublication(pub.ID)
}

// PostOne creates a publication from the properties in the request body.
// func PostOne(db *postgres.DB, body io.Reader) (*Publication, error) {
// 	var publication Publication
// 	if err := json.NewDecoder(body).Decode(&publication); err != nil {
// 		return nil, err
// 	}

// 	if err := postWork(db, &publication); err != nil {
// 		return nil, err
// 	}

// 	if err := postPublication(db, &publication); err != nil {
// 		return nil, err
// 	}

// 	return &publication, nil
// }

func (s *Service) getPublication(row interface {
	Scan(dest ...interface{}) error
}) (*Publication, error) {
	var pub Publication
	err := row.Scan(
		&pub.ID,
		&pub.EditionPubDate,
		&pub.Format,
		&pub.ImageURL,
		&pub.ISBN,
		&pub.ISBN13,
		&pub.Language,
		&pub.NumPages,
		&pub.Publisher,
		&pub.Work.ID,
		&pub.Work.AuthorID,
		&pub.Work.Description,
		&pub.Work.InitialPubDate,
		&pub.Work.OriginalLanguage,
		&pub.Work.Title,
	)
	return &pub, err
}

// func postPublication(db *postgres.DB, publication *Publication) error {
// 	row := db.QueryRow(
// 		`INSERT INTO publication
//                         (edition_pub_date, format, image_url, isbn, isbn13,
//                                 language, num_pages, publisher, work_id)
//                 VALUES
//                         ($1, $2, $3, $4, $5, $6, $7, $8, $9)
//                 RETURNING id`,
// 		publication.EditionPubDate,
// 		publication.Format,
// 		publication.ImageURL,
// 		publication.ISBN,
// 		publication.ISBN13,
// 		publication.Language,
// 		publication.NumPages,
// 		publication.Publisher,
// 		publication.WorkID,
// 	)
// 	return row.Scan(&(publication.ID))
// }

// func postWork(db *postgres.DB, publication *Publication) error {
// 	row := db.QueryRow(
// 		`INSERT INTO work
//                         (author, description, initial_pub_date, original_language, title)
//                 VALUES
//                         ($1, $2, $3, $4, $5)
//                 RETURNING id`,
// 		publication.Author,
// 		publication.Description,
// 		publication.InitialPubDate,
// 		publication.OriginalLanguage,
// 		publication.Title,
// 	)
// 	return row.Scan(&(publication.WorkID))
// }

func (s *Service) updatePublication(pub *Publication) error {
	db := s.DB

	p := map[string]interface{}{
		"edition_pub_date": pub.EditionPubDate,
		"format":           pub.Format,
		"image_url":        pub.ImageURL,
		"isbn":             pub.ISBN,
		"isbn13":           pub.ISBN13,
		"language":         pub.Language,
		"num_pages":        pub.NumPages,
		"publisher":        pub.Publisher,
	}

	queryBuilder := s.QueryMap()["PatchPublication"]
	if patchPublication := queryBuilder.(func(map[string]interface{}) string)(p); patchPublication != "" {
		if _, err := db.Exec(patchPublication, pub.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) updateWork(pub *Publication) error {
	db := s.DB

	w := map[string]interface{}{
		"author_id":         pub.Work.AuthorID,
		"description":       pub.Work.Description,
		"initial_pub_date":  pub.Work.InitialPubDate,
		"original_language": pub.Work.OriginalLanguage,
		"title":             pub.Work.Title,
	}

	queryBuilder := s.QueryMap()["PatchWork"]
	if patchWork := queryBuilder.(func(map[string]interface{}) string)(w); patchWork != "" {
		if _, err := db.Exec(patchWork, pub.ID); err != nil {
			return err
		}
	}
	return nil
}
