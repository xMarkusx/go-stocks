package query

import (
	"reflect"
	"stock-monitor/infrastructure"
	"testing"
)

func TestProjectionAppliesSharesAddedToPortfolioEvents(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 20}},
	}

	portfolioQuery := RunProjection(&infrastructure.InMemoryEventStream{events})
	got := portfolioQuery.GetPositions()["MO"]
	want := Position{"MO", 30}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}

func TestProjectionAppliesSharesRemovedFromPortfolioEvents(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5}},
	}

	portfolioQuery := RunProjection(&infrastructure.InMemoryEventStream{events})
	got := portfolioQuery.GetPositions()["MO"]
	want := Position{"MO", 5}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#vq, want: %#vq", got, want)
	}
}

func TestPositionIsRemovedWhenCompletelySold(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20}},
	}

	portfolioQuery := RunProjection(&infrastructure.InMemoryEventStream{events})
	_, found := portfolioQuery.GetPositions()["MO"]

	if found {
		t.Errorf("Expected no position of MO in portfolio but found one")
	}
}

func TestPortfolioGivesTotalAmountOfInvestedMoney(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.00, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 30.00, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 5}},
	}

	portfolioQuery := RunProjection(&infrastructure.InMemoryEventStream{events})

	got := portfolioQuery.GetTotalInvestedMoney()
	expected := float32(300)

	if got != expected {
		t.Errorf("Expected total invested money: %v, got %v", expected, got)
	}
}
