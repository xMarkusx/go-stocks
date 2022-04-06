package totalInvestedMoney

import (
	"stock-monitor/infrastructure"
	"testing"
)

func TestPortfolioGivesTotalAmountOfInvestedMoney(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.00, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 30.00, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 5}},
	}

	totalInvestedMoneyQuery := TotalInvestedMoneyQuery{&infrastructure.InMemoryEventStream{events}}

	got := totalInvestedMoneyQuery.GetTotalInvestedMoney()
	expected := float32(300)

	if got != expected {
		t.Errorf("Expected total invested money: %v, got %v", expected, got)
	}
}
