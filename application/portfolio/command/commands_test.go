package command

import (
	"reflect"
	"testing"
	"time"
)

func TestAddSharesToPortfolioCommandHasTodayAsDefaultDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	command := NewAddSharesToPortfolioCommand("MO", 20, 19.99)
	expected := AddSharesToPortfolioCommand{"MO", 20, 19.99, today}

	if reflect.DeepEqual(command, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", command, expected)
	}
}

func TestRemoveSharesFromPortfolioCommandHasTodayAsDefaultDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	command := NewRemoveSharesFromPortfolioCommand("MO", 20, 19.99)
	expected := RemoveSharesFromPortfolioCommand{"MO", 20, 19.99, today}

	if reflect.DeepEqual(command, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", command, expected)
	}
}
