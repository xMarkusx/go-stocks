package persistence

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
)

type PortfolioRepository interface {
	Load() portfolio.Portfolio
}

type EventSourcedPortfolioRepository struct {
	eventStream infrastructure.EventStream
}

func NewEventSourcedPortfolioRepository(eventStream infrastructure.EventStream) EventSourcedPortfolioRepository {
	return EventSourcedPortfolioRepository{eventStream: eventStream}
}

func (repository *EventSourcedPortfolioRepository) Load() portfolio.Portfolio {
	p := portfolio.NewPortfolio()
	for _, event := range repository.eventStream.Get() {
		ticker := event.Payload["ticker"].(string)
		shares := event.Payload["shares"].(int)
		date, _ := event.Payload["date"].(string)

		if event.Name == portfolio.SharesAddedToPortfolioEventName {
			domainEvent := portfolio.NewSharesAddedToPortfolioEvent(ticker, shares, 0.0, date)
			p.Apply(&domainEvent)
			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEventName {
			domainEvent := portfolio.NewSharesRemovedFromPortfolioEvent(ticker, shares, 0.0, date)
			p.Apply(&domainEvent)
			continue
		}
	}
	return p
}
