package portfolio

import (
	"testing"
	"reflect"
	"stock-monitor/infrastructure"
)

func TestCanNotBuyZeroShares(t *testing.T) {
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{})
	err := p.AddBuyOrder("MO", 20.45, 0)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{})
	err := p.AddBuyOrder("MO", 20.45, -10)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	err := p.AddSellOrder("MO", 20.45, 21)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestPortfolioCanBeInitializedWithEvents(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 24.00, "shares": 5}},
	}
	p := ReconstitueFromStream(&infrastructure.InMemoryEventStream{events})
	got := p.GetPositions()
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
	p.AddBuyOrder("MO", 20.45, 20)
	p.AddSellOrder("MO", 20.45, 10)
	p.AddSellOrder("MO", 20.45, 5)

	got := eventStream.Get()
	want := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 10}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 5}},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Unexpected event stream. got: %#v, want: %#v", got, want)
	}
}
