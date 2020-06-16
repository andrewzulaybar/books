package helpers

import (
	"reflect"
	"testing"

	"github.com/andrewzulaybar/books/api/pkg/location"
)

// AssertEqual compares `want` and `got` and raises an error if they are not deeply equal.
func AssertEqual(t *testing.T, want interface{}, got interface{}) {
	t.Helper()

	switch want.(type) {
	case location.Location, *location.Location:
		assertEqualLocation(t, want, got)
	default:
		assertEqual(t, want, got)
	}
}

// AssertNil raises an error if `got` is non-nil.
func AssertNil(t *testing.T, got interface{}) {
	t.Helper()

	if !reflect.ValueOf(got).IsNil() {
		t.Errorf("\nWant: %v\nGot:  %v\n", nil, got)
	}
}

func assertEqual(t *testing.T, want interface{}, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nWant: %v\nGot:  %v\n", want, got)
	}
}

func assertEqualLocation(t *testing.T, want, got interface{}) {
	t.Helper()

	switch wantLocation := want.(type) {
	case location.Location:
		gotLocation, ok := got.(location.Location)
		if !ok {
			t.Errorf("\nWant: %v\nGot:  %v\n", wantLocation, got)
		}
		wantLocation.ID = gotLocation.ID
		assertEqual(t, wantLocation, gotLocation)
	case *location.Location:
		gotLocation, ok := got.(*location.Location)
		if !ok {
			t.Errorf("\nWant: %v\nGot:  %v\n", wantLocation, got)
		}
		wl := *wantLocation
		gl := *gotLocation
		wl.ID = gl.ID
		assertEqual(t, &wl, &gl)
	}
}
