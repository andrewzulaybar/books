package location

import (
	"database/sql"
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
	GetLocation
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
	case GetLocation:
		return fmt.Sprintf("SELECT id, %s FROM location WHERE id = $1", Columns)
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

// GetLocation retrieves the location from the database matching the given id.
func (s *Service) GetLocation(id int) (*status.Status, *Location) {
	db := s.DB
	getLocation := s.Query(GetLocation)

	var location Location
	row := db.QueryRow(getLocation, id)
	if err := row.Scan(
		&location.ID,
		&location.City,
		&location.Country,
		&location.Region,
	); err != nil {
		if err == sql.ErrNoRows {
			return status.Newf(status.NotFound, "Location with id = %d does not exist", id), nil
		}
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), &location
}

// PostLocation creates an entry in the location table with the given attributes.
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
