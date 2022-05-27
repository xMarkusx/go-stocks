package infrastructure

import (
	"testing"
)

func TestUnsupportedDateFormatError(t *testing.T) {
	err := &UnsupportedDateFormatError{"Error text"}

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestInvalidDateError(t *testing.T) {
	err := &InvalidDateError{"Error text"}

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}
