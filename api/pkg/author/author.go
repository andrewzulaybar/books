package author

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/lib/pq"
)

// Columns is the comma-separated list of columns found in the author table.
const Columns string = "first_name, last_name, gender, date_of_birth, place_of_birth"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	DeleteAuthor
	FindAuthor
	GetAuthor
	GetAuthors
	PatchAuthor
	PostAuthor
)

// Author represents a writer of a work.
type Author struct {
	ID           int               `json:"id"`
	FirstName    string            `json:"firstName"`
	LastName     string            `json:"lastName"`
	Gender       string            `json:"gender"`
	DateOfBirth  string            `json:"dateOfBirth"`
	PlaceOfBirth location.Location `json:"placeOfBirth"`
}

// Authors represents a list of authors.
type Authors []Author

// Service wraps the database and other dependencies.
type Service struct {
	DB postgres.DB

	LocationService location.Service
}

// Query returns a SQL statement based on the postgres.Query value passed in.
func (s *Service) Query(query postgres.Query, args ...interface{}) string {
	switch query {
	case DeleteAuthor:
		return "DELETE FROM author WHERE id = $1"
	case FindAuthor:
		firstName := args[0].(string)
		lastName := args[1].(string)
		dateOfBirth := args[2].(string)
		return fmt.Sprintf(
			`SELECT id, %s
                        FROM author
                        WHERE first_name = '%s' AND last_name = '%s' AND date_of_birth = '%s'`,
			Columns, firstName, lastName, dateOfBirth,
		)
	case GetAuthor:
		return fmt.Sprintf(
			`SELECT author.id, %s, %s
                        FROM author
                        JOIN location ON author.place_of_birth=location.id
                        WHERE author.id = $1`,
			Columns,
			location.Columns,
		)
	case GetAuthors:
		return fmt.Sprintf(
			`SELECT author.id, %s, %s
                        FROM author
                        JOIN location ON author.place_of_birth=location.id`,
			Columns,
			location.Columns,
		)
	case PatchAuthor:
		var hasUpdate bool
		query := "UPDATE author SET "
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
			return strings.TrimSuffix(query, ",") + fmt.Sprintf(" WHERE id = $1 RETURNING id, %s", Columns)
		}
		return ""
	case PostAuthor:
		return fmt.Sprintf(
			`INSERT INTO author (%s)
                        VALUES ($1, $2, $3, $4, $5)
                        RETURNING *`,
			Columns,
		)
	default:
		return ""
	}
}

