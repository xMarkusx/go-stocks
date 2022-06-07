package dividend_test

import (
	"reflect"
	"stock-monitor/domain"
	"stock-monitor/domain/dividend"
	"stock-monitor/domain/portfolio"
	"testing"
)

func TestCanRecordADividend(t *testing.T) {
	d := dividend.NewDividend()
	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	d.Apply(&sharesAddedEvent)

	err := d.RecordDividend("MO", 20.00, 30.00, "2000-01-02")

	events := d.GetRecordedEvents()

	expectedEvent := dividend.NewDividendRecordedEvent("MO", 20.00, 30.00, "2000-01-02")
	expectedEvents := []domain.DomainEvent{
		&expectedEvent,
	}

	if err != nil {
		t.Errorf("Unexpected error")
	}
	if reflect.DeepEqual(events, expectedEvents) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEvents, events)
	}
}

func TestCanNotRecordADividendWhenTickerWasNotAddedToPortfolio(t *testing.T) {
	d := dividend.NewDividend()

	err := d.RecordDividend("MO", 20.00, 30.00, "2000-01-01")

	_, ok := err.(*dividend.TickerUnknownError)
	if !ok {
		t.Errorf("Expected TickerUnknownError but got %#v", err)
	}
}

func TestCanNotRecordADividendWhenTickerWasAddedToPortfolioLaterThanDividendDate(t *testing.T) {
	d := dividend.NewDividend()
	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	d.Apply(&sharesAddedEvent)

	err := d.RecordDividend("MO", 20.00, 30.00, "2000-01-01")

	_, ok := err.(*dividend.DividendDateBeforeSharesWereAddedToPortfolioError)
	if !ok {
		t.Errorf("Expected DividendDateBeforeSharesWereAddedToPortfolioError but got %#v", err)
	}
}

func TestDividendNetHasToBeGreaterThanZero(t *testing.T) {
	d := dividend.NewDividend()
	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	d.Apply(&sharesAddedEvent)

	err := d.RecordDividend("MO", 0, 30.00, "2000-01-02")

	_, ok := err.(*dividend.DividendNetZeroOrNegativeError)
	if !ok {
		t.Errorf("Expected DividendNetZeroOrNegativeError but got %#v", err)
	}
}

func TestDividendGrossHasToBeGreaterThanZero(t *testing.T) {
	d := dividend.NewDividend()
	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	d.Apply(&sharesAddedEvent)

	err := d.RecordDividend("MO", 20.00, 0, "2000-01-02")

	_, ok := err.(*dividend.DividendGrossZeroOrNegativeError)
	if !ok {
		t.Errorf("Expected DividendGrossZeroOrNegativeError but got %#v", err)
	}
}

func TestDateOfLaterAddedSharesIsIgnoredForDividendDateValidation(t *testing.T) {
	d := dividend.NewDividend()
	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	sharesAddedEvent2 := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2001-01-01")
	d.Apply(&sharesAddedEvent)
	d.Apply(&sharesAddedEvent2)

	err := d.RecordDividend("MO", 20.00, 30.00, "2000-01-02")

	if err != nil {
		t.Errorf("Unexpected error. Got %#v", err.Error())
	}
}

func TestTickerRenamesAreHandledWhenCheckingDividendDate(t *testing.T) {
	d := dividend.NewDividend()
	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	sharesAddedEvent2 := portfolio.NewTickerRenamedEvent("MO", "FOO")
	d.Apply(&sharesAddedEvent)
	d.Apply(&sharesAddedEvent2)

	err := d.RecordDividend("FOO", 20.00, 30.00, "2000-01-02")

	if err != nil {
		t.Errorf("Unexpected error. Got %#v", err.Error())
	}
}
