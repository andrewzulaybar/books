package helpers

import (
	"testing"

	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

// DeleteWork is a test helper that calls (WorkService).DeleteWork.
func DeleteWork(t *testing.T, ws *work.Service, id int) {
	t.Helper()

	gotStatus := ws.DeleteWork(id)
	wantStatus := status.New(status.NoContent, "")
	AssertEqual(t, wantStatus, gotStatus)
}

// FindWork is a test helper that calls (WorkService).FindWork with the given title and authorID.
func FindWork(t *testing.T, ws *work.Service, wantWork *work.Work, title string, authorID int) *work.Work {
	t.Helper()

	gotStatus, gotWork := ws.FindWork(title, authorID)
	wantStatus := status.New(status.OK, "")
	AssertEqual(t, gotStatus, wantStatus)
	AssertEqual(t, gotWork, wantWork)

	return gotWork
}

// PostWork is a test helper that calls (WorkService).PostWork.
func PostWork(t *testing.T, ws *work.Service, tw *work.Work) *work.Work {
	t.Helper()

	gotStatus, gotWork := ws.PostWork(tw)
	wantStatus := status.New(status.Created, "")
	AssertEqual(t, wantStatus, gotStatus)
	AssertEqual(t, tw, gotWork)

	return gotWork
}

// PostWorks is a test helper that calls (WorkServive).PostWorks for each work in `works`.
func PostWorks(t *testing.T, ws *work.Service, tw work.Works) work.Works {
	t.Helper()

	var works work.Works
	for i := range tw {
		work := PostWork(t, ws, &tw[i])
		if work != nil {
			works = append(works, *work)
		}
	}
	return works
}
