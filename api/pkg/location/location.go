package location

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/lib/pq"
)

// Columns is the comma-separated list of columns found in the location table.
const Columns string = "city, country, region"

// Enum constants representing types of SQL statements.
const (
	Unknown postgres.Query = iota
	DeleteLocation
	FindLocation
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
	case DeleteLocation:
		return "DELETE FROM location WHERE ID = $1"
	case FindLocation:
		city := args[0].(string)
		country := args[1].(string)
		return fmt.Sprintf(
			`SELECT id, %s
                        FROM location
                        WHERE city = '%s' AND country = '%s'`,
			Columns, city, country,
		)
	case GetLocation:
		return fmt.Sprintf("SELECT id, %s FROM location WHERE id = $1", Columns)
	case PostLocation:
		return fmt.Sprintf(
			`INSERT INTO location (%s)
                        VALUES ($1, $2, $3)
                        RETURNING *`,
			Columns,
		)
	default:
		return ""
	}
}

// DeleteLocation removes the entry in the location table matching the given id.
func (s *Service) DeleteLocation(id int) *status.Status {
	db := s.DB
	deleteLocation := s.Query(DeleteLocation)

	res, err := db.Exec(deleteLocation, id)
	if err != nil {
		log.Printf("[DeleteLocation] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		log.Printf("[DeleteLocation] %s", err)
		return status.New(status.InternalServerError, err.Error())
	}
	if numDeleted == 0 {
		msg := fmt.Sprintf("Location with id = %d does not exist", id)
		log.Printf("[DeleteLocation] %s", msg)
		return status.New(status.OK, msg)
	}
	return status.New(status.NoContent, "")
}

// FindLocation retrieves the location from the database matching the given city and country.
func (s *Service) FindLocation(city string, country string) (*status.Status, *Location) {
	db := s.DB
	findLocation := s.Query(FindLocation, city, country)

	var location Location
	row := db.QueryRow(findLocation)
	if err := row.Scan(&location.ID, &location.City, &location.Country, &location.Region); err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Location ('%s', '%s') does not exist", city, country)
			log.Printf("[FindLocation] %s", msg)
			return status.Newf(status.NotFound, msg), nil
		}
		log.Printf("[FindLocation] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), &location
}

// GetLocation retrieves the location from the database matching the given id.
func (s *Service) GetLocation(id int) (*status.Status, *Location) {
	db := s.DB
	getLocation := s.Query(GetLocation)

	var location Location
	row := db.QueryRow(getLocation, id)
	if err := row.Scan(&location.ID, &location.City, &location.Country, &location.Region); err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Location with id = %d does not exist", id)
			log.Printf("[GetLocation] %s", msg)
			return status.Newf(status.NotFound, msg), nil
		}
		log.Printf("[GetLocation] %s", err)
		return status.New(status.InternalServerError, err.Error()), nil
	}
	return status.New(status.OK, ""), &location
}

// PostLocation creates an entry in the location table with the given attributes.
func (s *Service) PostLocation(location *Location) (*status.Status, *Location) {
	if location.ID != 0 {
		if s, l := s.GetLocation(location.ID); s.Code() == status.OK {
			msg := fmt.Sprintf("Location with id = %d already exists", location.ID)
			log.Printf("[PostLocation] %s", msg)
			return status.Newf(status.Conflict, msg), l
		}
	}

	db := s.DB
	postLocation := s.Query(PostLocation)

	var loc Location
	row := db.QueryRow(postLocation, location.City, location.Country, location.Region)
	if err := row.Scan(&loc.ID, &loc.City, &loc.Country, &loc.Region); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			_, l := s.FindLocation(location.City, location.Country)
			log.Printf("[PostLocation] %s", err)
			return status.New(status.Conflict, err.Error()), l
		}
		log.Printf("[PostLocation] %s", err)
		return status.New(status.UnprocessableEntity, err.Error()), nil
	}
	return status.New(status.Created, ""), &loc
}
