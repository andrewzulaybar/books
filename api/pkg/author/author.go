package author

import (
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

// PostAuthor creates an entry in the author table with the given attributes.
func (s *Service) PostAuthor(author *Author) (*status.Status, *Author) {
	if author.PlaceOfBirth != (location.Location{}) {
		s.LocationService.PostLocation(&author.PlaceOfBirth)
	}

	db := s.DB
	a := map[string]interface{}{
		"first_name":     author.FirstName,
		"last_name":      author.LastName,
		"gender":         author.Gender,
		"date_of_birth":  author.DateOfBirth,
		"place_of_birth": author.PlaceOfBirth.ID,
	}

	postAuthor := s.Query(PostAuthor, a)
	if err := db.QueryRow(postAuthor).Scan(&(author.ID)); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return status.New(status.Conflict, err.Error()), nil
		}
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), author
}
