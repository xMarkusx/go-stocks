package command_test

import (
	"reflect"
	"stock-monitor/application/dividend/command"
	"testing"
	"time"
)

func TestRecordDividendCommandHasTodayAsDefaultDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	recordDividendCommand := command.NewRecordDividendCommand("MO", 20.00, 19.99)
	expected := command.RecordDividendCommand{"MO", 20.00, 19.99, today}

	if reflect.DeepEqual(recordDividendCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", recordDividendCommand, expected)
	}
}
