package dividend_history_test

import (
	"reflect"
	"stock-monitor/domain/dividend"
	"stock-monitor/infrastructure"
	"stock-monitor/query/dividend-history"
	"testing"
)

func TestDividendHistoryProvidesAllDividends(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MCD", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	got := dividendHistoryQuery.GetDividends()
	want := []dividend_history.Dividend{
		{"MO", 12.34, 23.45, "2001-01-02"},
		{"PG", 12.34, 23.45, "2001-01-02"},
		{"MCD", 12.34, 23.45, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestDividendHistoryCanHandleFloat32Values(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": float32(12.34), "gross": float32(23.45), "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	got := dividendHistoryQuery.GetDividends()
	want := []dividend_history.Dividend{
		{"MO", 12.34, 23.45, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestDividendHistoryCanBeFilteredByYear(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MCD", "net": 12.34, "gross": 23.45, "date": "2002-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	dividendHistoryQuery.SetYearFilter(2001)
	got := dividendHistoryQuery.GetDividends()
	want := []dividend_history.Dividend{
		{"MO", 12.34, 23.45, "2001-01-02"},
		{"PG", 12.34, 23.45, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestDividendHistoryCanBeFilteredByTicker(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MCD", "net": 12.34, "gross": 23.45, "date": "2002-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	dividendHistoryQuery.SetTickerFilter("PG")
	got := dividendHistoryQuery.GetDividends()
	want := []dividend_history.Dividend{
		{"PG", 12.34, 23.45, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestDividendHistoryCanBeFilteredByTickerAndYear(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 12.34, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 12.34, "gross": 23.45, "date": "2001-02-02"},
			map[string]interface{}{"occurred_at": "2001-02-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 12.34, "gross": 23.45, "date": "2002-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	dividendHistoryQuery.SetTickerFilter("PG")
	dividendHistoryQuery.SetYearFilter(2001)
	got := dividendHistoryQuery.GetDividends()
	want := []dividend_history.Dividend{
		{"PG", 12.34, 23.45, "2001-01-02"},
		{"PG", 12.34, 23.45, "2001-02-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}

func TestDividendHistoryProvidesSumOfDividendsInNet(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 10.01, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 20.02, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MCD", "net": 30.03, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	got := dividendHistoryQuery.GetSum()
	want := float32(60.06)

	if got != want {
		t.Errorf("Dividend sum not matching: %#v, want: %#v", got, want)
	}
}

func TestDividendHistorySumCanBeFilteredByYear(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 10.01, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 20.02, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MCD", "net": 30.03, "gross": 23.45, "date": "2002-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	dividendHistoryQuery.SetYearFilter(2001)
	got := dividendHistoryQuery.GetSum()
	want := float32(30.03)

	if got != want {
		t.Errorf("Dividend sum not matching: %#v, want: %#v", got, want)
	}
}

func TestDividendHistorySumCanBeFilteredByTicker(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 10.01, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 20.02, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MCD", "net": 30.03, "gross": 23.45, "date": "2002-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	dividendHistoryQuery.SetTickerFilter("PG")
	got := dividendHistoryQuery.GetSum()
	want := float32(20.02)

	if got != want {
		t.Errorf("Dividend sum not matching: %#v, want: %#v", got, want)
	}
}

func TestDividendHistorySumCanBeFilteredByTickerAndYear(t *testing.T) {
	events := []infrastructure.Event{
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "MO", "net": 10.01, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 20.02, "gross": 23.45, "date": "2001-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 31.02, "gross": 23.45, "date": "2001-02-02"},
			map[string]interface{}{"occurred_at": "2001-02-02"},
		},
		{
			dividend.DividendRecordedEventName,
			map[string]interface{}{"ticker": "PG", "net": 40.04, "gross": 23.45, "date": "2002-01-02"},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	dividendHistoryQuery := dividend_history.NewDividendHistoryQuery(&infrastructure.InMemoryEventStream{events})
	dividendHistoryQuery.SetTickerFilter("PG")
	dividendHistoryQuery.SetYearFilter(2001)
	got := dividendHistoryQuery.GetSum()
	want := float32(51.04)

	if got != want {
		t.Errorf("Dividend sum not matching: %#v, want: %#v", got, want)
	}
}
