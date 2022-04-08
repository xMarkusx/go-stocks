package portfolio

import (
	"testing"
	"time"
	"reflect"
	"stock-monitor/infrastructure"
)

func TestCanNotBuyZeroShares(t *testing.T) {
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{})
	command := AddSharesToPortfolioCommand("MO", 0, 20.45)
	err := p.AddSharesToPortfolio(command)

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{})
	command := AddSharesToPortfolioCommand("MO", -10, 20.45)
	err := p.AddSharesToPortfolio(command)

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{})
	addSharesCommand := AddSharesToPortfolioCommand("MO", 20, 20.45)
	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 21, 20.45)

	p.AddSharesToPortfolio(addSharesCommand)
	err := p.RemoveSharesFromPortfolio(removeSharesCommand)

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestPortfolioCanBeInitializedWithEvents(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 24.00, "shares": 5}},
	}
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{events})
	got := p.positions
	expected := map[string]Position{
		"MO": {"MO", 5},
		"PG": {"PG", 20},
	}

	if reflect.DeepEqual(got, expected) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, expected)
	}
}

func TestEventsWillBeAddedToEventStream(t *testing.T) {
	eventStream := &infrastructure.InMemoryEventStream{}
	p := ReconstitueFromStream(eventStream)
	addSharesCommand := AddSharesToPortfolioCommand("MO", 20, 20.45)
	addSharesCommand.Date = "2000-01-02"
	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand.Date = "2000-01-02"
	removeSharesCommand2 := RemoveSharesFromPortfolioCommand("MO", 5, 20.45)
	removeSharesCommand2.Date = "2000-01-02"
	p.AddSharesToPortfolio(addSharesCommand)
	p.RemoveSharesFromPortfolio(removeSharesCommand)
	p.RemoveSharesFromPortfolio(removeSharesCommand2)

	got := eventStream.Get()
	want := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 20, "date": "2000-01-02"}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 10, "date": "2000-01-02"}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 5, "date": "2000-01-02"}},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Unexpected event stream. got: %#v, want: %#v", got, want)
	}
}

func TestDateHasToBeInValidFormatWhenAddingShares(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	p := ReconstitueFromStream(&eventStream)
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = "Foo"
	err := p.AddSharesToPortfolio(command)

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenRemovingShares(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	p := ReconstitueFromStream(&eventStream)
	addSharesCommand := AddSharesToPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand.Date = "Foo"

	p.AddSharesToPortfolio(addSharesCommand)
	err := p.RemoveSharesFromPortfolio(removeSharesCommand)

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenAddingShares(t *testing.T) {
	today := time.Now()
	eventStream := infrastructure.InMemoryEventStream{}
	p := ReconstitueFromStream(&eventStream)
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = today.AddDate(0,0,1).Format("2006-01-02")
	err := p.AddSharesToPortfolio(command)

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenRemovingShares(t *testing.T) {
	today := time.Now()
	eventStream := infrastructure.InMemoryEventStream{}
	p := ReconstitueFromStream(&eventStream)
	addSharesCommand := AddSharesToPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand.Date = "Foo"
	removeSharesCommand.Date = today.AddDate(0,0,1).Format("2006-01-02")

	p.AddSharesToPortfolio(addSharesCommand)
	err := p.RemoveSharesFromPortfolio(removeSharesCommand)

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenAddingShares(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10, "date": "2020-01-02"}},
	}
	eventStream := infrastructure.InMemoryEventStream{events}
	p := ReconstitueFromStream(&eventStream)
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = "2020-01-01"
	err := p.AddSharesToPortfolio(command)

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenRemovingShares(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10, "date": "2020-01-02"}},
	}
	eventStream := infrastructure.InMemoryEventStream{events}
	p := ReconstitueFromStream(&eventStream)
	command := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	command.Date = "2020-01-01"
	err := p.RemoveSharesFromPortfolio(command)

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}
