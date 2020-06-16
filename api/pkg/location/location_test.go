package location_test

import (
	"testing"

	"github.com/andrewzulaybar/books/api/pkg/internal/test/helpers"
	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
	"github.com/andrewzulaybar/books/api/test/data"
)

func TestDeleteLocation(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ls := services.LocationService
	locations := data.GetLocations()
	helpers.PostLocations(t, ls, locations)

	t.Run("ExistingID", func(t *testing.T) {
		got := ls.DeleteLocation(1)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("AlreadyDeletedID", func(t *testing.T) {
		got := ls.DeleteLocation(2)
		want := status.New(status.NoContent, "")
		helpers.AssertEqual(t, want, got)

		got = ls.DeleteLocation(2)
		want = status.New(status.OK, "Location with id = 2 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		got := ls.DeleteLocation(-1)
		want := status.New(status.OK, "Location with id = -1 does not exist")
		helpers.AssertEqual(t, want, got)
	})

	cleanup()
}

func TestFindLocation(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ls := services.LocationService
	locations := data.GetLocations()
	helpers.PostLocations(t, ls, locations)

	t.Run("ExistingLocation", func(t *testing.T) {
		tl := &locations[0]
		gotStatus, gotLocation := ls.FindLocation(tl.City, tl.Country)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tl, gotLocation)
	})

	t.Run("NonExistentLocation", func(t *testing.T) {
		tl := &location.Location{City: "A", Country: "B", Region: "C"}
		gotStatus, gotLocation := ls.FindLocation(tl.City, tl.Country)
		wantStatus := status.New(status.NotFound, "Location ('A', 'B') does not exist")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotLocation)
	})

	cleanup()
}

func TestGetLocation(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ls := services.LocationService
	locations := data.GetLocations()
	helpers.PostLocations(t, ls, locations)

	t.Run("ExistingID", func(t *testing.T) {
		gotStatus, gotLocation := ls.GetLocation(1)
		wantStatus := status.New(status.OK, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, &locations[0], gotLocation)
	})

	t.Run("NonExistentID", func(t *testing.T) {
		gotStatus, gotLocation := ls.GetLocation(-1)
		wantStatus := status.New(status.NotFound, "Location with id = -1 does not exist")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertNil(t, gotLocation)
	})

	cleanup()
}

func TestPostLocation(t *testing.T) {
	services, cleanup := helpers.Setup(t)
	ls := services.LocationService
	locations := data.GetLocations()

	t.Run("AllFieldsValid", func(t *testing.T) {
		tl := &locations[0]
		gotStatus, gotLocation := ls.PostLocation(tl)
		defer helpers.DeleteLocation(t, ls, gotLocation.ID)

		wantStatus := status.New(status.Created, "")
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tl, gotLocation)
	})

	t.Run("ExistingID", func(t *testing.T) {
		tl := &locations[0]
		gl := helpers.PostLocation(t, ls, tl)
		defer helpers.DeleteLocation(t, ls, gl.ID)

		gotStatus, gotLocation := ls.PostLocation(gl)
		wantStatus := status.Newf(status.Conflict, "Location with id = %d already exists", gl.ID)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tl, gotLocation)
	})

	t.Run("DuplicateCityCountry", func(t *testing.T) {
		tl := &locations[0]
		gl := helpers.PostLocation(t, ls, tl)
		defer helpers.DeleteLocation(t, ls, gl.ID)

		gotStatus, gotLocation := ls.PostLocation(tl)
		msg := "pq: duplicate key value violates unique constraint \"location_city_country_key\""
		wantStatus := status.New(status.Conflict, msg)
		helpers.AssertEqual(t, wantStatus, gotStatus)
		helpers.AssertEqual(t, tl, gotLocation)
	})

	cleanup()
}