// DeleteAuthor removes the entry in the author table matching the given id.
func (s *Service) DeleteAuthor(id int) *status.Status {
	db := s.DB
	deleteAuthor := s.Query(DeleteAuthor)

	res, err := db.Exec(deleteAuthor, id)
	if err != nil {
		log.Printf("[DeleteAuthor] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		log.Printf("[DeleteAuthor] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		msg := fmt.Sprintf("Author with id = %d does not exist", id)
		log.Printf("[DeleteAuthor] %s", msg)
		return status.New(status.OK, msg)
	}
	return status.New(status.NoContent, "")
}

// DeleteAuthors removes the entries in the author table matching the given ids.
func (s *Service) DeleteAuthors(ids []int) (*status.Status, []int) {
	notFound := []int{}
	for _, id := range ids {
		if s := s.DeleteAuthor(id); s.Code() != status.NoContent {
			notFound = append(notFound, id)
		}
	}

	if len(notFound) > 0 {
		msg := fmt.Sprintf("The following authors could not be found: %v", notFound)
		log.Printf("[DeleteAuthors] %s", msg)
		return status.New(status.OK, msg), notFound
	}
	return status.New(status.NoContent, ""), nil
}

// FindAuthor retrieves the author from the database matching the given firstName, lastName, and dateOfBirth.
func (s *Service) FindAuthor(firstName string, lastName string, dateOfBirth string) (*status.Status, *Author) {
	db := s.DB
	findAuthor := s.Query(FindAuthor, firstName, lastName, dateOfBirth)

	var au Author
	row := db.QueryRow(findAuthor)
	if err := row.Scan(
		&au.ID, &au.FirstName, &au.LastName, &au.Gender, &au.DateOfBirth, &au.PlaceOfBirth.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Author ('%s', '%s', '%s') does not exist", firstName, lastName, dateOfBirth)
			log.Printf("[FindAuthor] %s", msg)
			return status.Newf(status.NotFound, msg), nil
		}
		log.Printf("[FindAuthor] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), &au
}

// GetAuthor retrieves the author from the database matching the given id.
func (s *Service) GetAuthor(id int) (*status.Status, *Author) {
	db := s.DB
	getAuthor := s.Query(GetAuthor)

	row := db.QueryRow(getAuthor, id)
	author, err := s.getAuthor(row)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Author with id = %d does not exist", id)
			log.Printf("[GetAuthor] %s", msg)
			return status.Newf(status.NotFound, msg), nil
		}
		log.Printf("[GetAuthor] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), author
}

// GetAuthors retrieves the entire list of authors from the database.
func (s *Service) GetAuthors() (*status.Status, Authors) {
	db := s.DB
	getAuthors := s.Query(GetAuthors)

	rows, err := db.Query(getAuthors)
	if err != nil {
		log.Printf("[GetAuthors] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}

	authors := Authors{}
	for rows.Next() {
		author, err := s.getAuthor(rows)
		if err != nil {
			log.Printf("[GetAuthors] %s", err)
			return status.New(status.InternalServerError, err.Error()), nil
		}
		authors = append(authors, *author)
	}
	return status.New(status.OK, ""), authors
}

// PatchAuthor updates the entry in the database matching author.id with the given attributes.
func (s *Service) PatchAuthor(author *Author) (*status.Status, *Author) {
	if author.PlaceOfBirth != (location.Location{}) {
		if s := s.handleLocation(author); s.Err() != nil {
			log.Printf("[PostAuthor] %s", s.Err())
			return status.New(s.Code(), s.Message()), nil
		}
	}

	db := s.DB
	a := map[string]interface{}{
		"first_name":     author.FirstName,
		"last_name":      author.LastName,
		"gender":         author.Gender,
		"date_of_birth":  author.DateOfBirth,
		"place_of_birth": author.PlaceOfBirth.ID,
	}

	patchAuthor := s.Query(PatchAuthor, a)
	if patchAuthor != "" {
		var au Author
		row := db.QueryRow(patchAuthor, author.ID)
		if err := row.Scan(
			&au.ID, &au.FirstName, &au.LastName, &au.Gender, &au.DateOfBirth, &au.PlaceOfBirth.ID,
		); err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
				log.Printf("[PostAuthor] %s", err)
				return status.New(status.Conflict, err.Error()), nil
			}
			log.Printf("[PatchAuthor] %s", err)
			return status.New(status.UnprocessableEntity, err.Error()), nil
		}
		return status.New(status.OK, ""), &au
	}
	return status.New(status.BadRequest, "No fields in author to update"), nil
}

// PostAuthor creates an entry in the author table with the given attributes.
func (s *Service) PostAuthor(author *Author) (*status.Status, *Author) {
	if author.ID != 0 {
		if s, l := s.GetAuthor(author.ID); s.Code() == status.OK {
			msg := fmt.Sprintf("Author with id = %d already exists", author.ID)
			log.Printf("[PostAuthor] %s", msg)
			return status.Newf(status.Conflict, msg), l
		}
	}

	if author.DateOfBirth == "" {
		author.DateOfBirth = "1970-01-01T00:00:00Z"
	}

	if author.PlaceOfBirth != (location.Location{}) {
		if s := s.handleLocation(author); s.Err() != nil {
			log.Printf("[PostAuthor] %s", s.Err())
			return status.New(s.Code(), s.Message()), nil
		}
	}

	db := s.DB
	postAuthor := s.Query(PostAuthor)

	var au Author
	row := db.QueryRow(
		postAuthor,
		author.FirstName,
		author.LastName,
		author.Gender,
		author.DateOfBirth,
		author.PlaceOfBirth.ID,
	)
	if err := row.Scan(
		&au.ID, &au.FirstName, &au.LastName, &au.Gender, &au.DateOfBirth, &au.PlaceOfBirth.ID,
	); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			log.Printf("[PostAuthor] %s", err)
			return status.New(status.Conflict, err.Error()), nil
		}
		log.Printf("[PostAuthor] %s", err)
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), &au
}

func (s *Service) getAuthor(row interface {
	Scan(dest ...interface{}) error
}) (*Author, error) {
	var author Author
	err := row.Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Gender,
		&author.DateOfBirth,
		&author.PlaceOfBirth.ID,
		&author.PlaceOfBirth.City,
		&author.PlaceOfBirth.Country,
		&author.PlaceOfBirth.Region,
	)
	return &author, err
}

func (s *Service) handleLocation(author *Author) *status.Status {
	pob := &author.PlaceOfBirth

	stat, location := s.LocationService.PostLocation(pob)
	if stat.Err() != nil {
		if stat.Code() != status.Conflict {
			return status.New(stat.Code(), stat.Message())
		}

		if stat, location = s.LocationService.GetLocation(pob.ID); stat.Err() != nil {
			if stat, location = s.LocationService.FindLocation(pob.City, pob.Country); stat.Err() != nil {
				return status.New(stat.Code(), stat.Message())
			}
		}
	}

	pob.ID = location.ID
	return status.New(status.OK, "")
}
