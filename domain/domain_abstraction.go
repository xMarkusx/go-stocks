package domain

type Aggregate interface {
	Apply(event DomainEvent)
	GetRecordedEvents() []DomainEvent
}

type DomainEvent interface {
	Name() string
	Payload() map[string]interface{}
}
