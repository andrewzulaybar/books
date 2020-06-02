package work_test

import (
	"reflect"
	"testing"

	"github.com/andrewzulaybar/books/api/config"
	"github.com/andrewzulaybar/books/api/internal/postgres"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/pkg/work"
	"github.com/andrewzulaybar/books/api/test/data"
)

func getDB(t *testing.T) (*postgres.DB, func()) {
	t.Helper()

	conf, err := config.Load("../../config/test.env")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return postgres.Setup(conf.ConnectionString, "../../internal/sql/")
}

func TestDeleteWork(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	w := &work.Service{DB: *db}
	data.LoadWorks(w)

	type Expected struct {
		status *status.Status
	}

	cases := []struct {
		name     string
		id       int
		expected Expected
	}{
		{
			name: "ValidId",
			id:   1,
			expected: Expected{
				status: status.New(status.NoContent, ""),
			},
		},
		{
			name: "AlreadyDeletedId",
			id:   1,
			expected: Expected{
				status: status.New(status.OK, "Work with id = 1 does not exist"),
			},
		},
		{
			name: "InvalidId",
			id:   -1,
			expected: Expected{
				status: status.New(status.OK, "Work with id = -1 does not exist"),
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s := w.DeleteWork(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
		})
	}
}

func TestDeleteWorks(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	w := &work.Service{DB: *db}
	data.LoadWorks(w)

	type Expected struct {
		status *status.Status
		ids    []int
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
				status: status.New(status.NoContent, ""),
				ids:    nil,
			},
		},
		{
			name: "MultipleIds",
			ids:  []int{2, 3, 4},
			expected: Expected{
				status: status.New(status.NoContent, ""),
				ids:    nil,
			},
		},
		{
			name: "AlreadyDeletedIds",
			ids:  []int{1, 2, 3},
			expected: Expected{
				status: status.New(status.OK, "The following works could not be found: [1 2 3]"),
				ids:    []int{1, 2, 3},
			},
		},
		{
			name: "IncludesIdNotFound",
			ids:  []int{5, 6, -1},
			expected: Expected{
				status: status.New(status.OK, "The following works could not be found: [-1]"),
				ids:    []int{-1},
			},
		},
		{
			name: "AllIdsNotFound",
			ids:  []int{-1, -2, -3},
			expected: Expected{
				status: status.New(status.OK, "The following works could not be found: [-1 -2 -3]"),
				ids:    []int{-1, -2, -3},
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, ids := w.DeleteWorks(c.ids)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.ids, ids) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.ids, ids)
			}
		})
	}
}

func TestGetWork(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	w := &work.Service{DB: *db}
	works := data.LoadWorks(w)

	type Expected struct {
		status *status.Status
		work   *work.Work
	}

	cases := []struct {
		name     string
		id       int
		expected Expected
	}{
		{
			name: "ValidId",
			id:   1,
			expected: Expected{
				status: status.New(status.OK, ""),
				work:   &works[0],
			},
		},
		{
			name: "InvalidId",
			id:   -1,
			expected: Expected{
				status: status.New(status.NotFound, "Work with id = -1 does not exist"),
				work:   nil,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, work := w.GetWork(c.id)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			if !reflect.DeepEqual(exp.work, work) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.work, work)
			}
		})
	}
}

func TestGetWorks(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	w := &work.Service{DB: *db}
	works := data.LoadWorks(w)

	type Expected struct {
		status *status.Status
		works  work.Works
	}

	exp := Expected{
		status: status.New(status.OK, ""),
		works:  works,
	}
	t.Run("AllWorks", func(t *testing.T) {
		s, works := w.GetWorks()
		if !reflect.DeepEqual(exp.status, s) {
			t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
		}
		if !reflect.DeepEqual(exp.works, works) {
			t.Errorf("\nExpected: %v\nActual: %v\n", exp.works, works)
		}
	})
}

func TestPatchWork(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	w := &work.Service{DB: *db}
	works := data.LoadWorks(w)

	type Expected struct {
		status *status.Status
		id     int
	}

	cases := []struct {
		name     string
		work     *work.Work
		expected Expected
	}{
		{
			name: "EmptyWork",
			work: &work.Work{},
			expected: Expected{
				status: status.New(status.OK, "No fields in work to update"),
				id:     0,
			},
		},
		{
			name: "UpdateAllFields",
			work: &work.Work{
				ID:               1,
				AuthorID:         2,
				Description:      "This is a test description.",
				InitialPubDate:   "1900-01-01T00:00:00Z",
				OriginalLanguage: "French",
				Title:            "Testing",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "UpdateAuthor",
			work: &work.Work{
				ID:       1,
				AuthorID: 2,
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "UpdateDescription",
			work: &work.Work{
				ID:          1,
				Description: "This is a test description.",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "UpdateInitialPubDate",
			work: &work.Work{
				ID:             1,
				InitialPubDate: "1900-01-01T00:00:00Z",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "UpdateOriginalLanguage",
			work: &work.Work{
				ID:               1,
				OriginalLanguage: "French",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "UpdateTitle",
			work: &work.Work{
				ID:    1,
				Title: "Testing",
			},
			expected: Expected{
				status: status.New(status.OK, ""),
				id:     1,
			},
		},
		{
			name: "DuplicateAuthorTitle",
			work: &work.Work{
				ID:       1,
				AuthorID: works[1].AuthorID,
				Title:    works[1].Title,
			},
			expected: Expected{
				status: status.New(
					status.Conflict,
					"pq: duplicate key value violates unique constraint \"work_author_id_title_key\"",
				),
				id: 0,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, work := w.PatchWork(c.work)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			_, want := w.GetWork(exp.id)
			if !reflect.DeepEqual(want, work) {
				t.Errorf("\nExpected: %v\nActual: %v\n", want, work)
			}
		})
	}
}

func TestPostWork(t *testing.T) {
	db, dc := getDB(t)
	defer dc()

	w := &work.Service{DB: *db}
	works := data.GetWorks()

	type Expected struct {
		status *status.Status
		id     int
	}

	cases := []struct {
		name     string
		work     *work.Work
		expected Expected
	}{
		{
			name: "AllFieldsValid",
			work: func() *work.Work {
				work := works[0]
				work.ID = 0
				return &work
			}(),
			expected: Expected{
				status: status.New(status.Created, ""),
				id:     1,
			},
		},
		{
			name: "DuplicateAuthorTitle",
			work: func() *work.Work {
				work := works[1]
				work.ID = 0
				work.AuthorID = works[0].AuthorID
				work.Title = works[0].Title
				return &work
			}(),
			expected: Expected{
				status: status.New(
					status.Conflict,
					"pq: duplicate key value violates unique constraint \"work_author_id_title_key\"",
				),
				id: 0,
			},
		},
	}

	for _, c := range cases {
		exp := c.expected
		t.Run(c.name, func(t *testing.T) {
			s, work := w.PostWork(c.work)
			if !reflect.DeepEqual(exp.status, s) {
				t.Errorf("\nExpected: %v\nActual: %v\n", exp.status, s)
			}
			_, want := w.GetWork(exp.id)
			if !reflect.DeepEqual(want, work) {
				t.Errorf("\nExpected: %v\nActual: %v\n", want, work)
			}
		})
	}
}
