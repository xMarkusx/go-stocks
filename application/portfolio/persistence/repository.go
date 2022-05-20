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
	state := repository.buildState()
	return portfolio.NewPortfolio(&state)
}

func (repository *EventSourcedPortfolioRepository) buildState() portfolio.PortfolioState {
	state := portfolio.NewPortfolioState()
	for _, event := range repository.eventStream.Get() {
		ticker := event.Payload["ticker"].(string)
		shares := event.Payload["shares"].(int)
		date, _ := event.Payload["date"].(string)

		if event.Name == portfolio.SharesAddedToPortfolioEvent {
			state.AddShares(ticker, shares, date)
			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEvent {
			state.RemoveShares(ticker, shares, date)
			continue
		}
	}

	return state
}
