package author_test

import (
	"reflect"
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/test/data"
)

type TestData struct {
	GetLocations func(*location.Service) location.Locations
	GetAuthors   func(*author.Service) author.Authors
}

func getDB(t *testing.T) (*postgres.DB, func()) {
	t.Helper()

	conf, err := config.Load("../../config/test.env")
	if err != nil {
		panic(err)
	}

	return postgres.Setup(conf.ConnectionString, "../../internal/sql/")
}

func setup(t *testing.T, db *postgres.DB, td *TestData) (*author.Service, author.Authors) {
	t.Helper()

	l := &location.Service{DB: *db}
	a := &author.Service{
		DB:              *db,
		LocationService: *l,
	}
	locations := td.GetLocations(l)
	authors := td.GetAuthors(a)

	for i := range authors {
		for j := range locations {
			if authors[i].PlaceOfBirth.ID == locations[j].ID {
				authors[i].PlaceOfBirth = locations[j]
			}
		}
	}

	return a, authors
}

func TestPostAuthor(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadLocations, data.GetAuthors}
	a, authors := setup(t, db, td)

	type Expected struct {
		status *status.Status
		author *author.Author
	}

	cases := []struct {
		name     string
		author   *author.Author
		expected Expected
	}{
		{
			name:   "AllFieldsValid",
			author: &authors[0],
			expected: Expected{
				status: status.New(status.Created, ""),
				author: &authors[0],
			},
		},
		{
			name: "OnlyRequiredFields",
			author: func() *author.Author {
				authors[1].DateOfBirth = ""
				authors[1].PlaceOfBirth = location.Location{}
				return &authors[1]
			}(),
			expected: Expected{
				status: status.New(status.Created, ""),
				author: &authors[1],
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, author := a.PostAuthor(c.author)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.author, author) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.author, author)
			}
		})
	}
}
