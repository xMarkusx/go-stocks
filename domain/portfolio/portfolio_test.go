package portfolio

import (
	"reflect"
	"stock-monitor/domain"
	"testing"
	"time"
)

func TestCanAddShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", 10, 9.99, "2000-01-01")
	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}

	expectedEvent := NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")
	expectedEventArray := []domain.DomainEvent{
		&expectedEvent,
	}
	got := portfolio.GetRecordedEvents()

	if reflect.DeepEqual(got, expectedEventArray) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEventArray, got)
	}
}

func TestCanRemoveShares(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)
	err := portfolio.RemoveSharesFromPortfolio("MO", 10, 9.99, "2000-01-02")
	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}

	expectedEvent := NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99, "2000-01-02")
	expectedEventArray := []domain.DomainEvent{
		&expectedEvent,
	}
	got := portfolio.GetRecordedEvents()

	if reflect.DeepEqual(got, expectedEventArray) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEventArray, got)
	}
}

func TestSharesAddedToPortfolioEventCanBeApplied(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99, "2000-01-01")
	sharesAddedEvent2 := NewSharesAddedToPortfolioEvent("MO", 9, 9.99, "2000-01-02")
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&sharesAddedEvent2)

	err := portfolio.RemoveSharesFromPortfolio("MO", 20, 9.99, "2000-01-03")

	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}
}

func TestSharesRemovedFromPortfolioEventCanBeApplied(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99, "2000-01-01")
	sharesRemovedEvent := NewSharesRemovedFromPortfolioEvent("MO", 11, 9.99, "2000-01-02")
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&sharesRemovedEvent)

	err := portfolio.RemoveSharesFromPortfolio("MO", 1, 9.99, "2000-01-03")

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", 0, 9.99, "2000-01-02")

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", -10, 9.99, "2000-01-02")

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.RemoveSharesFromPortfolio("MO", 21, 9.99, "2000-01-02")

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenAddingShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", 10, 9.99, "Foo")

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenRemovingShares(t *testing.T) {
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RemoveSharesFromPortfolio("MO", 10, 9.99, "Foo")

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenAddingShares(t *testing.T) {
	today := time.Now()
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", 10, 9.99, today.AddDate(0, 0, 1).Format("2006-01-02"))

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenRemovingShares(t *testing.T) {
	today := time.Now()
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RemoveSharesFromPortfolio("MO", 10, 9.99, today.AddDate(0, 0, 1).Format("2006-01-02"))

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenAddingShares(t *testing.T) {
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2020-01-02")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.AddSharesToPortfolio("MO", 10, 9.99, "2020-01-01")

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenRemovingShares(t *testing.T) {
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2020-01-02")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RemoveSharesFromPortfolio("MO", 10, 9.99, "2020-01-01")

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestTickerCanBeRenamed(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)

	portfolio.RenameTicker("MO", "FOO", "2000-01-02")

	expectedEvent := NewTickerRenamedEvent("MO", "FOO", "2000-01-02")
	expectedEventArray := []domain.DomainEvent{
		&expectedEvent,
	}
	got := portfolio.GetRecordedEvents()

	if reflect.DeepEqual(got, expectedEventArray) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEventArray, got)
	}
}

func TestTickerHasToBePresentInPortfolioToBeRenamed(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RenameTicker("PG", "FOO", "2000-01-02")

	_, ok := err.(*TickerNotInPortfolioError)
	if !ok {
		t.Errorf("Expected TickerNotInPortfolioError but got %#v", err)
	}
}

func TestTickerCanBeRenamedEvenIfThereAreNoSharesHeld(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 1, 9.99, "2000-01-01")
	removeSharesEvent := NewSharesRemovedFromPortfolioEvent("MO", 1, 9.99, "2000-01-02")
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&removeSharesEvent)

	err := portfolio.RenameTicker("MO", "FOO", "2000-01-02")

	if err != nil {
		t.Errorf("Got unexpected error: %#v", err)
	}
}

func TestNewTickerMustNotBeAlreadyInPortfolio(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99, "2000-01-01")
	sharesAddedEvent2 := NewSharesAddedToPortfolioEvent("PG", 11, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&sharesAddedEvent2)

	err := portfolio.RenameTicker("MO", "PG", "2000-01-02")

	_, ok := err.(*TickerAlreadyUsedError)
	if !ok {
		t.Errorf("Expected TickerAlreadyUsedError but got %#v", err)
	}
}

func TestNewTickerWillBeUsedForAnyNewPortfolioCommands(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 1, 9.99, "2000-01-01")
	renameEvent := NewTickerRenamedEvent("MO", "FOO", "2000-01-02")
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&renameEvent)

	err := portfolio.RemoveSharesFromPortfolio("FOO", 1, 9.99, "2000-01-03")

	if err != nil {
		t.Errorf("Got unexpected error: %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenRenamingTicker(t *testing.T) {
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RenameTicker("MO", "FOO", "Foo")

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenRenamingTicker(t *testing.T) {
	today := time.Now()
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2000-01-01")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RenameTicker("MO", "FOO", today.AddDate(0, 0, 1).Format("2006-01-02"))

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenRenamingTicker(t *testing.T) {
	portfolio := NewPortfolio()
	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 20, 9.99, "2020-01-02")
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RenameTicker("MO", "FOO", "2020-01-01")

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}
