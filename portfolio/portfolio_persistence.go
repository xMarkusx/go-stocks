package portfolio

import (
	"stock-monitor/infrastructure"
)

type EventBasedPortfolioState struct {
	positions map[string]position
	eventStream infrastructure.EventStream
	lastOrderDate string
}

type position struct {
	ticker string
	shares int
}

func NewEventBasedPortfolioState(eventStream infrastructure.EventStream) EventBasedPortfolioState {
	p := EventBasedPortfolioState{map[string]position{}, eventStream, ""}
	for _, event := range p.eventStream.Get() {
		p.apply(event)
	}
	return p
}

func (portfolioState *EventBasedPortfolioState) GetNumberOfSharesForTicker(ticker string) int {
	return portfolioState.positions[ticker].shares
}

func (portfolioState *EventBasedPortfolioState) GetDateOfLastOrder() string {
	return portfolioState.lastOrderDate
}

func (portfolioState *EventBasedPortfolioState) AddShares(command addSharesToPortfolioCommand) {
	sharesAddedToPortfolioEvent := infrastructure.Event{
		SharesAddedToPortfolioEvent,
		map[string]interface{}{
			"ticker": command.Ticker,
			"shares": command.NumberOfShares,
			"price": command.Price,
			"date": command.Date,
		},
	}

	portfolioState.eventStream.Add(sharesAddedToPortfolioEvent)
	portfolioState.apply(sharesAddedToPortfolioEvent)
}

func (portfolioState *EventBasedPortfolioState) RemoveShares(command removeSharesFromPortfolioCommand) {
	sharesRemovedFromPortfolioEvent := infrastructure.Event{
		SharesRemovedFromPortfolioEvent,
		map[string]interface{}{
			"ticker": command.Ticker,
			"shares": command.NumberOfShares,
			"price": command.Price,
			"date": command.Date,
		},
	}

	portfolioState.eventStream.Add(sharesRemovedFromPortfolioEvent)
	portfolioState.apply(sharesRemovedFromPortfolioEvent)
}

func (portfolioState *EventBasedPortfolioState) apply(event infrastructure.Event) {
	ticker := event.Payload["ticker"].(string)
	shares := event.Payload["shares"].(int)

	portfolioState.lastOrderDate, _ = event.Payload["date"].(string)

	if event.Name == SharesAddedToPortfolioEvent {
		p, found := portfolioState.positions[ticker]
		if !found {
			portfolioState.positions[ticker] = position{ticker, shares}
		} else {
			p.shares += shares
			portfolioState.positions[ticker] = p
		}

		return
	}

	if event.Name == SharesRemovedFromPortfolioEvent {
		p := portfolioState.positions[ticker]

		p.shares -= shares
		portfolioState.positions[ticker] = p

		return
	}
}
