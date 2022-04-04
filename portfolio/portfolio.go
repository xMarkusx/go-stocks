package portfolio

import (
	"errors"
	"stock-monitor/infrastructure"
)

type Position struct {
	ticker string
	shares int
}

func (pos Position) Dto() (string, int) {
	return pos.ticker, pos.shares
}

type Portfolio struct {
	positions    map[string]Position
	eventStream infrastructure.EventStream
}

func ReconstitueFromStream(eventStream infrastructure.EventStream) Portfolio {
	p := Portfolio{map[string]Position{}, eventStream}
	for _, event := range p.eventStream.Get() {
		p.apply(event)
	}
	return p
}

func (portfolio *Portfolio) AddBuyOrder(ticker string, price float32, shares int) error {
	if shares <= 0 {
		return errors.New("number of shares must be greater than 0")
	}

	sharesAddedToPortfolioEvent := infrastructure.Event{
		"Portfolio.SharesAddedToPortfolio",
		map[string]interface{}{
			"ticker": ticker,
			"shares": shares,
			"price": price,
		},
	}

	portfolio.eventStream.Add(sharesAddedToPortfolioEvent)

	portfolio.apply(sharesAddedToPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) AddSellOrder(ticker string, price float32, shares int) error {
	position := portfolio.positions[ticker]
	if position.shares < shares {
		return errors.New("not allowed to sell more shares than currently in portfolio")
	}

	sharesRemovedFromPortfolioEvent := infrastructure.Event{
		"Portfolio.SharesRemovedFromPortfolio",
		map[string]interface{}{
			"ticker": ticker,
			"shares": shares,
			"price": price,
		},
	}

	portfolio.eventStream.Add(sharesRemovedFromPortfolioEvent)

	portfolio.apply(sharesRemovedFromPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) GetPositions() map[string]Position {
	return portfolio.positions
}

func (portfolio *Portfolio) apply(event infrastructure.Event) {
	ticker, shares, _ := extractEventData(event)

	if event.Name == "Portfolio.SharesAddedToPortfolio" {
		position, found := portfolio.positions[ticker]
		if !found {
			portfolio.positions[ticker] = Position{ticker, shares}
		} else {
			position.shares += shares
			portfolio.positions[ticker] = position
		}

		return
	}

	if event.Name == "Portfolio.SharesRemovedFromPortfolio" {
		position := portfolio.positions[ticker]

		position.shares -= shares
		portfolio.positions[ticker] = position
		return
	}
}

func extractEventData(event infrastructure.Event) (string, int, float32) {
	ticker := event.Payload["ticker"].(string)
	shares := event.Payload["shares"].(int)
	price, ok := event.Payload["price"].(float32)
	if !ok {
		price = float32(event.Payload["price"].(float64))
	}

	return ticker, shares, price
}
