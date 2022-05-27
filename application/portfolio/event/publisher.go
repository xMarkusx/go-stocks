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

func (publisher *PortfolioEventPublisher) PublishDomainEvents(events []domain.DomainEvent, occurredAt string) error {
	for _, event := range events {
		genericEvent := infrastructure.Event{
			event.Name(),
			event.Payload(),
			map[string]interface{}{"occurred_at": occurredAt},
		}
		err := publisher.eventStream.Add(genericEvent)
		if err != nil {
			return err
		}
	}

	return nil
}
