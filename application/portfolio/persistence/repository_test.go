package persistence

import (
	"reflect"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestCanCreateRepository(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}

	repository := NewEventSourcedPortfolioRepository(&eventStream)

	expected := EventSourcedPortfolioRepository{eventStream: &eventStream}

	if reflect.DeepEqual(repository, expected) == false {
		t.Errorf("Unexpected repository. Expected:%#v Got:%#v", expected, repository)
	}
}

func TestEventsWillBeAppliedWhenLoadingPortfolio(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{portfolio.SharesRemovedFromPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  10.00,
				"date":   "2000-01-01",
			}},
			{portfolio.SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "PG",
				"shares": 20,
				"price":  10.00,
				"date":   "2000-01-02",
			}},
			{portfolio.SharesRemovedFromPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 10,
				"price":  10.00,
				"date":   "2000-01-03",
			}},
		},
	}
	repository := NewEventSourcedPortfolioRepository(&eventStream)

	p := repository.Load()

	expectedState := portfolio.NewPortfolioState()
	expectedState.RemoveShares("MO", 20, "2000-01-01")
	expectedState.AddShares("PG", 20, "2000-01-01")
	expectedState.RemoveShares("MO", 10, "2000-01-03")
	expectedPortfolio := portfolio.NewPortfolio(&expectedState)

	if reflect.DeepEqual(p, expectedPortfolio) == false {
		t.Errorf("Unexpected portfolio state. Expected:%#v Got:%#v", expectedPortfolio, p)
	}
}

func TestEventsWithoutDateWillBeHandledAsEmptyDate(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{portfolio.SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  10.00,
			}},
		},
	}
	repository := NewEventSourcedPortfolioRepository(&eventStream)

	p := repository.Load()

	expectedState := portfolio.NewPortfolioState()
	expectedState.AddShares("MO", 20, "")
	expectedPortfolio := portfolio.NewPortfolio(&expectedState)

	if reflect.DeepEqual(p, expectedPortfolio) == false {
		t.Errorf("Unexpected portfolio state. Expected:%#v Got:%#v", expectedPortfolio, p)
	}
}
