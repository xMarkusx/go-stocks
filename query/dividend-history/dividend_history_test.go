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

	dividendHistoryQuery := dividend_history.DividendHistoryQuery{&infrastructure.InMemoryEventStream{events}}
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

	dividendHistoryQuery := dividend_history.DividendHistoryQuery{&infrastructure.InMemoryEventStream{events}}
	got := dividendHistoryQuery.GetDividends()
	want := []dividend_history.Dividend{
		{"MO", 12.34, 23.45, "2001-01-02"},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Orders unequal got: %#v, want: %#v", got, want)
	}
}
