package command_test

import (
	"reflect"
	"stock-monitor/application/portfolio/command"
	"testing"
)

func TestNewAddSharesToPortfolioCommand(t *testing.T) {
	addSharesToPortfolioCommand := command.NewAddSharesToPortfolioCommand("MO", 20, 19.99, "2001-01-02")
	expected := command.AddSharesToPortfolioCommand{"MO", 20, 19.99, "2001-01-02"}

	if reflect.DeepEqual(addSharesToPortfolioCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", addSharesToPortfolioCommand, expected)
	}
}

func TestRemoveSharesFromPortfolioCommand(t *testing.T) {
	removeSharesFromPortfolioCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 20, 19.99, "2001-01-02")
	expected := command.RemoveSharesFromPortfolioCommand{"MO", 20, 19.99, "2001-01-02"}

	if reflect.DeepEqual(removeSharesFromPortfolioCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", removeSharesFromPortfolioCommand, expected)
	}
}

func TestRenameTickerCommandHasTodayAsDefaultDate(t *testing.T) {
	renameCommand := command.NewRenameTickerCommand("MO", "FOO", "2001-01-02")
	expected := command.RenameTickerCommand{"MO", "FOO", "2001-01-02"}

	if reflect.DeepEqual(renameCommand, expected) == false {
		t.Errorf("Unexpected command. got: %#v, want: %#v", renameCommand, expected)
	}
}
