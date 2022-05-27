package orderHistory_test

import (
	"reflect"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	orderHistory "stock-monitor/query/order-history"
	"testing"
)

func TestOrderHistoryProvidesAllOrders(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "PG", "price": 40.00, "shares": 20},
			map[string]interface{}{"occurred_at": "2001-01-03"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2001-01-04"},
		},
	}

	orderHistoryQuery := orderHistory.OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []orderHistory.Order{
		{"BUY", "MO", []string{}, 10, 20.45, "2001-01-02"},
		{"BUY", "PG", []string{}, 20, 40.00, "2001-01-03"},
		{"SELL", "MO", []string{}, 5, 40.00, "2001-01-04"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestOrderHistoryCanHandleOrdersWithoutDate(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{},
		},
	}

	orderHistoryQuery := orderHistory.OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []orderHistory.Order{
		{"BUY", "MO", []string{}, 10, 20.45, ""},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestOrderHistoryContainsRenamedTickers(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "MO", "new": "FOO"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	orderHistoryQuery := orderHistory.OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []orderHistory.Order{
		{"BUY", "FOO", []string{"MO"}, 10, 20.45, "2001-01-02"},
		{"SELL", "FOO", []string{"MO"}, 5, 40.00, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestMultipleRenamesOnTheSameTickerAreHandled(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "MO", "new": "FOO"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "FOO", "new": "BAR"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	orderHistoryQuery := orderHistory.OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []orderHistory.Order{
		{"BUY", "BAR", []string{"MO", "FOO"}, 10, 20.45, "2001-01-02"},
		{"SELL", "BAR", []string{"MO", "FOO"}, 5, 40.00, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestOrdersAfterRenameDoNotContainPreviousTickers(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "MO", "new": "FOO"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "FOO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-03"},
		},
	}

	orderHistoryQuery := orderHistory.OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []orderHistory.Order{
		{"BUY", "FOO", []string{"MO"}, 10, 20.45, "2001-01-02"},
		{"SELL", "FOO", []string{"MO"}, 5, 40.00, "2001-01-02"},
		{"BUY", "FOO", []string{}, 10, 20.45, "2001-01-03"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}
