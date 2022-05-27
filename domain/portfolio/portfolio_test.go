package portfolio_test

import (
	"reflect"
	"stock-monitor/domain"
	"stock-monitor/domain/portfolio"
	"testing"
)

func TestCanAddShares(t *testing.T) {
	p := portfolio.NewPortfolio()

	err := p.AddSharesToPortfolio("MO", 10, 9.99)
	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}

	expectedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99)
	expectedEventArray := []domain.DomainEvent{
		&expectedEvent,
	}
	got := p.GetRecordedEvents()

	if reflect.DeepEqual(got, expectedEventArray) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEventArray, got)
	}
}

func TestCanRemoveShares(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	p.Apply(&sharesAddedEvent)
	err := p.RemoveSharesFromPortfolio("MO", 10, 9.99)
	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}

	expectedEvent := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99)
	expectedEventArray := []domain.DomainEvent{
		&expectedEvent,
	}
	got := p.GetRecordedEvents()

	if reflect.DeepEqual(got, expectedEventArray) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEventArray, got)
	}
}

func TestSharesAddedToPortfolioEventCanBeApplied(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	sharesAddedEvent2 := portfolio.NewSharesAddedToPortfolioEvent("MO", 9, 9.99)
	p.Apply(&sharesAddedEvent)
	p.Apply(&sharesAddedEvent2)

	err := p.RemoveSharesFromPortfolio("MO", 20, 9.99)

	if err != nil {
		t.Errorf("Unexpected Error. %#v", err)
	}
}

func TestSharesRemovedFromPortfolioEventCanBeApplied(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	sharesRemovedEvent := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 11, 9.99)
	p.Apply(&sharesAddedEvent)
	p.Apply(&sharesRemovedEvent)

	err := p.RemoveSharesFromPortfolio("MO", 1, 9.99)

	_, ok := err.(*portfolio.CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	p := portfolio.NewPortfolio()

	err := p.AddSharesToPortfolio("MO", 0, 9.99)

	_, ok := err.(*portfolio.InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	p := portfolio.NewPortfolio()

	err := p.AddSharesToPortfolio("MO", -10, 9.99)

	_, ok := err.(*portfolio.InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	p := portfolio.NewPortfolio()

	err := p.RemoveSharesFromPortfolio("MO", 21, 9.99)

	_, ok := err.(*portfolio.CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestTickerCanBeRenamed(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	p.Apply(&sharesAddedEvent)

	p.RenameTicker("MO", "FOO")

	expectedEvent := portfolio.NewTickerRenamedEvent("MO", "FOO")
	expectedEventArray := []domain.DomainEvent{
		&expectedEvent,
	}
	got := p.GetRecordedEvents()

	if reflect.DeepEqual(got, expectedEventArray) == false {
		t.Errorf("Expected domain event missing. Expected:%#v Got:%#v", expectedEventArray, got)
	}
}

func TestTickerHasToBePresentInPortfolioToBeRenamed(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	p.Apply(&sharesAddedEvent)

	err := p.RenameTicker("PG", "FOO")

	_, ok := err.(*portfolio.TickerNotInPortfolioError)
	if !ok {
		t.Errorf("Expected TickerNotInPortfolioError but got %#v", err)
	}
}

func TestTickerCanBeRenamedEvenIfThereAreNoSharesHeld(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 1, 9.99)
	removeSharesEvent := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 1, 9.99)
	p.Apply(&sharesAddedEvent)
	p.Apply(&removeSharesEvent)

	err := p.RenameTicker("MO", "FOO")

	if err != nil {
		t.Errorf("Got unexpected error: %#v", err)
	}
}

func TestNewTickerMustNotBeAlreadyInPortfolio(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 11, 9.99)
	sharesAddedEvent2 := portfolio.NewSharesAddedToPortfolioEvent("PG", 11, 9.99)
	p.Apply(&sharesAddedEvent)
	p.Apply(&sharesAddedEvent2)

	err := p.RenameTicker("MO", "PG")

	_, ok := err.(*portfolio.TickerAlreadyUsedError)
	if !ok {
		t.Errorf("Expected TickerAlreadyUsedError but got %#v", err)
	}
}

func TestNewTickerWillBeUsedForAnyNewPortfolioCommands(t *testing.T) {
	p := portfolio.NewPortfolio()

	sharesAddedEvent := portfolio.NewSharesAddedToPortfolioEvent("MO", 1, 9.99)
	renameEvent := portfolio.NewTickerRenamedEvent("MO", "FOO")
	p.Apply(&sharesAddedEvent)
	p.Apply(&renameEvent)

	err := p.RemoveSharesFromPortfolio("FOO", 1, 9.99)

	if err != nil {
		t.Errorf("Got unexpected error: %#v", err)
	}
}
