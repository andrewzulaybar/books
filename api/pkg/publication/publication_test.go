package publication_test

import (
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

var pubs publication.Publications = publication.Publications{
	{
		1,
		"2019-04-16T00:00:00Z",
		"Hardcover",
		"https://images-na.ssl-images-amazon.com/images/I/81X4R7QhFkL.jpg",
		"1984822179",
		"9781984822178",
		"English",
		288,
		"Hogarth",
		work.Work{},
	},
	{
		2,
		"2017-09-12T00:00:00Z",
		"Hardcover",
		"https://images-na.ssl-images-amazon.com/images/I/91twTG-CQ8L.jpg",
		"0735224293",
		"9780735224292",
		"English",
		352,
		"Penguin Press",
		work.Work{},
	},
	{
		3,
		"2018-08-14T00:00:00Z",
		"Hardcover",
		"https://images-na.ssl-images-amazon.com/images/I/51j5p18mJNL.jpg",
		"0735219095",
		"9780735219090",
		"English",
		384,
		"G.P. Putnam's Sons",
		work.Work{},
	},
	{
		4,
		"2004-09-30T00:00:00Z",
		"Paperback",
		"https://images-na.ssl-images-amazon.com/images/I/81af+MCATTL.jpg",
		"0743273567",
		"9780743273565",
		"English",
		180,
		"Scribner",
		work.Work{},
	},
	{
		5,
		"2020-01-21T00:00:00Z",
		"Hardcover",
		"https://images-na.ssl-images-amazon.com/images/I/81iVsj91eQL.jpg",
		"1250209765",
		"9781250209764",
		"English",
		400,
		"Flatiron Books",
		work.Work{},
	},
	{
		6,
		"2018-09-16T00:00:00Z",
		"Hardcover",
		"https://images-na.ssl-images-amazon.com/images/I/91Xq+S+F2jL.jpg",
		"0735211299",
		"9780735211292",
		"English",
		320,
		"Avery",
		work.Work{},
	},
}

func getDB() *postgres.DB {
	conf, err := config.Load("../../config/test.env")
	if err != nil {
		panic(err)
	}

	return postgres.Setup(conf.ConnectionString, "../../internal/sql/")
}

func TestDeletePublication(t *testing.T) {
	db := getDB()
	defer db.Disconnect()

	p := &publication.Service{DB: *db}

	cases := []struct {
		name     string
		id       int
		expected int
	}{
		{
			name:     "NegativeId",
			id:       -1,
			expected: status.NotFound,
		},
		{
			name:     "ValidId",
			id:       1,
			expected: status.NoContent,
		},
		{
			name:     "AlreadyDeletedId",
			id:       1,
			expected: status.NotFound,
		},
		{
			name:     "IdNotFound",
			id:       1000000,
			expected: status.NotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := p.DeletePublication(c.id)
			if code := s.Code(); code != c.expected {
				t.Errorf("\nExpected: %d\nActual: %d\n", c.expected, code)
			}
		})
	}
}

func TestDeletePublications(t *testing.T) {
	db := getDB()
	defer db.Disconnect()

	p := &publication.Service{DB: *db}

	type Expected struct {
		code int
		ids  []int
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
				code: status.NoContent,
				ids:  nil,
			},
		},
		{
			name: "MultipleIds",
			ids:  []int{2, 3, 4},
			expected: Expected{
				code: status.NoContent,
				ids:  nil,
			},
		},
		{
			name: "AlreadyDeletedIds",
			ids:  []int{1, 2, 3},
			expected: Expected{
				code: status.OK,
				ids:  []int{1, 2, 3},
			},
		},
		{
			name: "IncludesIdNotFound",
			ids:  []int{5, 6, -1},
			expected: Expected{
				code: status.OK,
				ids:  []int{-1},
			},
		},
		{
			name: "AllIdsNotFound",
			ids:  []int{-1, -2, -3},
			expected: Expected{
				code: status.OK,
				ids:  []int{-1, -2, -3},
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, ids := p.DeletePublications(c.ids)
			if code := s.Code(); code != exp.code {
				t.Errorf("\nExpected: %d\nActual: %d\n", exp.code, code)
			}
			for i, id := range ids {
				if id != exp.ids[i] {
					t.Errorf("\nExpected: %d\nActual: %d\n", exp.ids[i], id)
				}
			}
		})
	}
}

func TestGetPublication(t *testing.T) {
	db := getDB()
	defer db.Disconnect()

	p := &publication.Service{DB: *db}

	type Expected struct {
		code int
		pub  *publication.Publication
	}

	cases := []struct {
		name     string
		id       int
		expected Expected
	}{
		{
			"ValidId",
			1,
			Expected{
				code: status.OK,
				pub:  &pubs[0],
			},
		},
		{
			"InvalidId",
			-1,
			Expected{
				code: status.NotFound,
				pub:  nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, pub := p.GetPublication(c.id)
			if code := s.Code(); code != exp.code {
				t.Errorf("\nExpected: %d\nActual: %d\n", exp.code, code)
			}
			if exp.pub != nil && pub != nil {
				pub.Work = work.Work{}
				if *exp.pub != *pub {
					t.Errorf("\nExpected: %v\nActual: %v\n", *exp.pub, *pub)
				}
			} else if exp.pub != pub {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.pub, pub)
			}
		})
	}
}

func TestGetPublications(t *testing.T) {
	db := getDB()
	defer db.Disconnect()

	p := &publication.Service{DB: *db}

	type Expected struct {
		code int
		pubs publication.Publications
	}

	c := struct {
		name     string
		expected Expected
	}{
		"AllPublications",
		Expected{
			code: status.OK,
			pubs: pubs,
		},
	}

	exp := c.expected
	t.Run(c.name, func(t *testing.T) {
		s, pubs := p.GetPublications()
		if code := s.Code(); code != exp.code {
			t.Errorf("\nExpected: %d\nActual: %d\n", exp.code, code)
		}
		if pubs != nil {
			for i, pub := range pubs {
				pub.Work = work.Work{}
				if pub != exp.pubs[i] {
					t.Errorf("\nExpected: %v\nActual: %v\n", exp.pubs[i], pub)
				}
			}
		}
	})
}
