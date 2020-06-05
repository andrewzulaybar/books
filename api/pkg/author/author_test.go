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

func TestDeleteAuthor(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadLocations, data.LoadAuthors}
	a, _ := setup(t, db, td)

	type Expected struct {
		status *status.Status
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
				status: status.New(status.NoContent, ""),
			},
		},
		{
			name: "AlreadyDeletedId",
			id:   1,
			expected: Expected{
				status: status.New(status.NotFound, "Author with id = 1 does not exist"),
			},
		},
		{
			name: "NegativeId",
			id:   -1,
			expected: Expected{
				status: status.New(status.NotFound, "Author with id = -1 does not exist"),
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s := a.DeleteAuthor(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
		})
	}
}

func TestDeleteAuthors(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadLocations, data.LoadAuthors}
	a, _ := setup(t, db, td)

	type Expected struct {
		status *status.Status
		ids    []int
	}

	cases := []struct {
		name     string
		ids      []int
		expected Expected
	}{
		{
			name: "OneId",
			ids:  []int{1},
			expected: Expected{
				status: status.New(status.NoContent, ""),
				ids:    nil,
			},
		},
		{
			name: "MultipleIds",
			ids:  []int{2, 3, 4},
			expected: Expected{
				status: status.New(status.NoContent, ""),
				ids:    nil,
			},
		},
		{
			name: "AlreadyDeletedIds",
			ids:  []int{1, 2, 3},
			expected: Expected{
				status: status.New(
					status.OK,
					"The following authors could not be found: [1 2 3]",
				),
				ids: []int{1, 2, 3},
			},
		},
		{
			name: "IncludesIdNotFound",
			ids:  []int{5, 6, -1},
			expected: Expected{
				status: status.New(
					status.OK,
					"The following authors could not be found: [-1]",
				),
				ids: []int{-1},
			},
		},
		{
			name: "AllIdsNotFound",
			ids:  []int{-1, -2, -3},
			expected: Expected{
				status: status.New(
					status.OK,
					"The following authors could not be found: [-1 -2 -3]",
				),
				ids: []int{-1, -2, -3},
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, ids := a.DeleteAuthors(c.ids)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.ids, ids) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.ids, ids)
			}
		})
	}
}

func TestGetAuthor(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadLocations, data.LoadAuthors}
	a, authors := setup(t, db, td)

	type Expected struct {
		status *status.Status
		author *author.Author
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
				status: status.New(status.OK, ""),
				author: &authors[0],
			},
		},
		{
			name: "InvalidId",
			id:   -1,
			expected: Expected{
				status: status.New(status.NotFound, "Author with id = -1 does not exist"),
				author: nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, author := a.GetAuthor(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.author, author) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.author, author)
			}
		})
	}
}

func TestGetAuthors(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadLocations, data.LoadAuthors}
	a, authors := setup(t, db, td)

	type Expected struct {
		status  *status.Status
		authors author.Authors
	}

	c := struct {
		name     string
		expected Expected
	}{
		"AllAuthors",
		Expected{
			status:  status.New(status.OK, ""),
			authors: authors,
		},
	}

	exp := c.expected
	t.Run(c.name, func(t *testing.T) {
		s, authors := a.GetAuthors()
		if !reflect.DeepEqual(exp.status, s) {
			t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
		}
		if !reflect.DeepEqual(exp.authors, authors) {
			t.Errorf("\nExpected: %v\nActual: %v\n", exp.authors, authors)
		}
	})
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
				authors[1].DateOfBirth = nil
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
