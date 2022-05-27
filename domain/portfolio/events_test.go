package portfolio

import (
	"reflect"
	"testing"
)

func TestSharesAddedToPortfolioEventCanBeCreated(t *testing.T) {
	event := NewSharesAddedToPortfolioEvent("MO", 10, 9.99)

	if event.Name() != SharesAddedToPortfolioEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", SharesAddedToPortfolioEventName, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"ticker": "MO",
		"shares": 10,
		"price":  float32(9.99),
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}

func TestSharesRemovedFromPortfolioEventCanBeCreated(t *testing.T) {
	event := NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99)

	if event.Name() != SharesRemovedFromPortfolioEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", SharesRemovedFromPortfolioEventName, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"ticker": "MO",
		"shares": 10,
		"price":  float32(9.99),
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}

func TestTickerRenamedEventEventCanBeCreated(t *testing.T) {
	event := NewTickerRenamedEvent("MO", "FOO")

	if event.Name() != TickerRenamedEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", TickerRenamedEvent{}, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"old": "MO",
		"new": "FOO",
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}
