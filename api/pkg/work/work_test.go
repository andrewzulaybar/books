package work_test

import (
	"testing"

	"github.com/andrewzulaybar/books/api/internal/test/helpers"
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/andrewzulaybar/books/api/test/data"
)

func TestDeleteWork(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	works := data.GetWorks(ws)
	helpers.PostWorks(t, ws, works)

	t.Run("ValidID", func(t *testing.T) {
		got := ws.DeleteWork(1)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("AlreadyDeletedID", func(t *testing.T) {
		got := ws.DeleteWork(2)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)

		got = ws.DeleteWork(2)
		want = status.New(status.OK, "Work with id = 2 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		got := ws.DeleteWork(-1)
		want := status.New(status.OK, "Work with id = -1 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	cleanup()
}

func TestDeleteWorks(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	works := data.GetWorks(ws)
	helpers.PostWorks(t, ws, works)

	t.Run("OneID", func(t *testing.T) {
		gotStatus, gotIDs := ws.DeleteWorks([]int{1})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)
	})

	t.Run("MultipleIDs", func(t *testing.T) {
		gotStatus, gotIDs := ws.DeleteWorks([]int{2, 3})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)
	})

	t.Run("AlreadyDeletedIDs", func(t *testing.T) {
		gotStatus, gotIDs := ws.DeleteWorks([]int{4, 5})
		wantStatus := status.New(status.NoContent, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotIDs)

		gotStatus, gotIDs = ws.DeleteWorks([]int{4, 5})
		wantStatus = status.New(status.OK, "The following works could not be found: [4 5]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{4, 5}, gotIDs)
	})

	t.Run("IncludesNonExistentID", func(t *testing.T) {
		gotStatus, gotIDs := ws.DeleteWorks([]int{6, -1})
		wantStatus := status.New(status.OK, "The following works could not be found: [-1]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{-1}, gotIDs)
	})

	t.Run("AllNonExistentIDs", func(t *testing.T) {
		gotStatus, gotIDs := ws.DeleteWorks([]int{1, 2, 3})
		wantStatus := status.New(status.OK, "The following works could not be found: [1 2 3]")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, []int{1, 2, 3}, gotIDs)
	})

	cleanup()
}

func TestFindWork(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	works := data.GetWorks(ws)
	helpers.PostWorks(t, ws, works)

	t.Run("ExistingWork", func(t *testing.T) {
		tw := &works[0]
		gotStatus, gotWork := ws.FindWork(tw.Title, tw.Author.ID)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tw, gotWork)
	})

	t.Run("NonExistentWork", func(t *testing.T) {
		tw := work.Work{Title: "A", Author: author.Author{ID: 1000}}
		gotStatus, gotWork := ws.FindWork(tw.Title, tw.Author.ID)
		wantStatus := status.New(status.NotFound, "Work ('A', 1000) does not exist")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	cleanup()
}

func TestGetWork(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	works := data.GetWorks(ws)
	helpers.PostWorks(t, ws, works)

	t.Run("ValidID", func(t *testing.T) {
		gotStatus, gotWork := ws.GetWork(1)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, &works[0], gotWork)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		gotStatus, gotWork := ws.GetWork(-1)
		wantStatus := status.New(status.NotFound, "Work with id = -1 does not exist")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	cleanup()
}

func TestGetWorks(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	tws := data.GetWorks(ws)
	works := helpers.PostWorks(t, ws, tws)

	t.Run("AllWorks", func(t *testing.T) {
		gotStatus, gotWorks := ws.GetWorks()
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		for i := range works {
			helpers.AssertEqual(t, works[i], gotWorks[i])
		}
	})

	cleanup()
}

func TestPatchWork(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ws := services.WorkService
	locations := data.GetLocations()
	authors := data.GetAuthors(&ws.AuthorService)
	works := data.GetWorks(ws)

	t.Run("NoFieldsToUpdate", func(t *testing.T) {
		tw := helpers.PostWork(t, ws, &works[0])
		defer helpers.DeleteWork(t, ws, tw.ID)

		updates := &work.Work{ID: tw.ID}

		gotStatus, gotWork := ws.PatchWork(updates)
		wantStatus := status.New(status.OK, "No fields in work to update")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	t.Run("Description", func(t *testing.T) {
		tw := helpers.PostWork(t, ws, &works[0])
		defer helpers.DeleteWork(t, ws, tw.ID)

		updates := &work.Work{ID: tw.ID, Description: "This is a test description..."}

		gotStatus, gotWork := ws.PatchWork(updates)
		wantStatus := status.New(status.OK, "")
		tw.Description = "This is a test description..."
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tw, gotWork)
	})

	t.Run("InitialPubDate", func(t *testing.T) {
		tw := helpers.PostWork(t, ws, &works[0])
		defer helpers.DeleteWork(t, ws, tw.ID)

		updates := &work.Work{ID: tw.ID, InitialPubDate: "1900-01-01T00:00:00Z"}

		gotStatus, gotWork := ws.PatchWork(updates)
		wantStatus := status.New(status.OK, "")
		tw.InitialPubDate = "1900-01-01T00:00:00Z"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tw, gotWork)
	})

	t.Run("OriginalLanguage", func(t *testing.T) {
		tw := helpers.PostWork(t, ws, &works[0])
		defer helpers.DeleteWork(t, ws, tw.ID)

		updates := &work.Work{ID: tw.ID, OriginalLanguage: "French"}

		gotStatus, gotWork := ws.PatchWork(updates)
		wantStatus := status.New(status.OK, "")
		tw.OriginalLanguage = "French"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tw, gotWork)
	})

	t.Run("Title", func(t *testing.T) {
		tw := helpers.PostWork(t, ws, &works[0])
		defer helpers.DeleteWork(t, ws, tw.ID)

		updates := &work.Work{ID: tw.ID, Title: "Testing"}

		gotStatus, gotWork := ws.PatchWork(updates)
		wantStatus := status.New(status.OK, "")
		tw.Title = "Testing"
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tw, gotWork)
	})

	t.Run("Author", func(t *testing.T) {
		t.Run("ExistingID", func(t *testing.T) {
			tw := helpers.PostWork(t, ws, &works[0])
			defer helpers.DeleteWork(t, ws, tw.ID)

			ta := author.Author{ID: len(authors) - 1}
			updates := &work.Work{ID: tw.ID, Author: ta}

			gotStatus, gotWork := ws.PatchWork(updates)
			wantStatus := status.New(status.OK, "")
			tw.Author.ID = ta.ID
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tw, gotWork)
		})

		t.Run("NewFirstNameLastNameDateOfBirth", func(t *testing.T) {
			tw := helpers.PostWork(t, ws, &works[0])
			defer helpers.DeleteWork(t, ws, tw.ID)

			ta := author.Author{
				FirstName:    "A",
				LastName:     "B",
				Gender:       "C",
				DateOfBirth:  "1900-01-01T00:00:00Z",
				PlaceOfBirth: locations[0],
			}
			updates := &work.Work{ID: tw.ID, Author: ta}

			gotStatus, gotWork := ws.PatchWork(updates)
			wantStatus := status.New(status.OK, "")
			tw.Author = author.Author{ID: len(authors) + 1}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tw, gotWork)
		})

		t.Run("ExistingFirstNameLastNameDateOfBirth", func(t *testing.T) {
			tw := helpers.PostWork(t, ws, &works[0])
			defer helpers.DeleteWork(t, ws, tw.ID)

			ta := authors[0]
			updates := &work.Work{ID: tw.ID, Author: ta}

			gotStatus, gotWork := ws.PatchWork(updates)
			wantStatus := status.New(status.OK, "")
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tw, gotWork)
		})
	})

	t.Run("DuplicateAuthorTitle", func(t *testing.T) {
		tw := helpers.PostWork(t, ws, &works[0])
		defer helpers.DeleteWork(t, ws, tw.ID)
		dup := helpers.PostWork(t, ws, &works[1])
		defer helpers.DeleteWork(t, ws, dup.ID)

		updates := &work.Work{ID: tw.ID, Title: dup.Title, Author: dup.Author}
		gotStatus, gotWork := ws.PatchWork(updates)
		msg := "pq: duplicate key value violates unique constraint \"work_author_id_title_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	cleanup()
}

func TestPostWork(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	as := services.AuthorService
	ws := services.WorkService
	authors := data.GetAuthors(&ws.AuthorService)
	works := data.GetWorks(ws)

	t.Run("AllFieldsValid", func(t *testing.T) {
		tw := &works[0]
		gotStatus, gotWork := ws.PostWork(tw)
		defer helpers.DeleteWork(t, ws, gotWork.ID)

		wantStatus := status.New(status.Created, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tw, gotWork)
	})

	t.Run("Author", func(t *testing.T) {
		t.Run("ExistingID", func(t *testing.T) {
			tw := &works[0]
			tw.Author = author.Author{ID: 1}
			gotStatus, gotWork := ws.PostWork(tw)
			defer helpers.DeleteWork(t, ws, gotWork.ID)

			wantStatus := status.New(status.Created, "")
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tw, gotWork)
		})

		t.Run("ExistingFirstNameLastNameDateOfBirth", func(t *testing.T) {
			tw := &works[0]
			ta := &authors[0]
			tw.Author = *ta
			gotStatus, gotWork := ws.PostWork(tw)
			defer helpers.DeleteWork(t, ws, gotWork.ID)

			wantStatus := status.New(status.Created, "")
			au := helpers.FindAuthor(t, as, ta, ta.FirstName, ta.LastName, ta.DateOfBirth)
			tw.Author = author.Author{ID: au.ID}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tw, gotWork)
		})

		t.Run("NewFirstNameLastNameDateOfBirth", func(t *testing.T) {
			tw := &works[0]
			ta := &author.Author{
				FirstName:    "A",
				LastName:     "B",
				DateOfBirth:  "1900-01-01T00:00:00Z",
				PlaceOfBirth: location.Location{ID: 1},
			}
			tw.Author = *ta
			gotStatus, gotWork := ws.PostWork(tw)
			defer helpers.DeleteWork(t, ws, gotWork.ID)

			wantStatus := status.New(status.Created, "")
			au := helpers.FindAuthor(t, as, ta, "A", "B", "1900-01-01T00:00:00Z")
			tw.Author = author.Author{ID: au.ID}
			helpers.AssertEqual(t, wantStatus, gotStatus)
			helpers.AssertEqual(t, tw, gotWork)
		})
	})

	t.Run("DuplicateTitleAuthorID", func(t *testing.T) {
		tw := &works[0]
		gw := helpers.PostWork(t, ws, tw)
		defer helpers.DeleteWork(t, ws, gw.ID)

		gotStatus, gotWork := ws.PostWork(tw)
		msg := "pq: duplicate key value violates unique constraint \"work_author_id_title_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotWork)
	})

	cleanup()
}
