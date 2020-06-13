package test

import (
	"reflect"
	"testing"
)

// AssertEqual compares `want` and `got` and raises an error if they are not deeply equal.
func AssertEqual(t *testing.T, want interface{}, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nWant: %v\nGot:  %v\n", want, got)
	}
}

// AssertNil raises an error if `got` is non-nil.
func AssertNil(t *testing.T, got interface{}) {
	t.Helper()

	if !reflect.ValueOf(got).IsNil() {
		t.Errorf("\nWant: %v\nGot:  %v\n", nil, got)
	}
}
