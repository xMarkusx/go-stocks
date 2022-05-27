package totalInvestedMoney_test

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	totalInvestedMoney "stock-monitor/query/total-invested-money"
	"testing"
)

func TestPortfolioGivesTotalAmountOfInvestedMoney(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.00, "shares": 20},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 30.00, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	totalInvestedMoneyQuery := totalInvestedMoney.TotalInvestedMoneyQuery{&infrastructure.InMemoryEventStream{events}}

	got := totalInvestedMoneyQuery.GetTotalInvestedMoney()
	expected := float32(300)

	if got != expected {
		t.Errorf("Expected total invested money: %v, got %v", expected, got)
	}
}
