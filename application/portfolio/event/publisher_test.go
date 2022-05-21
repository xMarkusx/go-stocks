package event

import (
	"reflect"
	"stock-monitor/domain"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestItPublishesMultipleDomainEvents(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := NewPortfolioEventPublisher(&eventStream)
	event1 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 20, 0.0, "2000-01-01")
	event2 := portfolio.NewSharesAddedToPortfolioEvent("PG", 20, 0.0, "2000-01-02")
	event3 := portfolio.NewSharesRemovedFromPortfolioEvent("MO", 10, 0.0, "2000-01-03")

	publisher.PublishDomainEvents([]domain.DomainEvent{&event1, &event2, &event3})

	expectedEvents := []infrastructure.Event{
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  float32(0.0),
				"date":   "2000-01-01",
			},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{
				"ticker": "PG",
				"shares": 20,
				"price":  float32(0.0),
				"date":   "2000-01-02",
			},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{
				"ticker": "MO",
				"shares": 10,
				"price":  float32(0.0),
				"date":   "2000-01-03",
			},
		},
	}
	got := eventStream.Events

	if reflect.DeepEqual(got, expectedEvents) == false {
		t.Errorf("Unexpected events published. Expected:%#v Got:%#v", expectedEvents, got)
	}
}
