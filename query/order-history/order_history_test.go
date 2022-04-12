package orderHistory

import (
	"reflect"
	"stock-monitor/infrastructure"
	"stock-monitor/portfolio"
	"testing"
)

func TestOrderHistoryProvidesAllOrders(t *testing.T) {
	events := []infrastructure.Event{
		{portfolio.SharesAddedToPortfolioEvent, map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10, "date": "2001-01-02"}},
		{portfolio.SharesAddedToPortfolioEvent, map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 20, "date": "2001-02-02"}},
		{portfolio.SharesRemovedFromPortfolioEvent, map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5, "date": "2002-01-02"}},
	}

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
		{"BUY", "MO", 10, 20.45, "2001-01-02"},
		{"BUY", "PG", 20, 40.00, "2001-02-02"},
		{"SELL", "MO", 5, 40.00, "2002-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestOrderHistoryCanHandleOrdersWithoutDate(t *testing.T) {
	events := []infrastructure.Event{
		{portfolio.SharesAddedToPortfolioEvent, map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
	}

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
		{"BUY", "MO", 10, 20.45, ""},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}
