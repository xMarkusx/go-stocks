package infrastructure_test

import (
	"stock-monitor/infrastructure"
	"testing"
)

func TestUnsupportedDateFormatError(t *testing.T) {
	err := infrastructure.NewUnsupportedDateFormatError("Error text")

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestInvalidDateError(t *testing.T) {
	err := infrastructure.NewInvalidDateError("Error text")

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}
