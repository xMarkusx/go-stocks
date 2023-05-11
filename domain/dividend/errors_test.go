package dividend_test

import (
	"stock-monitor/domain/dividend"
	"testing"
)

func TestTickerUnknownError(t *testing.T) {
	err := dividend.NewTickerUnknownError("FOO")

	expected := "ticker not added to portfolio. ticker: FOO"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestDividendDateBeforeSharesWereAddedToPortfolioError(t *testing.T) {
	err := dividend.NewDividendDateBeforeSharesWereAddedToPortfolioError("FOO", "2000-01-01")

	expected := "dividend date is before shares were added to portfolio. ticker: FOO date: 2000-01-01"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestDividendNetZeroOrNegativeError(t *testing.T) {
	err := dividend.DividendNetZeroOrNegativeError{}

	expected := "dividend net must be greater than zero"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}

func TestDividendGrossZeroOrNegativeError(t *testing.T) {
	err := dividend.DividendGrossZeroOrNegativeError{}

	expected := "dividend gross must be greater than zero"
	got := err.Error()

	if expected != got {
		t.Errorf("Unexpected error text. Expected:%#v Got:%#v", expected, got)
	}
}
