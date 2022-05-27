package persistence_test

import (
	"reflect"
	"stock-monitor/application/portfolio/persistence"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestEventsWillBeAppliedWhenLoadingPortfolio(t *testing.T) {
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
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)

	p := repository.Load()

	expectedPortfolio := portfolio.NewPortfolio()
	event1 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 20, 10.00)
	event2 := portfolio.NewSharesAddedToPortfolioEvent("PG", 20, 10.00)
	event3 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 10, 10.00)
	event4 := portfolio.NewTickerRenamedEvent("MO", "FOO")
	expectedPortfolio.Apply(&event1)
	expectedPortfolio.Apply(&event2)
	expectedPortfolio.Apply(&event3)
	expectedPortfolio.Apply(&event4)

	if reflect.DeepEqual(p, expectedPortfolio) == false {
		t.Errorf("Unexpected portfolio state. Expected:%#v Got:%#v", expectedPortfolio, p)
	}
}
