package portfolio

import (
	"reflect"
	"stock-monitor/domain"
	"testing"
)

func TestCanAddShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", 10, 9.99)
	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}

	expectedEvent := NewSharesAddedToPortfolioEvent("MO", 10, 9.99)
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

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	portfolio.Apply(&sharesAddedEvent)
	err := portfolio.RemoveSharesFromPortfolio("MO", 10, 9.99)
	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}

	expectedEvent := NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99)
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

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	sharesAddedEvent2 := NewSharesAddedToPortfolioEvent("MO", 9, 9.99)
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&sharesAddedEvent2)

	err := portfolio.RemoveSharesFromPortfolio("MO", 20, 9.99)

	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}
}

func TestSharesRemovedFromPortfolioEventCanBeApplied(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	sharesRemovedEvent := NewSharesRemovedFromPortfolioEvent("MO", 11, 9.99)
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&sharesRemovedEvent)

	err := portfolio.RemoveSharesFromPortfolio("MO", 1, 9.99)

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", 0, 9.99)

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.AddSharesToPortfolio("MO", -10, 9.99)

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	portfolio := NewPortfolio()

	err := portfolio.RemoveSharesFromPortfolio("MO", 21, 9.99)

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestTickerCanBeRenamed(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	portfolio.Apply(&sharesAddedEvent)

	portfolio.RenameTicker("MO", "FOO")

	expectedEvent := NewTickerRenamedEvent("MO", "FOO")
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

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	portfolio.Apply(&sharesAddedEvent)

	err := portfolio.RenameTicker("PG", "FOO")

	_, ok := err.(*TickerNotInPortfolioError)
	if !ok {
		t.Errorf("Expected TickerNotInPortfolioError but got %#v", err)
	}
}

func TestTickerCanBeRenamedEvenIfThereAreNoSharesHeld(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 1, 9.99)
	removeSharesEvent := NewSharesRemovedFromPortfolioEvent("MO", 1, 9.99)
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&removeSharesEvent)

	err := portfolio.RenameTicker("MO", "FOO")

	if err != nil {
		t.Errorf("Got unexpected error: %#v", err)
	}
}

func TestNewTickerMustNotBeAlreadyInPortfolio(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	sharesAddedEvent2 := NewSharesAddedToPortfolioEvent("PG", 11, 9.99)
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&sharesAddedEvent2)

	err := portfolio.RenameTicker("MO", "PG")

	_, ok := err.(*TickerAlreadyUsedError)
	if !ok {
		t.Errorf("Expected TickerAlreadyUsedError but got %#v", err)
	}
}

func TestNewTickerWillBeUsedForAnyNewPortfolioCommands(t *testing.T) {
	portfolio := NewPortfolio()

	sharesAddedEvent := NewSharesAddedToPortfolioEvent("MO", 1, 9.99)
	renameEvent := NewTickerRenamedEvent("MO", "FOO")
	portfolio.Apply(&sharesAddedEvent)
	portfolio.Apply(&renameEvent)

	err := portfolio.RemoveSharesFromPortfolio("FOO", 1, 9.99)

	if err != nil {
		t.Errorf("Got unexpected error: %#v", err)
	}
}
