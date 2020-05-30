package work

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
)

// Columns is the comma-separated list of columns found in the work table.
const Columns string = "author_id, description, initial_pub_date, original_language, title"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	GetWork
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
	case GetWork:
		return `SELECT *
                        FROM work
                        WHERE id = $1`
	case PatchWork:
		var hasUpdate bool
		query := "UPDATE work SET "
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

// GetWork retrieves the work from the database matching the given id.
func (s *Service) GetWork(id int) (*status.Status, *Work) {
	db := s.DB

	var work Work
	err := db.QueryRow(s.Query(GetWork), id).Scan(
		&work.ID,
		&work.AuthorID,
		&work.Description,
		&work.InitialPubDate,
		&work.OriginalLanguage,
		&work.Title,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Newf(status.NotFound, "Work with id = %d does not exist", id), nil
		}
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), &work
}

// PatchWork updates the entry in the database matching work.id with the given attributes.
func (s *Service) PatchWork(work *Work) (*status.Status, *Work) {
	if work != nil && *work == (Work{}) {
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
			return status.New(status.InternalServerError, err.Error()), nil
		}
	}
	return s.GetWork(work.ID)
}

// PostWork creates an entry in the work table with the given attributes.
func (s *Service) PostWork(work *Work) (*status.Status, *Work) {
	db := s.DB

	if work.ID != 0 {
		s, w := s.GetWork(work.ID)
		if s.Err() != nil {
			return status.New(status.UnprocessableEntity, s.Message()), nil
		}
		return status.New(status.OK, ""), w
	}

	if err := db.QueryRow(
		s.Query(PostWork),
		work.AuthorID,
		work.Description,
		work.InitialPubDate,
		work.OriginalLanguage,
		work.Title,
	).Scan(&(work.ID)); err != nil {
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), work
}
