package persistence_test

import (
	"reflect"
	"stock-monitor/application/dividend/persistence"
	"stock-monitor/domain/dividend"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestPortfolioEventsWillBeAppliedWhenLoadingDividend(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesRemovedFromPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "PG",
					"shares": 20,
					"price":  10.00,
					"date":   "2000-01-01",
				},
				map[string]interface{}{"occurred_at": "2000-01-02"},
			},
			{
				portfolio.SharesRemovedFromPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 10,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-03"},
			},
			{
				portfolio.TickerRenamedEventName,
				map[string]interface{}{
					"old": "MO",
					"new": "FOO",
				},
				map[string]interface{}{"occurred_at": "2000-01-04"},
			},
		},
	}
	repository := persistence.NewEventSourcedDividendRepository(&eventStream)

	d := repository.Load()

	expectedDividend := dividend.NewDividend()
	event1 := portfolio.NewSharesAddedToPortfolioEvent("PG", 20, 10.00, "2000-01-01")
	event2 := portfolio.NewTickerRenamedEvent("MO", "FOO")
	expectedDividend.Apply(&event1)
	expectedDividend.Apply(&event2)

	if reflect.DeepEqual(d, expectedDividend) == false {
		t.Errorf("Unexpected portfolio state. Expected:%#v Got:%#v", expectedDividend, d)
	}
}
