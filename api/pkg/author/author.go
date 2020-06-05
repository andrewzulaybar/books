package author

import (
	"database/sql"
	"fmt"
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
	GetAuthor
	GetAuthors
	PostAuthor
)

// Author represents a writer of a work.
type Author struct {
	ID           int               `json:"id"`
	FirstName    string            `json:"firstName"`
	LastName     string            `json:"lastName"`
	Gender       string            `json:"gender"`
	DateOfBirth  *string           `json:"dateOfBirth"`
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
	case PostAuthor:
		query := "INSERT INTO author ("
		values := []interface{}{}
		for column, value := range args[0].(map[string]interface{}) {
			if value != "" && value != 0 {
				query += fmt.Sprintf("%s, ", column)
				values = append(values, value)
			}
		}
		query = strings.TrimSuffix(query, ", ") + ") VALUES ("
		for _, value := range values {
			switch value.(type) {
			case string:
				query += fmt.Sprintf("'%s', ", value)
			case int:
				query += fmt.Sprintf("%d, ", value)
			}
		}
		query = strings.TrimSuffix(query, ", ") + ") RETURNING id"
		return query
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
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		return status.Newf(status.NotFound, "Author with id = %d does not exist", id)
	}

	return status.New(status.NoContent, "")
}

// DeleteAuthors removes the entries in the author table matching the given ids.
func (s *Service) DeleteAuthors(ids []int) (*status.Status, []int) {
	notFound := []int{}
	for _, id := range ids {
		if s := s.DeleteAuthor(id); s.Code() == status.NotFound {
			notFound = append(notFound, id)
		}
	}

	if len(notFound) > 0 {
		return status.Newf(status.OK, "The following authors could not be found: %v", notFound), notFound
	}
	return status.New(status.NoContent, ""), nil
}

// GetAuthor retrieves the author from the database matching the given id.
func (s *Service) GetAuthor(id int) (*status.Status, *Author) {
	db := s.DB
	getAuthor := s.Query(GetAuthor)

	row := db.QueryRow(getAuthor, id)
	author, err := s.getAuthor(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Newf(status.NotFound, "Author with id = %d does not exist", id), nil
		}
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
		return status.New(status.InternalServerError, err.Error()), nil
	}

	authors := Authors{}
	for rows.Next() {
		author, err := s.getAuthor(rows)
		if err != nil {
			return status.New(status.InternalServerError, err.Error()), nil
		}
		authors = append(authors, *author)
	}
	return status.New(status.OK, ""), authors
}

// PostAuthor creates an entry in the author table with the given attributes.
func (s *Service) PostAuthor(author *Author) (*status.Status, *Author) {
	if author.PlaceOfBirth != (location.Location{}) {
		s.LocationService.PostLocation(&author.PlaceOfBirth)
	}

	var dateOfBirth string = ""
	if author.DateOfBirth != nil {
		dateOfBirth = *author.DateOfBirth
	}
	a := map[string]interface{}{
		"first_name":     author.FirstName,
		"last_name":      author.LastName,
		"gender":         author.Gender,
		"date_of_birth":  dateOfBirth,
		"place_of_birth": author.PlaceOfBirth.ID,
	}

	db := s.DB
	postAuthor := s.Query(PostAuthor, a)
	if err := db.QueryRow(postAuthor).Scan(&(author.ID)); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return status.New(status.Conflict, err.Error()), nil
		}
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), author
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
