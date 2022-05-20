package event

import (
	"reflect"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestSharesAddedToPortfolioEventCanBePublished(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := NewPortfolioEventPublisher(&eventStream)
	addSharesCommand := command.NewAddSharesToPortfolioCommand("MO", 10, 9.99)
	addSharesCommand.Date = "2000-01-01"

	publisher.PublishSharesAddedToPortfolioEvent(addSharesCommand)

	expectedEvent := infrastructure.Event{
		portfolio.SharesAddedToPortfolioEvent,
		map[string]interface{}{
			"ticker": "MO",
			"shares": 10,
			"price":  float32(9.99),
			"date":   "2000-01-01",
		},
	}
	got := eventStream.Events[0]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestSharesRemovedFromPortfolioEventCanBePublished(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := NewPortfolioEventPublisher(&eventStream)
	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 10, 9.99)
	removeSharesCommand.Date = "2000-01-01"

	publisher.PublishSharesRemovedFromPortfolioEvent(removeSharesCommand)

	expectedEvent := infrastructure.Event{
		portfolio.SharesRemovedFromPortfolioEvent,
		map[string]interface{}{
			"ticker": "MO",
			"shares": 10,
			"price":  float32(9.99),
			"date":   "2000-01-01",
		},
	}
	got := eventStream.Events[0]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}
