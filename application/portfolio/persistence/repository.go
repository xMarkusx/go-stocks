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

		if event.Name == portfolio.SharesAddedToPortfolioEventName {
			ticker := event.Payload["ticker"].(string)
			shares := event.Payload["shares"].(int)
			domainEvent := portfolio.NewSharesAddedToPortfolioEvent(ticker, shares, 0.0)
			p.Apply(&domainEvent)
			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEventName {
			ticker := event.Payload["ticker"].(string)
			shares := event.Payload["shares"].(int)
			domainEvent := portfolio.NewSharesRemovedFromPortfolioEvent(ticker, shares, 0.0)
			p.Apply(&domainEvent)
			continue
		}

		if event.Name == portfolio.TickerRenamedEventName {
			oldSymbol := event.Payload["old"].(string)
			newSymbol := event.Payload["new"].(string)
			domainEvent := portfolio.NewTickerRenamedEvent(oldSymbol, newSymbol)
			p.Apply(&domainEvent)
			continue
		}
	}
	return p
}
