package portfolio

import "testing"
import "reflect"

func TestBuyOrder(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	got := p.GetPositions()["MO"]
	want := Position{"MO", 20}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}

func TestMultipleBuyOrderForSameTicker(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	p.AddBuyOrder("MO", 30.45, 20)
	got := p.GetPositions()["MO"]
	want := Position{"MO", 40}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	err := p.AddBuyOrder("MO", 20.45, 0)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	err := p.AddBuyOrder("MO", 20.45, -10)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestSellOrder(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	p.AddSellOrder("MO", 20.45, 10)
	got := p.GetPositions()["MO"]
	want := Position{"MO", 10}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}

func TestMultipleSellOrdersOnSamePosition(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	p.AddSellOrder("MO", 20.45, 10)
	p.AddSellOrder("MO", 20.45, 5)
	got := p.GetPositions()["MO"]
	want := Position{"MO", 5}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#vq, want: %#vq", got, want)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	err := p.AddSellOrder("MO", 20.45, 21)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestPositionIsRemovedWhenCompletelySold(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20.45, 20)
	p.AddSellOrder("MO", 20.45, 20)
	_, found := p.GetPositions()["MO"]

	if found {
		t.Errorf("Expected no position of MO in portfolio but found one")
	}
}

func TestPortfolioGivesTotalAmountOfInvestedMoney(t *testing.T) {
	p := InitPortfolio(&InMemoryEventStream{})
	p.AddBuyOrder("MO", 20, 20)
	p.AddSellOrder("MO", 30, 10)
	p.AddBuyOrder("PG", 40, 5)

	got := p.GetTotalInvestedMoney()
	expected := float32(300)

	if got != expected {
		t.Errorf("Expected total invested money: %v, got %v", expected, got)
	}
}

func TestPortfolioCanBeInitializedWithEvents(t *testing.T) {
	events := []Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 24.00, "shares": 5}},
	}
	p := InitPortfolio(&InMemoryEventStream{events})
	got := p.GetPositions()
	expected := map[string]Position{
		"MO": {"MO", 5},
		"PG": {"PG", 20},
	}

	if reflect.DeepEqual(got, expected) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, expected)
	}
}

func TestOrdersWillBeAddedToStorage(t *testing.T) {
	os := &InMemoryEventStream{}
	p := InitPortfolio(os)
	p.AddBuyOrder("MO", 20.45, 20)
	p.AddSellOrder("MO", 20.45, 10)
	p.AddSellOrder("MO", 20.45, 5)

	got := os.Get()
	want := []Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 10}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": float32(20.45), "shares": 5}},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Unexpected event stream. got: %#v, want: %#v", got, want)
	}
}
