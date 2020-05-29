package publication_test

import (
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
)

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
			for i, val := range exp.ids {
				if ids[i] != val {
					t.Errorf("\nExpected: %d\nActual: %d\n", val, ids[i])
				}
			}
		})
	}
}
