package portfolio

import (
	"stock-monitor/domain"
)

type Portfolio struct {
	state  PortfolioState
	events []domain.DomainEvent
}

func NewPortfolio() Portfolio {
	state := NewPortfolioState()
	return Portfolio{state, []domain.DomainEvent{}}
}

func (portfolio *Portfolio) AddSharesToPortfolio(ticker string, shares int, price float32) error {
	if shares <= 0 {
		return &InvalidNumbersOfSharesError{"number of shares must be greater than 0"}
	}

	sharesAddedToPortfolioEvent := NewSharesAddedToPortfolioEvent(ticker, shares, price)
	portfolio.events = append(portfolio.events, &sharesAddedToPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) RemoveSharesFromPortfolio(ticker string, shares int, price float32) error {
	if portfolio.state.GetNumberOfSharesForTicker(ticker) < shares {
		return &CantSellMoreSharesThanExistingError{"not allowed to sell more shares than currently in portfolio"}
	}

	sharesRemovedFromPortfolioEvent := NewSharesRemovedFromPortfolioEvent(ticker, shares, price)
	portfolio.events = append(portfolio.events, &sharesRemovedFromPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) RenameTicker(old string, new string) error {
	_, foundOld := portfolio.state.positions[old]
	if !foundOld {
		return &TickerNotInPortfolioError{"Ticker to be renamed not found. Ticker: " + old}
	}
	_, foundNew := portfolio.state.positions[new]
	if foundNew {
		return &TickerAlreadyUsedError{"New ticker symbol already in use. Ticker: " + new}
	}

	tickerRenamedEvent := NewTickerRenamedEvent(old, new)
	portfolio.events = append(portfolio.events, &tickerRenamedEvent)

	return nil
}

func (portfolio *Portfolio) Apply(event domain.DomainEvent) {
	if event.Name() == SharesAddedToPortfolioEventName {
		sharesAddedToPortfolioEvent := event.(*SharesAddedToPortfolioEvent)
		portfolio.state.AddShares(
			sharesAddedToPortfolioEvent.ticker,
			sharesAddedToPortfolioEvent.shares,
		)
		return
	}
	if event.Name() == SharesRemovedFromPortfolioEventName {
		sharesRemovedFromPortfolioEvent := event.(*SharesRemovedFromPortfolioEvent)
		portfolio.state.RemoveShares(
			sharesRemovedFromPortfolioEvent.ticker,
			sharesRemovedFromPortfolioEvent.shares,
		)
		return
	}
	if event.Name() == TickerRenamedEventName {
		tickerRenamedEvent := event.(*TickerRenamedEvent)
		portfolio.state.positions[tickerRenamedEvent.new] = portfolio.state.positions[tickerRenamedEvent.old]
		delete(portfolio.state.positions, tickerRenamedEvent.old)
	}
}

func (portfolio *Portfolio) GetRecordedEvents() []domain.DomainEvent {
	return portfolio.events
}
