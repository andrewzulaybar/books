package location

import (
	"fmt"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/lib/pq"
)

// Columns is the comma-separated list of columns found in the location table.
const Columns string = "city, country, region"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	PostLocation
)

// Location represents a geographic city location.
type Location struct {
	ID      int    `json:"id"`
	City    string `json:"city"`
	Country string `json:"country"`
	Region  string `json:"region"`
}

// Locations represents a list of locations.
type Locations []Location

// Service wraps the database and other dependencies.
type Service struct {
	DB postgres.DB
}

// Query returns a SQL statement based on the postgres.Query value passed in.
func (s *Service) Query(query postgres.Query, args ...interface{}) string {
	switch query {
	case PostLocation:
		return fmt.Sprintf(
			`INSERT INTO location (%s)
                        VALUES ($1, $2, $3)
                        RETURNING id`,
			Columns,
		)
	default:
		return ""
	}
}

func (s *Service) PostLocation(loc *Location) (*status.Status, *Location) {
	db := s.DB
	if err := db.QueryRow(
		s.Query(PostLocation),
		loc.City,
		loc.Country,
		loc.Region,
	).Scan(&(loc.ID)); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return status.New(status.Conflict, err.Error()), nil
		}
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), loc
}
