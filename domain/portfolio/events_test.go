package portfolio_test

import (
	"reflect"
	"stock-monitor/domain/portfolio"
	"testing"
)

func TestSharesAddedToPortfolioEventCanBeCreated(t *testing.T) {
	event := portfolio.NewSharesAddedToPortfolioEvent("MO", 10, 9.99)

	if event.Name() != portfolio.SharesAddedToPortfolioEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", portfolio.SharesAddedToPortfolioEventName, event.Name())
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
	event := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99)

	if event.Name() != portfolio.SharesRemovedFromPortfolioEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", portfolio.SharesRemovedFromPortfolioEventName, event.Name())
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
	event := portfolio.NewTickerRenamedEvent("MO", "FOO")

	if event.Name() != portfolio.TickerRenamedEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", portfolio.TickerRenamedEvent{}, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"old": "MO",
		"new": "FOO",
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}
