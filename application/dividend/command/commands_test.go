package command_test

import (
	"reflect"
	"stock-monitor/application/dividend/command"
	"testing"
)

func TestRecordDividendCommand(t *testing.T) {
	recordDividendCommand := command.NewRecordDividendCommand("MO", 20.00, 19.99, "2001-01-01")
	expected := command.RecordDividendCommand{"MO", 20.00, 19.99, "2001-01-01"}

	if reflect.DeepEqual(recordDividendCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", recordDividendCommand, expected)
	}
}
