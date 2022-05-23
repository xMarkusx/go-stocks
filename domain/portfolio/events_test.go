package portfolio

import (
	"reflect"
	"testing"
)

func TestSharesAddedToPortfolioEventCanBeCreated(t *testing.T) {
	event := NewSharesAddedToPortfolioEvent("MO", 10, 9.99, "2000-01-01")

	if event.Name() != SharesAddedToPortfolioEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", SharesAddedToPortfolioEventName, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"ticker": "MO",
		"shares": 10,
		"price":  float32(9.99),
		"date":   "2000-01-01",
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}

func TestSharesRemovedFromPortfolioEventCanBeCreated(t *testing.T) {
	event := NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99, "2000-01-01")

	if event.Name() != SharesRemovedFromPortfolioEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", SharesRemovedFromPortfolioEventName, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"ticker": "MO",
		"shares": 10,
		"price":  float32(9.99),
		"date":   "2000-01-01",
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}

func TestTickerRenamedEventEventCanBeCreated(t *testing.T) {
	event := NewTickerRenamedEvent("MO", "FOO", "2000-01-01")

	if event.Name() != TickerRenamedEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", TickerRenamedEvent{}, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"old":  "MO",
		"new":  "FOO",
		"date": "2000-01-01",
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}
