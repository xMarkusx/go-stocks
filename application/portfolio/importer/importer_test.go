package importer_test

import (
	"reflect"
	"stock-monitor/application/portfolio/importer"
	"testing"
)

func TestParse(t *testing.T) {
	csv := [][]string{
		{
			"buy",
			"2000-01-01",
			"MO",
			"",
			"12.3456",
			"100",
		},
		{
			"sell",
			"2000-01-01",
			"MO",
			"",
			"12.3456",
			"100",
		},
		{
			"rename",
			"2000-01-01",
			"MO",
			"FOO",
			"",
			"",
		},
	}

	importItems := importer.Parse(csv)

	expected := []importer.ImportItem{
		{
			"buy",
			"2000-01-01",
			"MO",
			"",
			12.3456,
			100,
		},
		{
			"sell",
			"2000-01-01",
			"MO",
			"",
			12.3456,
			100,
		},
		{
			"rename",
			"2000-01-01",
			"MO",
			"FOO",
			0,
			0,
		},
	}

	if reflect.DeepEqual(importItems, expected) == false {
		t.Errorf("Unexpected csv parse result. got: %#v, want: %#v", importItems, expected)
	}
}
