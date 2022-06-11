package importer_test

import (
	"reflect"
	"stock-monitor/application/dividend/importer"
	"testing"
)

func TestParse(t *testing.T) {
	csv := [][]string{
		{
			"2000-01-01",
			"MO",
			"12.3456",
			"23.4567",
		},
		{
			"2000-01-01",
			"PG",
			"0.12345",
			"1.2345",
		},
		{
			"2000-01-01",
			"GIS",
			"123",
			"456",
		},
	}

	importItems := importer.Parse(csv)

	expected := []importer.ImportItem{
		{
			"2000-01-01",
			"MO",
			12.3456,
			23.4567,
		},
		{
			"2000-01-01",
			"PG",
			0.12345,
			1.2345,
		},
		{
			"2000-01-01",
			"GIS",
			123,
			456,
		},
	}

	if reflect.DeepEqual(importItems, expected) == false {
		t.Errorf("Unexpected csv parse result. got: %#v, want: %#v", importItems, expected)
	}
}
