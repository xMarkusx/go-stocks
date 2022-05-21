package event

import (
	"stock-monitor/domain"
	"stock-monitor/infrastructure"
)

type PortfolioEventPublisher struct {
	eventStream infrastructure.EventStream
}

func NewPortfolioEventPublisher(eventStream infrastructure.EventStream) PortfolioEventPublisher {
	return PortfolioEventPublisher{eventStream: eventStream}
}

func (publisher *PortfolioEventPublisher) PublishDomainEvents(events []domain.DomainEvent) {
	for _, event := range events {
		genericEvent := infrastructure.Event{
			event.Name(),
			event.Payload(),
		}
		publisher.eventStream.Add(genericEvent)
	}
}
