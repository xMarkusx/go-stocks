package totalInvestedMoney

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestPortfolioGivesTotalAmountOfInvestedMoney(t *testing.T) {
	events := []infrastructure.Event{
		{portfolio.SharesAddedToPortfolioEventName, map[string]interface{}{"ticker": "MO", "price": 20.00, "shares": 20}},
		{portfolio.SharesRemovedFromPortfolioEventName, map[string]interface{}{"ticker": "MO", "price": 30.00, "shares": 10}},
		{portfolio.SharesAddedToPortfolioEventName, map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 5}},
	}

	totalInvestedMoneyQuery := TotalInvestedMoneyQuery{&infrastructure.InMemoryEventStream{events}}

	got := totalInvestedMoneyQuery.GetTotalInvestedMoney()
	expected := float32(300)

	if got != expected {
		t.Errorf("Expected total invested money: %v, got %v", expected, got)
	}
}
