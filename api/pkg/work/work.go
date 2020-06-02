package work

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/lib/pq"
)

// Columns is the comma-separated list of columns found in the work table.
const Columns string = "author_id, description, initial_pub_date, original_language, title"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	DeleteWork
	GetWork
	GetWorks
	PatchWork
	PostWork
)

// A Work represents a literary work.
type Work struct {
	ID               int    `json:"id"`
	AuthorID         int    `json:"authorId"`
	Description      string `json:"description"`
	InitialPubDate   string `json:"initialPubDate"`
	OriginalLanguage string `json:"originalLanguage"`
	Title            string `json:"title"`
}

// Works represents a list of works.
type Works []Work

// Service wraps the database and other dependencies.
type Service struct {
	DB postgres.DB
}

// Query returns a SQL statement based on the postgres.Query value passed in.
func (s *Service) Query(query postgres.Query, args ...interface{}) string {
	switch query {
	case DeleteWork:
		return "DELETE FROM work WHERE id = $1"
	case GetWork:
		return "SELECT * FROM work WHERE id = $1"
	case GetWorks:
		return "SELECT * FROM work ORDER BY id"
	case PatchWork:
		var hasUpdate bool
		query := "UPDATE work SET"
		for column, value := range args[0].(map[string]interface{}) {
			switch value.(type) {
			case string:
				if value != "" {
					query += fmt.Sprintf(" %s = '%s',", column, value)
					hasUpdate = true
				}
			case int:
				if value != 0 {
					query += fmt.Sprintf(" %s = %d,", column, value)
					hasUpdate = true
				}
			}
		}
		if hasUpdate {
			return strings.TrimSuffix(query, ",") + " WHERE id = $1"
		}
		return ""
	case PostWork:
		return fmt.Sprintf(
			`INSERT INTO work (%s)
                        VALUES ($1, $2, $3, $4, $5)
                        RETURNING id`,
			Columns,
		)
	default:
		return ""
	}
}

// DeleteWork removes the entry in the work table matching the given id.
func (s *Service) DeleteWork(id int) *status.Status {
	db := s.DB
	deleteWork := s.Query(DeleteWork)

	res, err := db.Exec(deleteWork, id)
	if err != nil {
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		return status.Newf(status.OK, "Work with id = %d does not exist", id)
	}

	return status.New(status.NoContent, "")
}

// DeleteWorks removes the entries in the work table matching the given ids.
func (s *Service) DeleteWorks(ids []int) (*status.Status, []int) {
	notFound := []int{}
	for _, id := range ids {
		if s := s.DeleteWork(id); s.Code() != status.NoContent {
			notFound = append(notFound, id)
		}
	}

	if len(notFound) > 0 {
		return status.Newf(status.OK, "The following works could not be found: %v", notFound), notFound
	}
	return status.New(status.NoContent, ""), nil
}

// GetWork retrieves the work from the database matching the given id.
func (s *Service) GetWork(id int) (*status.Status, *Work) {
	db := s.DB
	getWork := s.Query(GetWork)

	row := db.QueryRow(getWork, id)
	work, err := s.getWork(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Newf(status.NotFound, "Work with id = %d does not exist", id), nil
		}
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), work
}

// GetWorks retrieves the entire list of works from the database.
func (s *Service) GetWorks() (*status.Status, Works) {
	db := s.DB
	getWorks := s.Query(GetWorks)

	rows, err := db.Query(getWorks)
	if err != nil {
		return status.New(status.InternalServerError, err.Error()), nil
	}

	works := Works{}
	for rows.Next() {
		work, err := s.getWork(rows)
		if err != nil {
			return status.New(status.InternalServerError, err.Error()), nil
		}
		works = append(works, *work)
	}
	return status.New(status.OK, ""), works
}

// PatchWork updates the entry in the database matching work.id with the given attributes.
func (s *Service) PatchWork(work *Work) (*status.Status, *Work) {
	if work != nil && work.ID == 0 {
		return status.New(status.OK, "No fields in work to update"), nil
	}

	db := s.DB
	w := map[string]interface{}{
		"author_id":         work.AuthorID,
		"description":       work.Description,
		"initial_pub_date":  work.InitialPubDate,
		"original_language": work.OriginalLanguage,
		"title":             work.Title,
	}

	patchWork := s.Query(PatchWork, w)
	if patchWork != "" {
		if _, err := db.Exec(patchWork, work.ID); err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
				return status.New(status.Conflict, err.Error()), nil
			}
			return status.New(status.InternalServerError, err.Error()), nil
		}
	}
	return s.GetWork(work.ID)
}

// PostWork creates an entry in the work table with the given attributes.
func (s *Service) PostWork(work *Work) (*status.Status, *Work) {
	if work.ID != 0 {
		s, w := s.GetWork(work.ID)
		if s.Err() != nil {
			return status.New(status.UnprocessableEntity, s.Message()), nil
		}
		return status.New(status.OK, ""), w
	}

	db := s.DB
	if err := db.QueryRow(
		s.Query(PostWork),
		work.AuthorID,
		work.Description,
		work.InitialPubDate,
		work.OriginalLanguage,
		work.Title,
	).Scan(&(work.ID)); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return status.New(status.Conflict, err.Error()), nil
		}
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), work
}

func (s *Service) getWork(row interface {
	Scan(dest ...interface{}) error
}) (*Work, error) {
	var work Work
	if err := row.Scan(
		&work.ID,
		&work.AuthorID,
		&work.Description,
		&work.InitialPubDate,
		&work.OriginalLanguage,
		&work.Title,
	); err != nil {
		return nil, err
	}
	return &work, nil
}
