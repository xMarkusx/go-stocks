package command_test

import (
	"reflect"
	"stock-monitor/application/portfolio/command"
	"testing"
	"time"
)

func TestAddSharesToPortfolioCommandHasTodayAsDefaultDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	addSharesToPortfolioCommand := command.NewAddSharesToPortfolioCommand("MO", 20, 19.99)
	expected := command.AddSharesToPortfolioCommand{"MO", 20, 19.99, today}

	if reflect.DeepEqual(addSharesToPortfolioCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", addSharesToPortfolioCommand, expected)
	}
}

func TestRemoveSharesFromPortfolioCommandHasTodayAsDefaultDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	removeSharesFromPortfolioCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 20, 19.99)
	expected := command.RemoveSharesFromPortfolioCommand{"MO", 20, 19.99, today}

	if reflect.DeepEqual(removeSharesFromPortfolioCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", removeSharesFromPortfolioCommand, expected)
	}
}

func TestRenameTickerCommandHasTodayAsDefaultDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	renameCommand := command.NewRenameTickerCommand("MO", "FOO")
	expected := command.RenameTickerCommand{"MO", "FOO", today}

	if reflect.DeepEqual(renameCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", renameCommand, expected)
	}
}
