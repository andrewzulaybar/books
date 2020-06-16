package helpers

import (
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

// Services contains instances of services for testing.
type Services struct {
	LocationService    *location.Service
	AuthorService      *author.Service
	WorkService        *work.Service
	PublicationService *publication.Service
}

// Setup loads the test database and returns one instance of each service.
func Setup(t *testing.T) (*Services, func()) {
	t.Helper()

	conf, err := config.Load("test.env")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	db, dc := postgres.Setup(conf.ConnectionString)

	ls := &location.Service{DB: *db}
	as := &author.Service{DB: *db, LocationService: *ls}
	ws := &work.Service{DB: *db, AuthorService: *as}
	ps := &publication.Service{DB: *db, WorkService: *ws}

	services := &Services{
		LocationService:    ls,
		AuthorService:      as,
		WorkService:        ws,
		PublicationService: ps,
	}
	return services, dc
}
