package infrastructure_test

import (
	"reflect"
	"stock-monitor/infrastructure"
	"testing"
)

func TestReadCsv(t *testing.T) {
	csv := "test.csv"

	csvData, _ := infrastructure.ReadData(csv)

	expected := [][]string{
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

	if reflect.DeepEqual(csvData, expected) == false {
		t.Errorf("Unexpected csv parse result. got: %#v, want: %#v", csvData, expected)
	}
}
