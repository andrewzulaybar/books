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
			expected: status.DoesNotExist,
		},
		{
			name:     "ValidId",
			id:       1,
			expected: status.NoContent,
		},
		{
			name:     "AlreadyDeletedId",
			id:       1,
			expected: status.DoesNotExist,
		},
		{
			name:     "NotFoundId",
			id:       1000000,
			expected: status.DoesNotExist,
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
