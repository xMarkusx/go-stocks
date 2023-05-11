package event_test

import (
	"reflect"
	"stock-monitor/application/event"
	"stock-monitor/domain"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestItPublishesMultipleDomainEvents(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewEventPublisher(&eventStream)
	event1 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 20, 9.99)
	event2 := portfolio.NewSharesAddedToPortfolioEvent("PG", 20, 9.99, "2000-01-01")
	event3 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 10, 9.99)

	publisher.PublishDomainEvents([]domain.DomainEvent{&event1, &event2, &event3}, "2000-01-01")

	expectedEvents := []infrastructure.Event{
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  float32(9.99),
			},
			map[string]interface{}{"occurred_at": "2000-01-01"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{
				"ticker": "PG",
				"shares": 20,
				"price":  float32(9.99),
				"date":   "2000-01-01",
			},
			map[string]interface{}{"occurred_at": "2000-01-01"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{
				"ticker": "MO",
				"shares": 10,
				"price":  float32(9.99),
			},
			map[string]interface{}{"occurred_at": "2000-01-01"},
		},
	}
	got := eventStream.Events

	if reflect.DeepEqual(got, expectedEvents) == false {
		t.Errorf("Unexpected events published. Expected:%#v Got:%#v", expectedEvents, got)
	}
}

func TestItThrowsAnErrorIfAddingToEventStreamFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewEventPublisher(&eventStream)
	event1 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 20, 9.99)

	err := publisher.PublishDomainEvents([]domain.DomainEvent{&event1}, "FOO")

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}
