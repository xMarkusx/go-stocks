package shared_test

import (
	"stock-monitor/application/shared"
	"testing"
	"time"
)

func TestDefaultCommandDateIsToday(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	date := shared.CommandDate("")
	expected := today

	if date.Get() != expected {
		t.Errorf("Unexpected date. got: %#v, want: %#v", date.Get(), expected)
	}
}

func TestProvidesTheDateAsString(t *testing.T) {
	date := shared.CommandDate("2001-01-02")
	expected := "2001-01-02"

	if date.Get() != expected {
		t.Errorf("Unexpected date. got: %#v, want: %#v", date.Get(), expected)
	}
}
