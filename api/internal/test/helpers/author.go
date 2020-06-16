package helpers

import (
	"testing"

	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
)

// DeleteAuthor is a test helper that calls (AuthorService).DeleteAuthor.
func DeleteAuthor(t *testing.T, as *author.Service, id int) {
	t.Helper()

	gotStatus := as.DeleteAuthor(id)
	wantStatus := status.New(status.NoContent, "")
	AssertEqual(t, wantStatus, gotStatus)
}

// FindAuthor is a test helper that calls (AuthorService).FindAuthor with the given firstName, lastName, and dateOfBirth.
func FindAuthor(t *testing.T, as *author.Service, wantLocation location.Location,
	firstName string, lastName string, dateOfBirth string) *author.Author {
	t.Helper()

	gotStatus, gotAuthor := as.FindAuthor(firstName, lastName, dateOfBirth)
	wantStatus := status.New(status.OK, "")
	AssertEqual(t, gotStatus, wantStatus)
	AssertEqual(t, gotAuthor, wantLocation)

	return gotAuthor
}

// PostAuthor is a test helper that calls (AuthorService).PostAuthor.
func PostAuthor(t *testing.T, as *author.Service, ta *author.Author) *author.Author {
	t.Helper()

	gotStatus, gotAuthor := as.PostAuthor(ta)
	wantStatus := status.New(status.Created, "")
	AssertEqual(t, wantStatus, gotStatus)
	AssertEqual(t, ta, gotAuthor)

	return gotAuthor
}

// PostAuthors is a test helper that calls (AuthorService).PostAuthors for each author in `authors`.
func PostAuthors(t *testing.T, as *author.Service, ta author.Authors) author.Authors {
	t.Helper()

	var authors author.Authors
	for i := range ta {
		author := PostAuthor(t, as, &ta[i])
		if author != nil {
			authors = append(authors, *author)
		}
	}
	return authors
}
