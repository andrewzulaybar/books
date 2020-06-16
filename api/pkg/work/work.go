package work

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/lib/pq"
)

// Columns is the comma-separated list of columns found in the work table.
const Columns string = "description, initial_pub_date, original_language, title, author_id"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	DeleteWork
	FindWork
	GetWork
	GetWorks
	PatchWork
	PostWork
)

// A Work represents a literary work.
type Work struct {
	ID               int           `json:"id"`
	Description      string        `json:"description"`
	InitialPubDate   string        `json:"initialPubDate"`
	OriginalLanguage string        `json:"originalLanguage"`
	Title            string        `json:"title"`
	Author           author.Author `json:"author"`
}

// Works represents a list of works.
type Works []Work

// Service wraps the database and other dependencies.
type Service struct {
	DB postgres.DB

	AuthorService author.Service
}

// Query returns a SQL statement based on the postgres.Query value passed in.
func (s *Service) Query(query postgres.Query, args ...interface{}) string {
	switch query {
	case DeleteWork:
		return "DELETE FROM work WHERE id = $1"
	case FindWork:
		title := args[0].(string)
		authorID := args[1].(int)
		return fmt.Sprintf(
			`SELECT id, %s
                        FROM work
                        WHERE title = '%s' AND author_id = %d`,
			Columns, title, authorID,
		)
	case GetWork:
		return fmt.Sprintf(
			`SELECT work.id, %s, %s, %s
                        FROM work
                        JOIN author ON work.author_id=author.id
                        JOIN location ON author.place_of_birth=location.id
                        WHERE work.id = $1`,
			Columns,
			author.Columns,
			location.Columns,
		)
	case GetWorks:
		return fmt.Sprintf(
			`SELECT work.id, %s, %s, %s
                        FROM work
                        JOIN author ON work.author_id=author.id
                        JOIN location ON author.place_of_birth=location.id`,
			Columns,
			author.Columns,
			location.Columns,
		)
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
			return strings.TrimSuffix(query, ",") + fmt.Sprintf(" WHERE id = $1 RETURNING id, %s", Columns)
		}
		return ""
	case PostWork:
		return fmt.Sprintf(
			`INSERT INTO work (%s)
                        VALUES ($1, $2, $3, $4, $5)
                        RETURNING id, %s`,
			Columns,
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
		log.Printf("[DeleteWork] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		log.Printf("[DeleteWork] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		msg := fmt.Sprintf("Work with id = %d does not exist", id)
		log.Printf("[DeleteWork] %s", msg)
		return status.Newf(status.OK, msg)
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
		msg := fmt.Sprintf("The following works could not be found: %v", notFound)
		log.Printf("[DeleteWorks] %s", msg)
		return status.Newf(status.OK, msg), notFound
	}
	return status.New(status.NoContent, ""), nil
}

// FindWork retrieves the work from the database matching the given title and authorID.
func (s *Service) FindWork(title string, authorID int) (*status.Status, *Work) {
	db := s.DB
	findWork := s.Query(FindWork, title, authorID)

	var wk Work
	row := db.QueryRow(findWork)
	if err := row.Scan(
		&wk.ID, &wk.Description, &wk.InitialPubDate, &wk.OriginalLanguage, &wk.Title, &wk.Author.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Work ('%s', %d) does not exist", title, authorID)
			log.Printf("[FindWork] %s", msg)
			return status.Newf(status.NotFound, msg), nil
		}
		log.Printf("[FindWork] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), &wk
}

// GetWork retrieves the work from the database matching the given id.
func (s *Service) GetWork(id int) (*status.Status, *Work) {
	db := s.DB
	getWork := s.Query(GetWork)

	row := db.QueryRow(getWork, id)
	work, err := s.getWork(row)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Work with id = %d does not exist", id)
			log.Printf("[GetWork] %s", msg)
			return status.Newf(status.NotFound, msg), nil
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
		log.Printf("[GetWorks] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}

	works := Works{}
	for rows.Next() {
		work, err := s.getWork(rows)
		if err != nil {
			log.Printf("[GetWorks] %s", err)
			return status.New(status.InternalServerError, err.Error()), nil
		}
		works = append(works, *work)
	}
	return status.New(status.OK, ""), works
}

// PatchWork updates the entry in the database matching work.id with the given attributes.
func (s *Service) PatchWork(work *Work) (*status.Status, *Work) {
	if work.Author != (author.Author{}) {
		if s := s.handleAuthor(work); s.Err() != nil {
			log.Printf("[PatchWork] %s", s.Err())
			return status.New(s.Code(), s.Message()), nil
		}
	}

	db := s.DB
	w := map[string]interface{}{
		"description":       work.Description,
		"initial_pub_date":  work.InitialPubDate,
		"original_language": work.OriginalLanguage,
		"title":             work.Title,
		"author_id":         work.Author.ID,
	}

	patchWork := s.Query(PatchWork, w)
	if patchWork != "" {
		var wk Work
		row := db.QueryRow(patchWork, work.ID)
		if err := row.Scan(
			&wk.ID, &wk.Description, &wk.InitialPubDate, &wk.OriginalLanguage, &wk.Title, &wk.Author.ID,
		); err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
				log.Printf("[PatchWork] %s", err)
				return status.New(status.Conflict, err.Error()), nil
			}
			log.Printf("[PatchWork] %s", err)
			return status.New(status.InternalServerError, err.Error()), nil
		}
		return status.New(status.OK, ""), &wk
	}
	return status.New(status.OK, "No fields in work to update"), nil
}

// PostWork creates an entry in the work table with the given attributes.
func (s *Service) PostWork(work *Work) (*status.Status, *Work) {
	if work.ID != 0 {
		if s, l := s.GetWork(work.ID); s.Code() == status.OK {
			msg := fmt.Sprintf("Work with id = %d already exists", work.ID)
			log.Printf("[PostWork] %s", msg)
			return status.Newf(status.Conflict, msg), l
		}
	}

	if work.Author != (author.Author{}) {
		if s := s.handleAuthor(work); s.Err() != nil {
			log.Printf("[PostWork] %s", s.Err())
			return status.New(s.Code(), s.Message()), nil
		}
	}

	db := s.DB
	postWork := s.Query(PostWork)

	var wk Work
	row := db.QueryRow(
		postWork,
		work.Description,
		work.InitialPubDate,
		work.OriginalLanguage,
		work.Title,
		work.Author.ID,
	)
	if err := row.Scan(
		&wk.ID,
		&wk.Description,
		&wk.InitialPubDate,
		&wk.OriginalLanguage,
		&wk.Title,
		&wk.Author.ID,
	); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			log.Printf("[PostWork]: %s", err)
			return status.New(status.Conflict, err.Error()), nil
		}
		log.Printf("[PostWork]: %s", err)
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), &wk
}

func (s *Service) getWork(row interface {
	Scan(dest ...interface{}) error
}) (*Work, error) {
	var work Work
	author := &work.Author
	location := &author.PlaceOfBirth
	if err := row.Scan(
		&work.ID,
		&work.Description,
		&work.InitialPubDate,
		&work.OriginalLanguage,
		&work.Title,
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Gender,
		&author.DateOfBirth,
		&location.ID,
		&location.City,
		&location.Country,
		&location.Region,
	); err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *Service) handleAuthor(work *Work) *status.Status {
	au := &work.Author

	if au.ID != 0 {
		return status.New(status.OK, "")
	}

	stat, author := s.AuthorService.PostAuthor(au)
	if stat.Err() != nil {
		if stat.Code() != status.Conflict {
			return status.New(stat.Code(), stat.Message())
		}

		if stat, author = s.AuthorService.GetAuthor(au.ID); stat.Err() != nil {
			stat, author = s.AuthorService.FindAuthor(au.FirstName, au.LastName, au.DateOfBirth)
			if stat.Err() != nil {
				return status.New(stat.Code(), stat.Message())
			}
		}
	}

	au.ID = author.ID
	return status.New(status.OK, "")
}
