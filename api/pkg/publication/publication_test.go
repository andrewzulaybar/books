package publication_test

import (
	"reflect"
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/andrewzulaybar/books/api/test/data"
)

type TestData struct {
	GetWorks        func(*work.Service) work.Works
	GetPublications func(*publication.Service) publication.Publications
}

func getDB(t *testing.T) (*postgres.DB, func()) {
	t.Helper()

	conf, err := config.Load("../../config/test.env")
	if err != nil {
		panic(err)
	}

	return postgres.Setup(conf.ConnectionString, "../../internal/sql/")
}

func setup(t *testing.T, db *postgres.DB, td *TestData) (*publication.Service, publication.Publications) {
	t.Helper()

	l := &location.Service{DB: *db}
	a := &author.Service{
		DB:              *db,
		LocationService: *l,
	}
	data.LoadLocations(l)
	data.LoadAuthors(a)

	w := &work.Service{DB: *db}
	p := &publication.Service{
		DB:          *db,
		WorkService: *w,
	}
	works := td.GetWorks(w)
	publications := td.GetPublications(p)

	for i := range publications {
		for j := range works {
			if publications[i].Work.ID == works[j].ID {
				publications[i].Work = works[j]
			}
		}
	}

	return p, publications
}

func TestDeletePublication(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadWorks, data.LoadPublications}
	p, _ := setup(t, db, td)

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
				status: status.New(status.NotFound, "Publication with id = 1 does not exist"),
			},
		},
		{
			name: "NegativeId",
			id:   -1,
			expected: Expected{
				status: status.New(status.NotFound, "Publication with id = -1 does not exist"),
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s := p.DeletePublication(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
		})
	}
}

func TestDeletePublications(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadWorks, data.LoadPublications}
	p, _ := setup(t, db, td)

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
					"The following publications could not be found: [1 2 3]",
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
					"The following publications could not be found: [-1]",
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
					"The following publications could not be found: [-1 -2 -3]",
				),
				ids: []int{-1, -2, -3},
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, ids := p.DeletePublications(c.ids)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.ids, ids) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.ids, ids)
			}
		})
	}
}

func TestGetPublication(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadWorks, data.LoadPublications}
	p, pubs := setup(t, db, td)

	type Expected struct {
		status *status.Status
		pub    *publication.Publication
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
				pub:    &pubs[0],
			},
		},
		{
			name: "InvalidId",
			id:   -1,
			expected: Expected{
				status: status.New(status.NotFound, "Publication with id = -1 does not exist"),
				pub:    nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, pub := p.GetPublication(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.pub, pub) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.pub, pub)
			}
		})
	}
}

func TestGetPublications(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadWorks, data.LoadPublications}
	p, pubs := setup(t, db, td)

	type Expected struct {
		status *status.Status
		pubs   publication.Publications
	}

	c := struct {
		name     string
		expected Expected
	}{
		"AllPublications",
		Expected{
			status: status.New(status.OK, ""),
			pubs:   pubs,
		},
	}

	exp := c.expected
	t.Run(c.name, func(t *testing.T) {
		s, pubs := p.GetPublications()
		if !reflect.DeepEqual(exp.status, s) {
			t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
		}
		if !reflect.DeepEqual(exp.pubs, pubs) {
			t.Errorf("\nExpected: %v\nActual: %v\n", exp.pubs, pubs)
		}
	})
}

func TestPatchPublication(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadWorks, data.LoadPublications}
	p, _ := setup(t, db, td)

	type Expected struct {
		status *status.Status
		id     int
	}

	cases := []struct {
		name     string
		pub      *publication.Publication
		expected Expected
	}{
		{
			name: "PatchEditionPubDate",
			pub: &publication.Publication{
				ID:             1,
				EditionPubDate: "1900-01-01T00:00:00Z",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchFormat",
			pub: &publication.Publication{
				ID:     1,
				Format: "Paperback",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchImageURL",
			pub: &publication.Publication{
				ID:       1,
				ImageURL: "https://example.com",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchISBN",
			pub: &publication.Publication{
				ID:   1,
				ISBN: "0123456789",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchISBN13",
			pub: &publication.Publication{
				ID:     1,
				ISBN13: "0123456789012",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchLanguage",
			pub: &publication.Publication{
				ID:       1,
				Language: "French",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchNumPages",
			pub: &publication.Publication{
				ID:       1,
				NumPages: 100,
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "PatchPublisher",
			pub: &publication.Publication{
				ID:        1,
				Publisher: "Penguin",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, pub := p.PatchPublication(c.pub)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			_, want := p.GetPublication(exp.id)
			if !reflect.DeepEqual(want, pub) {
				t.Errorf("\nExpected: %v\nActual: %v\n", want, pub)
			}
		})
	}
}

func TestPostPublication(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	td := &TestData{data.LoadWorks, data.GetPublications}
	p, pubs := setup(t, db, td)

	type Expected struct {
		status *status.Status
		pub    *publication.Publication
	}

	cases := []struct {
		name     string
		pub      *publication.Publication
		expected Expected
	}{
		{
			name: "AllFieldsValid",
			pub:  &pubs[0],
			expected: Expected{
				status: status.New(status.Created, ""),
				pub:    &pubs[0],
			},
		},
		{
			name: "DuplicateImageURL",
			pub: func() *publication.Publication {
				pub := pubs[1]
				pub.ImageURL = pubs[0].ImageURL
				return &pub
			}(),
			expected: Expected{
				status: status.New(
					status.Conflict,
					"pq: duplicate key value violates unique constraint \"publication_image_url_key\"",
				),
				pub: nil,
			},
		},
		{
			name: "DuplicateISBN",
			pub: func() *publication.Publication {
				pub := pubs[1]
				pub.ISBN = pubs[0].ISBN
				return &pub
			}(),
			expected: Expected{
				status: status.New(
					status.Conflict,
					"pq: duplicate key value violates unique constraint \"publication_isbn_key\"",
				),
				pub: nil,
			},
		},
		{
			name: "DuplicateISBN13",
			pub: func() *publication.Publication {
				pub := pubs[1]
				pub.ISBN13 = pubs[0].ISBN13
				return &pub
			}(),
			expected: Expected{
				status: status.New(
					status.Conflict,
					"pq: duplicate key value violates unique constraint \"publication_isbn13_key\"",
				),
				pub: nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, pub := p.PostPublication(c.pub)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.pub, pub) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.pub, pub)
			}
		})
	}
}
