package author_test

import (
	"testing"

	"github.com/andrewzulaybar/books/api/internal/test/helpers"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/test/data"
)

func TestDeleteAuthor(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	authors := data.GetAuthors(as)
	helpers.PostAuthors(t, as, authors)

	t.Run("ExistingID", func(t *testing.T) {
		got := as.DeleteAuthor(1)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("AlreadyDeletedID", func(t *testing.T) {
		got := as.DeleteAuthor(2)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)

		got = as.DeleteAuthor(2)
		want = status.New(status.OK, "Author with id = 2 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		got := as.DeleteAuthor(-1)
		want := status.New(status.OK, "Author with id = -1 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	cleanup()
}

func TestDeleteAuthors(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	authors := data.GetAuthors(as)
	helpers.PostAuthors(t, as, authors)

	t.Run("OneID", func(t *testing.T) {
		gotStatus, gotIDs := as.DeleteAuthors([]int{1})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)
	})

	t.Run("MultipleIDs", func(t *testing.T) {
		gotStatus, gotIDs := as.DeleteAuthors([]int{2, 3})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)
	})

	t.Run("AlreadyDeletedIDs", func(t *testing.T) {
		gotStatus, gotIDs := as.DeleteAuthors([]int{4, 5})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)

		gotStatus, gotIDs = as.DeleteAuthors([]int{4, 5})
		wantStatus = status.New(status.OK, "The following authors could not be found: [4 5]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{4, 5}, gotIDs)
	})

	t.Run("IncludesNonExistentID", func(t *testing.T) {
		gotStatus, gotIDs := as.DeleteAuthors([]int{6, -1})
		wantStatus := status.New(status.OK, "The following authors could not be found: [-1]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{-1}, gotIDs)
	})

	t.Run("AllNonExistentIDs", func(t *testing.T) {
		gotStatus, gotIDs := as.DeleteAuthors([]int{-1, -2, -3})
		wantStatus := status.New(status.OK, "The following authors could not be found: [-1 -2 -3]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{-1, -2, -3}, gotIDs)
	})

	cleanup()
}

func TestFindAuthor(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	authors := data.GetAuthors(as)
	helpers.PostAuthors(t, as, authors)

	t.Run("ExistingAuthor", func(t *testing.T) {
		ta := &authors[0]
		gotStatus, gotAuthor := as.FindAuthor(ta.FirstName, ta.LastName, ta.DateOfBirth)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("NonExistentAuthor", func(t *testing.T) {
		ta := author.Author{FirstName: "A", LastName: "B", DateOfBirth: "1970-01-01T00:00:00Z"}
		gotStatus, gotAuthor := as.FindAuthor(ta.FirstName, ta.LastName, ta.DateOfBirth)
		msg := "Author ('A', 'B', '1970-01-01T00:00:00Z') does not exist"
		wantStatus := status.New(status.NotFound, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotAuthor)
	})

	cleanup()
}

func TestGetAuthor(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	authors := data.GetAuthors(as)
	helpers.PostAuthors(t, as, authors)

	t.Run("ExistingID", func(t *testing.T) {
		gotStatus, gotAuthor := as.GetAuthor(1)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, &authors[0], gotAuthor)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		gotStatus, gotAuthor := as.GetAuthor(-1)
		wantStatus := status.New(status.NotFound, "Author with id = -1 does not exist")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotAuthor)
	})

	cleanup()
}

func TestGetAuthors(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	tas := data.GetAuthors(as)
	authors := helpers.PostAuthors(t, as, tas)

	t.Run("AllAuthors", func(t *testing.T) {
		gotStatus, gotAuthors := as.GetAuthors()
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		for i := range authors {
			helpers.AssertEqual(t, authors[i], gotAuthors[i])
		}
	})

	cleanup()
}

func TestPatchAuthor(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	locations := data.GetLocations()
	authors := data.GetAuthors(as)

	t.Run("NoFieldsToUpdate", func(t *testing.T) {
		ta := helpers.PostAuthor(t, as, &authors[0])
		defer helpers.DeleteAuthor(t, as, ta.ID)

		updates := &author.Author{ID: ta.ID}

		gotStatus, gotAuthor := as.PatchAuthor(updates)
		wantStatus := status.New(status.BadRequest, "No fields in author to update")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotAuthor)
	})

	t.Run("FirstName", func(t *testing.T) {
		ta := helpers.PostAuthor(t, as, &authors[0])
		defer helpers.DeleteAuthor(t, as, ta.ID)

		firstName := "John"
		updates := &author.Author{ID: ta.ID, FirstName: firstName}

		gotStatus, gotAuthor := as.PatchAuthor(updates)
		wantStatus := status.New(status.OK, "")
		ta.FirstName = firstName
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("LastName", func(t *testing.T) {
		ta := helpers.PostAuthor(t, as, &authors[0])
		defer helpers.DeleteAuthor(t, as, ta.ID)

		lastName := "Steinbeck"
		updates := &author.Author{ID: ta.ID, LastName: lastName}

		gotStatus, gotAuthor := as.PatchAuthor(updates)
		wantStatus := status.New(status.OK, "")
		ta.LastName = lastName
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("Gender", func(t *testing.T) {
		ta := helpers.PostAuthor(t, as, &authors[0])
		defer helpers.DeleteAuthor(t, as, ta.ID)

		gender := "M"
		updates := &author.Author{ID: ta.ID, Gender: gender}

		gotStatus, gotAuthor := as.PatchAuthor(updates)
		wantStatus := status.New(status.OK, "")
		ta.Gender = gender
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("DateOfBirth", func(t *testing.T) {
		ta := helpers.PostAuthor(t, as, &authors[0])
		defer helpers.DeleteAuthor(t, as, ta.ID)

		dateOfBirth := "1970-01-01T00:00:00Z"
		updates := &author.Author{ID: ta.ID, DateOfBirth: dateOfBirth}

		gotStatus, gotAuthor := as.PatchAuthor(updates)
		wantStatus := status.New(status.OK, "")
		ta.DateOfBirth = dateOfBirth
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("PlaceOfBirth", func(t *testing.T) {
		t.Run("ExistingID", func(t *testing.T) {
			ta := helpers.PostAuthor(t, as, &authors[0])
			defer helpers.DeleteAuthor(t, as, ta.ID)

			tl := location.Location{ID: len(locations) - 1}
			updates := &author.Author{ID: ta.ID, PlaceOfBirth: tl}

			gotStatus, gotAuthor := as.PatchAuthor(updates)
			wantStatus := status.New(status.OK, "")
			ta.PlaceOfBirth.ID = tl.ID
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, ta, gotAuthor)
		})

		t.Run("NewCityCountry", func(t *testing.T) {
			ta := helpers.PostAuthor(t, as, &authors[0])
			defer helpers.DeleteAuthor(t, as, ta.ID)

			tl := location.Location{City: "A", Country: "B", Region: "C"}
			updates := &author.Author{ID: ta.ID, PlaceOfBirth: tl}

			gotStatus, gotAuthor := as.PatchAuthor(updates)
			wantStatus := status.New(status.OK, "")
			ta.PlaceOfBirth = location.Location{ID: len(locations) + 1}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, ta, gotAuthor)
		})

		t.Run("ExistingCityCountry", func(t *testing.T) {
			ta := helpers.PostAuthor(t, as, &authors[0])
			defer helpers.DeleteAuthor(t, as, ta.ID)

			tl := locations[0]
			updates := &author.Author{ID: ta.ID, PlaceOfBirth: tl}

			gotStatus, gotAuthor := as.PatchAuthor(updates)
			wantStatus := status.New(status.OK, "")
			ta.PlaceOfBirth = location.Location{ID: 1}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, ta, gotAuthor)
		})
	})

	t.Run("DuplicateFirstNameLastNameDateOfBirth", func(t *testing.T) {
		ta := helpers.PostAuthor(t, as, &authors[0])
		defer helpers.DeleteAuthor(t, as, ta.ID)
		dup := helpers.PostAuthor(t, as, &authors[1])
		defer helpers.DeleteAuthor(t, as, dup.ID)

		updates := &author.Author{
			ID:          ta.ID,
			FirstName:   dup.FirstName,
			LastName:    dup.LastName,
			DateOfBirth: dup.DateOfBirth,
		}

		gotStatus, gotAuthor := as.PatchAuthor(updates)
		msg := "pq: duplicate key value violates unique constraint \"author_first_name_last_name_date_of_birth_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotAuthor)
	})

	cleanup()
}

func TestPostAuthor(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ls := services.LocationService
	as := services.AuthorService
	authors := data.GetAuthors(as)

	t.Run("AllFieldsValid", func(t *testing.T) {
		ta := &authors[0]
		gotStatus, gotAuthor := as.PostAuthor(ta)
		defer helpers.DeleteAuthor(t, as, gotAuthor.ID)

		wantStatus := status.New(status.Created, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("ExistingID", func(t *testing.T) {
		ta := &authors[0]
		ga := helpers.PostAuthor(t, as, ta)
		defer helpers.DeleteAuthor(t, as, ga.ID)

		gotStatus, gotAuthor := as.PostAuthor(ga)
		wantStatus := status.Newf(status.Conflict, "Author with id = %d already exists", ga.ID)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("DateOfBirth", func(t *testing.T) {
		ta := &authors[0]
		ta.DateOfBirth = ""
		gotStatus, gotAuthor := as.PostAuthor(ta)
		defer helpers.DeleteAuthor(t, as, gotAuthor.ID)

		wantStatus := status.New(status.Created, "")
		ta.DateOfBirth = "1970-01-01T00:00:00Z"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, ta, gotAuthor)
	})

	t.Run("PlaceOfBirth", func(t *testing.T) {
		t.Run("ExistingID", func(t *testing.T) {
			ta := &authors[0]
			ta.PlaceOfBirth = location.Location{ID: 1}
			gotStatus, gotAuthor := as.PostAuthor(ta)
			defer helpers.DeleteAuthor(t, as, gotAuthor.ID)

			wantStatus := status.New(status.Created, "")
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, ta, gotAuthor)
		})

		t.Run("ExistingCityCountry", func(t *testing.T) {
			ta := &authors[0]
			ta.PlaceOfBirth = location.Location{City: "Toronto", Country: "Canada", Region: "Ontario"}
			gotStatus, gotAuthor := as.PostAuthor(ta)
			defer helpers.DeleteAuthor(t, as, gotAuthor.ID)

			wantStatus := status.New(status.Created, "")
			loc := helpers.FindLocation(t, ls, &ta.PlaceOfBirth, "Toronto", "Canada")
			ta.PlaceOfBirth = location.Location{ID: loc.ID}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, ta, gotAuthor)
		})

		t.Run("NewCityCountry", func(t *testing.T) {
			ta := &authors[0]
			ta.PlaceOfBirth = location.Location{City: "A", Country: "B", Region: "C"}
			gotStatus, gotAuthor := as.PostAuthor(ta)
			defer helpers.DeleteAuthor(t, as, gotAuthor.ID)

			wantStatus := status.New(status.Created, "")
			loc := helpers.FindLocation(t, ls, &ta.PlaceOfBirth, "A", "B")
			ta.PlaceOfBirth = location.Location{ID: loc.ID}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, ta, gotAuthor)
		})
	})

	t.Run("DuplicateFirstNameLastNameDateOfBirth", func(t *testing.T) {
		ta := &authors[0]
		ga := helpers.PostAuthor(t, as, ta)
		defer helpers.DeleteAuthor(t, as, ga.ID)

		gotStatus, gotAuthor := as.PostAuthor(ta)
		msg := "pq: duplicate key value violates unique constraint \"author_first_name_last_name_date_of_birth_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotAuthor)
	})

	cleanup()
}
