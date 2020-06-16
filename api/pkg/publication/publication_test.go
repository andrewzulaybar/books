package publication_test

import (
	"testing"

	"github.com/andrewzulaybar/books/api/internal/test/helpers"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/andrewzulaybar/books/api/test/data"
)

func TestDeletePublication(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ps := services.PublicationService
	publications := data.GetPublications(ps)
	helpers.PostPublications(t, ps, publications)

	t.Run("ValidID", func(t *testing.T) {
		got := ps.DeletePublication(1)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("AlreadyDeletedID", func(t *testing.T) {
		got := ps.DeletePublication(2)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)

		got = ps.DeletePublication(2)
		want = status.New(status.OK, "Publication with id = 2 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		got := ps.DeletePublication(-1)
		want := status.New(status.OK, "Publication with id = -1 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	cleanup()
}

func TestDeletePublications(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ps := services.PublicationService
	publications := data.GetPublications(ps)
	helpers.PostPublications(t, ps, publications)

	t.Run("OneID", func(t *testing.T) {
		gotStatus, gotIDs := ps.DeletePublications([]int{1})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)
	})

	t.Run("MultipleIDs", func(t *testing.T) {
		gotStatus, gotIDs := ps.DeletePublications([]int{2, 3})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)
	})

	t.Run("AlreadyDeletedIDs", func(t *testing.T) {
		gotStatus, gotIDs := ps.DeletePublications([]int{4, 5})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)

		gotStatus, gotIDs = ps.DeletePublications([]int{4, 5})
		wantStatus = status.New(status.OK, "The following publications could not be found: [4 5]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{4, 5}, gotIDs)
	})

	t.Run("IncludesNonExistentID", func(t *testing.T) {
		gotStatus, gotIDs := ps.DeletePublications([]int{6, -1})
		wantStatus := status.New(status.OK, "The following publications could not be found: [-1]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{-1}, gotIDs)
	})

	t.Run("AllNonExistentIDs", func(t *testing.T) {
		gotStatus, gotIDs := ps.DeletePublications([]int{1, 2, 3})
		wantStatus := status.New(status.OK, "The following publications could not be found: [1 2 3]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{1, 2, 3}, gotIDs)
	})

	cleanup()
}

func TestGetPublication(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ps := services.PublicationService
	publications := data.GetPublications(ps)
	helpers.PostPublications(t, ps, publications)

	t.Run("ValidID", func(t *testing.T) {
		gotStatus, gotWork := ps.GetPublication(1)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, &publications[0], gotWork)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		gotStatus, gotWork := ps.GetPublication(-1)
		wantStatus := status.New(status.NotFound, "Publication with id = -1 does not exist")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	cleanup()
}

func TestGetPublications(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ps := services.PublicationService
	publications := data.GetPublications(ps)
	helpers.PostPublications(t, ps, publications)

	t.Run("AllWorks", func(t *testing.T) {
		gotStatus, gotPublications := ps.GetPublications()
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		for i := range publications {
			helpers.AssertEqual(t, publications[i], gotPublications[i])
		}
	})

	cleanup()
}

func TestPatchPublication(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ps := services.PublicationService
	works := data.GetWorks(&ps.WorkService)
	publications := data.GetPublications(ps)

	t.Run("NoFieldsToUpdate", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "No fields in publication to update")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	t.Run("EditionPubDate", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, EditionPubDate: "1900-01-01T00:00:00Z"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.EditionPubDate = "1900-01-01T00:00:00Z"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("InitialPubDate", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, Format: "Paperback"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.Format = "Paperback"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("ImageUrl", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, ImageURL: "https://google.com"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.ImageURL = "https://google.com"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("ISBN", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, ISBN: "0123456789"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.ISBN = "0123456789"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("ISBN13", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, ISBN13: "0123456789123"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.ISBN13 = "0123456789123"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("Language", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, Language: "French"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.Language = "French"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("NumPages", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, NumPages: 100}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.NumPages = 100
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("Publisher", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		updates := &publication.Publication{ID: tp.ID, Publisher: "Penguin"}

		gotStatus, gotWork := ps.PatchPublication(updates)
		wantStatus := status.New(status.OK, "")
		tp.Publisher = "Penguin"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotWork)
	})

	t.Run("Work", func(t *testing.T) {
		t.Run("ExistingID", func(t *testing.T) {
			tp := helpers.PostPublication(t, ps, &publications[0])
			defer helpers.DeletePublication(t, ps, tp.ID)

			tw := work.Work{ID: len(works) - 1}
			updates := &publication.Publication{ID: tp.ID, Work: tw}

			gotStatus, gotPublication := ps.PatchPublication(updates)
			wantStatus := status.New(status.OK, "")
			tp.Work.ID = tw.ID
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tp, gotPublication)
		})

		t.Run("NewTitleAuthorID", func(t *testing.T) {
			tp := helpers.PostPublication(t, ps, &publications[0])
			defer helpers.DeletePublication(t, ps, tp.ID)

			ta := author.Author{ID: 1}
			tw := work.Work{Title: "A", InitialPubDate: "1900-01-01T00:00:00Z", Author: ta}
			updates := &publication.Publication{ID: tp.ID, Work: tw}

			gotStatus, gotPublication := ps.PatchPublication(updates)
			wantStatus := status.New(status.OK, "")
			tp.Work = work.Work{ID: len(works) + 1}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tp, gotPublication)
		})

		t.Run("ExistingTitleAuthorID", func(t *testing.T) {
			tp := helpers.PostPublication(t, ps, &publications[0])
			defer helpers.DeletePublication(t, ps, tp.ID)

			tw := works[0]
			updates := &publication.Publication{ID: tp.ID, Work: tw}

			gotStatus, gotPublication := ps.PatchPublication(updates)
			wantStatus := status.New(status.OK, "")
			tp.Work = work.Work{ID: 1}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tp, gotPublication)
		})
	})

	t.Run("DuplicateImageURL", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)
		dup := helpers.PostPublication(t, ps, &publications[1])
		defer helpers.DeletePublication(t, ps, dup.ID)

		updates := &publication.Publication{ID: tp.ID, ImageURL: dup.ImageURL}
		gotStatus, gotPublication := ps.PatchPublication(updates)
		msg := "pq: duplicate key value violates unique constraint \"publication_image_url_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotPublication)
	})

	t.Run("DuplicateISBN", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)
		dup := helpers.PostPublication(t, ps, &publications[1])
		defer helpers.DeletePublication(t, ps, dup.ID)

		updates := &publication.Publication{ID: tp.ID, ISBN: dup.ISBN}
		gotStatus, gotPublication := ps.PatchPublication(updates)
		msg := "pq: duplicate key value violates unique constraint \"publication_isbn_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotPublication)
	})

	t.Run("DuplicateISBN13", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)
		dup := helpers.PostPublication(t, ps, &publications[1])
		defer helpers.DeletePublication(t, ps, dup.ID)

		updates := &publication.Publication{ID: tp.ID, ISBN13: dup.ISBN13}
		gotStatus, gotPublication := ps.PatchPublication(updates)
		msg := "pq: duplicate key value violates unique constraint \"publication_isbn13_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotPublication)
	})

	cleanup()
}

func TestPostPublication(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	ps := services.PublicationService
	works := data.GetWorks(&ps.WorkService)
	publications := data.GetPublications(ps)

	t.Run("AllFieldsValid", func(t *testing.T) {
		tp := &publications[0]
		gotStatus, gotPublication := ps.PostPublication(tp)
		defer helpers.DeletePublication(t, ps, gotPublication.ID)

		wantStatus := status.New(status.Created, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tp, gotPublication)
	})

	t.Run("Work", func(t *testing.T) {
		t.Run("ExistingID", func(t *testing.T) {
			tp := &publications[0]
			tp.Work = work.Work{ID: 1}
			gotStatus, gotPublication := ps.PostPublication(tp)
			defer helpers.DeletePublication(t, ps, gotPublication.ID)

			wantStatus := status.New(status.Created, "")
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tp, gotPublication)
		})

		t.Run("ExistingTitleAuthorID", func(t *testing.T) {
			tp := &publications[0]
			tp.Work = works[0]
			gotStatus, gotPublication := ps.PostPublication(tp)
			defer helpers.DeletePublication(t, ps, gotPublication.ID)

			wantStatus := status.New(status.Created, "")
			wk := helpers.FindWork(t, ws, &tp.Work, tp.Work.Title, tp.Work.Author.ID)
			tp.Work = work.Work{ID: wk.ID}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tp, gotPublication)
		})

		t.Run("NewTitleAuthorID", func(t *testing.T) {
			tp := &publications[0]
			tp.Work = work.Work{
				Description:      "A",
				InitialPubDate:   "1900-01-01T00:00:00Z",
				OriginalLanguage: "B",
				Title:            "C",
				Author: author.Author{
					FirstName:    "D",
					LastName:     "E",
					Gender:       "F",
					DateOfBirth:  "1900-01-01T00:00:00Z",
					PlaceOfBirth: location.Location{City: "G", Country: "H", Region: "I"},
				},
			}
			gotStatus, gotPublication := ps.PostPublication(tp)
			defer helpers.DeletePublication(t, ps, gotPublication.ID)

			wantStatus := status.New(status.Created, "")
			wk := helpers.FindWork(t, ws, &tp.Work, tp.Work.Title, tp.Work.Author.ID)
			tp.Work = work.Work{ID: wk.ID}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tp, gotPublication)
		})
	})

	t.Run("DuplicateImageURL", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		dup := publications[1]
		dup.ImageURL = tp.ImageURL
		gotStatus, gotPublication := ps.PostPublication(&dup)

		msg := "pq: duplicate key value violates unique constraint \"publication_image_url_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotPublication)
	})

	t.Run("DuplicateISBN", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		dup := publications[1]
		dup.ISBN = tp.ISBN
		gotStatus, gotPublication := ps.PostPublication(&dup)

		msg := "pq: duplicate key value violates unique constraint \"publication_isbn_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotPublication)
	})

	t.Run("DuplicateISBN13", func(t *testing.T) {
		tp := helpers.PostPublication(t, ps, &publications[0])
		defer helpers.DeletePublication(t, ps, tp.ID)

		dup := publications[1]
		dup.ISBN13 = tp.ISBN13
		gotStatus, gotPublication := ps.PostPublication(&dup)

		msg := "pq: duplicate key value violates unique constraint \"publication_isbn13_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotPublication)
	})

	cleanup()
}
