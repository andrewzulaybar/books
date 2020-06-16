package helpers

import (
	"testing"

	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
)

// DeletePublication is a test helper that calls (PublicationService).DeletePublication.
func DeletePublication(t *testing.T, ps *publication.Service, id int) {
	t.Helper()

	gotStatus := ps.DeletePublication(id)
	wantStatus := status.New(status.NoContent, "")
	AssertEqual(t, wantStatus, gotStatus)
}

// PostPublication is a test helper that calls (PublicationService).PostPublication.
func PostPublication(t *testing.T, ps *publication.Service, tp *publication.Publication) *publication.Publication {
	t.Helper()

	gotStatus, gotPublication := ps.PostPublication(tp)
	wantStatus := status.New(status.Created, "")
	AssertEqual(t, wantStatus, gotStatus)
	AssertEqual(t, tp, gotPublication)

	return gotPublication
}

// PostPublications is a test helper that calls (PublicationService).PostPublications for each publication in `publications`.
func PostPublications(t *testing.T, ps *publication.Service, tp publication.Publications) publication.Publications {
	t.Helper()

	var publications publication.Publications
	for i := range tp {
		publication := PostPublication(t, ps, &tp[i])
		if publication != nil {
			publications = append(publications, *publication)
		}
	}
	return publications
}
