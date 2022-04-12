package portfolio

import (
	"reflect"
	"stock-monitor/infrastructure"
	"testing"
)

func TestCanCreateEmptyEventBasedPortfolioState(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}

	state := NewEventBasedPortfolioState(&eventStream)
	
	expected := EventBasedPortfolioState{map[string]position{}, &eventStream, ""}
	
	if reflect.DeepEqual(state, expected) == false {
		t.Errorf("Unexpected state. Expected:%#v Got:%#v", expected, state)
	}
}

func TestEventsWillBeAppliedWhenCreatingEventBasedPortfolioState(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price": 10.00,
				"date": "2000-01-01",
			}},
			{SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "PG",
				"shares": 20,
				"price": 10.00,
				"date": "2000-01-02",
			}},
			{SharesRemovedFromPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 10,
				"price": 10.00,
				"date": "2000-01-03",
			}},
		},
	}

	state := NewEventBasedPortfolioState(&eventStream)
	
	expected := EventBasedPortfolioState{map[string]position{"MO": {"MO", 10}, "PG": {"PG", 20}}, &eventStream, "2000-01-03"}
	
	if reflect.DeepEqual(state, expected) == false {
		t.Errorf("Unexpected state. Expected:%#v Got:%#v", expected, state)
	}
}

func TestCanGetNumberOfSharesForTicker(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price": 10.00,
				"date": "2000-01-01",
			}},
		},
	}

	state := NewEventBasedPortfolioState(&eventStream)
	
	got := state.GetNumberOfSharesForTicker("MO")
	want := 20
	
	if got != want {
		t.Errorf("Unexpected number of share. Expected:%#v Got:%#v", got, want)
	}
}

func TestCanGetDateOfLastOrder(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price": 10.00,
				"date": "2000-01-01",
			}},
		},
	}

	state := NewEventBasedPortfolioState(&eventStream)
	
	got := state.GetDateOfLastOrder()
	want := "2000-01-01"
	
	if got != want {
		t.Errorf("Unexpected date of last order. Expected:%#v Got:%#v", got, want)
	}
}

func TestItAppliesAddSharesToPortfolioCommandToState(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	state := NewEventBasedPortfolioState(&eventStream)
	
	state.AddShares(addSharesToPortfolioCommand{"MO", 20, 10.00, "2000-01-01"})
	
	numberOfShares := state.GetNumberOfSharesForTicker("MO")
	wantNumberOfShares := 20
	lastOrderDate := state.GetDateOfLastOrder()
	wantLastOrderDate := "2000-01-01"

	if numberOfShares != wantNumberOfShares {
		t.Errorf("Unexpected number of shares. Expected:%#v Got:%#v", wantNumberOfShares, numberOfShares)
	}

	if lastOrderDate != wantLastOrderDate {
		t.Errorf("Unexpected last order date. Expected:%#v Got:%#v", wantLastOrderDate, lastOrderDate)
	}
}

func TestItAppliesRemoveSharesFromPortfolioCommandToState(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price": 10.00,
				"date": "2000-01-01",
			}},
		},
	}
	state := NewEventBasedPortfolioState(&eventStream)
	
	state.RemoveShares(removeSharesFromPortfolioCommand{"MO", 19, 10.00, "2000-01-01"})
	
	numberOfShares := state.GetNumberOfSharesForTicker("MO")
	wantNumberOfShares := 1
	lastOrderDate := state.GetDateOfLastOrder()
	wantLastOrderDate := "2000-01-01"

	if numberOfShares != wantNumberOfShares {
		t.Errorf("Unexpected number of shares. Expected:%#v Got:%#v", wantNumberOfShares, numberOfShares)
	}

	if lastOrderDate != wantLastOrderDate {
		t.Errorf("Unexpected last order date. Expected:%#v Got:%#v", wantLastOrderDate, lastOrderDate)
	}
}

func TestItAddsEventToStreamOnAddSharesToPortfolioCommand(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	state := NewEventBasedPortfolioState(&eventStream)
	
	state.AddShares(addSharesToPortfolioCommand{"MO", 20, 10.00, "2000-01-01"})
	
	got := eventStream.Get()
	want := []infrastructure.Event{
		{SharesAddedToPortfolioEvent, map[string]interface{}{
			"ticker": "MO",
			"shares": 20,
			"price": float32(10.00),
			"date": "2000-01-01",
		}},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Unexpected event stream. Expected:%#v Got:%#v", want, got)
	}
}

func TestItAddsEventToStreamOnRemoveSharesFromPortfolioCommand(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price": float32(10.00),
				"date": "2000-01-01",
			}},
		},
	}
	state := NewEventBasedPortfolioState(&eventStream)
	
	state.RemoveShares(removeSharesFromPortfolioCommand{"MO", 20, 10.00, "2000-01-01"})
	
	got := eventStream.Get()
	want := []infrastructure.Event{
		{SharesAddedToPortfolioEvent, map[string]interface{}{
			"ticker": "MO",
			"shares": 20,
			"price": float32(10.00),
			"date": "2000-01-01",
		}},
		{SharesRemovedFromPortfolioEvent, map[string]interface{}{
			"ticker": "MO",
			"shares": 20,
			"price": float32(10.00),
			"date": "2000-01-01",
		}},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Unexpected event stream. Expected:%#v Got:%#v", want, got)
	}
}
