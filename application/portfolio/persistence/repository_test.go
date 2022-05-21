package persistence

import (
	"reflect"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestEventsWillBeAppliedWhenLoadingPortfolio(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{portfolio.SharesRemovedFromPortfolioEventName, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  10.00,
				"date":   "2000-01-01",
			}},
			{portfolio.SharesAddedToPortfolioEventName, map[string]interface{}{
				"ticker": "PG",
				"shares": 20,
				"price":  10.00,
				"date":   "2000-01-02",
			}},
			{portfolio.SharesRemovedFromPortfolioEventName, map[string]interface{}{
				"ticker": "MO",
				"shares": 10,
				"price":  10.00,
				"date":   "2000-01-03",
			}},
		},
	}
	repository := NewEventSourcedPortfolioRepository(&eventStream)

	p := repository.Load()

	expectedPortfolio := portfolio.NewPortfolio()
	event1 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 20, 10.00, "2000-01-01")
	event2 := portfolio.NewSharesAddedToPortfolioEvent("PG", 20, 10.00, "2000-01-02")
	event3 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 10, 10.00, "2000-01-03")
	expectedPortfolio.Apply(&event1)
	expectedPortfolio.Apply(&event2)
	expectedPortfolio.Apply(&event3)

	if reflect.DeepEqual(p, expectedPortfolio) == false {
		t.Errorf("Unexpected portfolio state. Expected:%#v Got:%#v", expectedPortfolio, p)
	}
}

func TestEventsWithoutDateWillBeHandledAsEmptyDate(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{portfolio.SharesAddedToPortfolioEventName, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  10.00,
			}},
		},
	}
	repository := NewEventSourcedPortfolioRepository(&eventStream)

	p := repository.Load()

	expectedPortfolio := portfolio.NewPortfolio()
	event1 := portfolio.NewSharesAddedToPortfolioEvent("MO", 20, 10.00, "")
	expectedPortfolio.Apply(&event1)

	if reflect.DeepEqual(p, expectedPortfolio) == false {
		t.Errorf("Unexpected portfolio state. Expected:%#v Got:%#v", expectedPortfolio, p)
	}
}
