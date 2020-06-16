package helpers

import (
	"testing"

	"github.com/andrewzulaybar/books/api/pkg/location"
	"github.com/andrewzulaybar/books/api/pkg/status"
)

// DeleteLocation is a test helper that calls (LocationService).DeleteLocation.
func DeleteLocation(t *testing.T, ls *location.Service, id int) {
	t.Helper()

	gotStatus := ls.DeleteLocation(id)
	wantStatus := status.New(status.NoContent, "")
	AssertEqual(t, wantStatus, gotStatus)
}

// FindLocation is a test helper that calls (LocationService).FindLocation with the given city and country.
func FindLocation(t *testing.T, ls *location.Service, wantLocation *location.Location,
	city, country string) *location.Location {
	t.Helper()

	gotStatus, gotLocation := ls.FindLocation(city, country)
	wantStatus := status.New(status.OK, "")
	AssertEqual(t, gotStatus, wantStatus)
	AssertEqual(t, gotLocation, wantLocation)

	return gotLocation
}

// PostLocation is a test helper that calls (LocationService).PostLocation.
func PostLocation(t *testing.T, ls *location.Service, wl *location.Location) *location.Location {
	t.Helper()

	tl := *wl
	gotStatus, gotLocation := ls.PostLocation(&tl)
	wantStatus := status.New(status.Created, "")
	AssertEqual(t, wantStatus, gotStatus)
	AssertEqual(t, &tl, gotLocation)

	return gotLocation
}

// PostLocations is a test helper that calls (LocationService).PostLocation for each location in `locations`.
func PostLocations(t *testing.T, ls *location.Service, locations location.Locations) {
	t.Helper()

	for _, location := range locations {
		PostLocation(t, ls, &location)
	}
}
