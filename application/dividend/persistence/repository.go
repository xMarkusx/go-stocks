package persistence

import (
	"stock-monitor/domain/dividend"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
)

type DividendRepository interface {
	Load() dividend.Dividend
}

type EventSourcedDividendRepository struct {
	portfolioEventStream infrastructure.EventStream
}

func NewEventSourcedDividendRepository(eventStream infrastructure.EventStream) EventSourcedDividendRepository {
	return EventSourcedDividendRepository{portfolioEventStream: eventStream}
}

func (repository *EventSourcedDividendRepository) Load() dividend.Dividend {
	d := dividend.NewDividend()
	for _, event := range repository.portfolioEventStream.Get() {

		if event.Name == portfolio.SharesAddedToPortfolioEventName {
			ticker := event.Payload["ticker"].(string)
			shares := event.Payload["shares"].(int)
			date := event.Payload["date"].(string)
			domainEvent := portfolio.NewSharesAddedToPortfolioEvent(ticker, shares, 0.0, date)
			d.Apply(&domainEvent)
			continue
		}

		if event.Name == portfolio.TickerRenamedEventName {
			oldSymbol := event.Payload["old"].(string)
			newSymbol := event.Payload["new"].(string)
			domainEvent := portfolio.NewTickerRenamedEvent(oldSymbol, newSymbol)
			d.Apply(&domainEvent)
			continue
		}
	}
	return d
}
