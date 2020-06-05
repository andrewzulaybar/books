package location_test

import (
	"reflect"
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/test/data"
)

func getDB(t *testing.T) (*postgres.DB, func()) {
	t.Helper()

	conf, err := config.Load("../../config/test.env")
	if err != nil {
		panic(err)
	}

	return postgres.Setup(conf.ConnectionString, "../../internal/sql/")
}

func TestGetWork(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	l := &location.Service{DB: *db}
	locations := data.LoadLocations(l)

	type Expected struct {
		status   *status.Status
		location *location.Location
	}

	cases := []struct {
		name     string
		id       int
		expected Expected
	}{
		{
			name: "ValidId",
			id:   1,
			expected: Expected{
				status:   status.New(status.OK, ""),
				location: &locations[0],
			},
		},
		{
			name: "InvalidId",
			id:   -1,
			expected: Expected{
				status:   status.New(status.NotFound, "Location with id = -1 does not exist"),
				location: nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, loc := l.GetLocation(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.location, loc) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.location, loc)
			}
		})
	}
}

func TestPostLocation(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	l := &location.Service{DB: *db}
	locations := data.GetLocations()

	type Expected struct {
		status *status.Status
		loc    *location.Location
	}

	cases := []struct {
		name     string
		loc      *location.Location
		expected Expected
	}{
		{
			name: "AllFieldsValid",
			loc:  &locations[0],
			expected: Expected{
				status: status.New(status.Created, ""),
				loc:    &locations[0],
			},
		},
		{
			name: "DuplicateCityCountry",
			loc: func() *location.Location {
				loc := locations[1]
				loc.City = locations[0].City
				loc.Country = locations[0].Country
				return &loc
			}(),
			expected: Expected{
				status: status.New(
					status.Conflict,
					"pq: duplicate key value violates unique constraint \"location_city_country_key\"",
				),
				loc: nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, loc := l.PostLocation(c.loc)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.loc, loc) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.loc, loc)
			}
		})
	}
}
