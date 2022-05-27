package portfolio

import (
	"testing"
)

func TestInvalidNumbersOfSharesError(t *testing.T) {
	err := &InvalidNumbersOfSharesError{"Error text"}

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestCantSellMoreSharesThanExistingError(t *testing.T) {
	err := &CantSellMoreSharesThanExistingError{"Error text"}

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestTickerNotInPortfolioError(t *testing.T) {
	err := &TickerNotInPortfolioError{"Error text"}

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestTickerAlreadyUsedError(t *testing.T) {
	err := &TickerAlreadyUsedError{"Error text"}

	expected := "Error text"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}
