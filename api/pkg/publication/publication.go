package publication

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

// Columns is the comma-separated list of columns found in the publication table.
const Columns string = "edition_pub_date, format, image_url, isbn, isbn13, language, num_pages, publisher, work_id"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	DeletePublication
	GetPublication
	GetPublications
	PatchPublication
	PostPublication
)

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

// Service wraps the database and other dependencies.
type Service struct {
	DB postgres.DB

	WorkService work.Service
}

// Query returns a SQL statement based on the postgres.Query value passed in.
func (s *Service) Query(query postgres.Query, args ...interface{}) string {
	switch query {
	case DeletePublication:
		return "DELETE FROM publication WHERE id = $1"
	case GetPublication:
		return fmt.Sprintf(
			`SELECT publication.id, %s, %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id
                        WHERE publication.id = $1`,
			Columns,
			work.Columns,
		)
	case GetPublications:
		return fmt.Sprintf(
			`SELECT publication.id, %s, %s
                        FROM publication
                        JOIN work ON publication.work_id=work.id`,
			Columns,
			work.Columns,
		)
	case PatchPublication:
		var hasUpdate bool
		query := "UPDATE publication SET "
		for column, value := range args[0].(map[string]interface{}) {
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
	case PostPublication:
		return fmt.Sprintf(
			`INSERT INTO publication (%s)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                        RETURNING id`,
			Columns,
		)
	default:
		return ""
	}
}

// DeletePublication removes the entry in the publication table matching the given id.
func (s *Service) DeletePublication(id int) *status.Status {
	db := s.DB
	deletePublication := s.Query(DeletePublication)

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
	getPublication := s.Query(GetPublication)

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
	getPublications := s.Query(GetPublications)

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

// PatchPublication updates the entry in the database matching pub.id with the given attributes.
func (s *Service) PatchPublication(pub *Publication) (*status.Status, *Publication) {
	stat, _ := s.WorkService.PatchWork(&pub.Work)
	if err := stat.Err(); err != nil {
		return status.New(stat.Code(), stat.Message()), nil
	}

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

	patchPublication := s.Query(PatchPublication, p)
	if patchPublication != "" {
		if _, err := db.Exec(patchPublication, pub.ID); err != nil {
			return status.New(status.UnprocessableEntity, err.Error()), nil
		}
	}
	return s.GetPublication(pub.ID)
}

// PostPublication creates an entry in the publication table with the given attributes.
func (s *Service) PostPublication(pub *Publication) (*status.Status, *Publication) {
	stat, _ := s.WorkService.PostWork(&pub.Work)
	if err := stat.Err(); err != nil {
		return status.New(stat.Code(), stat.Message()), nil
	}

	db := s.DB
	if err := db.QueryRow(
		s.Query(PostPublication),
		pub.EditionPubDate,
		pub.Format,
		pub.ImageURL,
		pub.ISBN,
		pub.ISBN13,
		pub.Language,
		pub.NumPages,
		pub.Publisher,
		pub.Work.ID,
	).Scan(&(pub.ID)); err != nil {
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), pub
}

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
