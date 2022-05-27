package portfolio_test

import (
	"stock-monitor/domain/portfolio"
	"testing"
)

func TestInvalidNumbersOfSharesError(t *testing.T) {
	err := &portfolio.InvalidNumbersOfSharesError{}

	expected := "number of shares must be greater than 0"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestCantSellMoreSharesThanExistingError(t *testing.T) {
	err := &portfolio.CantSellMoreSharesThanExistingError{}

	expected := "not allowed to sell more shares than currently in portfolio"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestTickerNotInPortfolioError(t *testing.T) {
	err := portfolio.NewTickerNotInPortfolioError("FOO")

	expected := "Ticker to be renamed not found. Ticker: FOO"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestTickerAlreadyUsedError(t *testing.T) {
	err := portfolio.NewTickerAlreadyUsedError("FOO")

	expected := "New ticker symbol already in use. Ticker: FOO"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}
