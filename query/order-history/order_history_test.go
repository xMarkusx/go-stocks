package orderHistory

import (
	"reflect"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
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

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
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

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
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

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
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

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
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

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := orderHistoryQuery.GetOrders()
	want := []Order{
		{"BUY", "FOO", []string{"MO"}, 10, 20.45, "2001-01-02"},
		{"SELL", "FOO", []string{"MO"}, 5, 40.00, "2001-01-02"},
		{"BUY", "FOO", []string{}, 10, 20.45, "2001-01-03"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestOrderDto(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	orderHistoryQuery := OrderHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	orderType, ticker, aliases, shares, price, date := orderHistoryQuery.GetOrders()[0].Dto()

	if orderType != "BUY" {
		t.Errorf("Unexpected orderType from dto output. got: %#v, want: %#v", orderType, "BUY")
	}
	if ticker != "MO" {
		t.Errorf("Unexpected ticker from dto output. got: %#v, want: %#v", ticker, "MO")
	}
	if reflect.DeepEqual(aliases, []string{}) == false {
		t.Errorf("Unexpected aliases from dto output. got: %#v, want: %#v", aliases, []string{})
	}
	if shares != 10 {
		t.Errorf("Unexpected shares from dto output. got: %#v, want: %#v", shares, 10)
	}
	if price != 20.45 {
		t.Errorf("Unexpected price from dto output. got: %#v, want: %#v", price, 20.45)
	}
	if date != "2001-01-02" {
		t.Errorf("Unexpected date from dto output. got: %#v, want: %#v", date, "2001-01-02")
	}
}
