package event

import (
	"stock-monitor/domain"
	"stock-monitor/infrastructure"
)

type EventPublisher struct {
	eventStream infrastructure.EventStream
}

func NewEventPublisher(eventStream infrastructure.EventStream) EventPublisher {
	return EventPublisher{eventStream: eventStream}
}

func (publisher *EventPublisher) PublishDomainEvents(events []domain.DomainEvent, occurredAt string) error {
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
