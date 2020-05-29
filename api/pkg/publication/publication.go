package publication

import (
	"database/sql"
	"fmt"

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
func (s Service) QueryMap() map[string]string {
	return map[string]string{
		"DeletePublication": `DELETE FROM publication WHERE id = $1`,
		"GetPublication": fmt.Sprintf(
			`SELECT %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id
                        WHERE publication.id = $1`,
			columns,
		),
	}
}

// DeletePublication removes the entry in the publication table matching the given id.
func (s *Service) DeletePublication(id int) *status.Status {
	db := s.DB
	deletePublication := s.QueryMap()["DeletePublication"]

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
	getPublication := s.QueryMap()["GetPublication"]

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

// GetMany retrieves the entire list of publications from the database.
// func GetMany(db *postgres.DB) (Publications, error) {
// 	query := fmt.Sprintf(
// 		`SELECT %s
//                 FROM publication
//                 JOIN work ON publication.work_id=work.id`,
// 		columns,
// 	)
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return Publications{}, err
// 	}

// 	var publications Publications = Publications{}
// 	for rows.Next() {
// 		publication, err := getPublication(rows)
// 		if err != nil {
// 			return Publications{}, err
// 		}
// 		publications = append(publications, publication)
// 	}
// 	return publications, nil
// }

// PatchOne updates the entry in the database matching the given ID
// with the attributes passed in the request body.
// func PatchOne(db *postgres.DB, body io.Reader, ID int) error {
// 	var publication Publication
// 	publication.ID = ID
// 	if err := json.NewDecoder(body).Decode(&publication); err != nil {
// 		return err
// 	}

// 	if err := updateWork(db, &publication); err != nil {
// 		return err
// 	}

// 	if err := updatePublication(db, &publication); err != nil {
// 		return err
// 	}

// 	return nil
// }

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

// func updatePublication(db *postgres.DB, publication *Publication) error {
// 	p := map[string]interface{}{
// 		"edition_pub_date": publication.EditionPubDate,
// 		"format":           publication.Format,
// 		"image_url":        publication.ImageURL,
// 		"isbn":             publication.ISBN,
// 		"isbn13":           publication.ISBN13,
// 		"language":         publication.Language,
// 		"num_pages":        publication.NumPages,
// 		"publisher":        publication.Publisher,
// 	}

// 	var hasUpdate bool
// 	query := "UPDATE publication SET "
// 	for column, value := range p {
// 		switch value.(type) {
// 		case string:
// 			if value != "" {
// 				query += fmt.Sprintf("%s = '%s',", column, value)
// 				hasUpdate = true
// 			}
// 		case int:
// 			if value != 0 {
// 				query += fmt.Sprintf("%s = %d,", column, value)
// 				hasUpdate = true
// 			}
// 		}
// 	}
// 	if hasUpdate {
// 		query = strings.TrimSuffix(query, ",") + " WHERE id = $1"
// 		if _, err := db.Exec(query, publication.ID); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func updateWork(db *postgres.DB, publication *Publication) error {
// 	w := map[string]string{
// 		"author":            publication.Author,
// 		"description":       publication.Description,
// 		"initial_pub_date":  publication.InitialPubDate,
// 		"original_language": publication.OriginalLanguage,
// 		"title":             publication.Title,
// 	}

// 	var hasUpdate bool
// 	query := "UPDATE work SET "
// 	for column, value := range w {
// 		if value != "" {
// 			query += fmt.Sprintf("%s = '%s',", column, value)
// 			hasUpdate = true
// 		}
// 	}
// 	if hasUpdate {
// 		query = strings.TrimSuffix(query, ",") +
// 			" WHERE id = (SELECT work_id FROM publication WHERE id = $1)"
// 		if _, err := db.Exec(query, publication.ID); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
