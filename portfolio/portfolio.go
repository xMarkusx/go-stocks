package portfolio

import "errors"

type Position struct {
	ticker string
	shares int
}

func (pos Position) Dto() (string, int) {
	return pos.ticker, pos.shares
}

func (pos Position) CurrentValue(valueTracker ValueTracker) float32 {
	return valueTracker.Current(pos.ticker) * float32(pos.shares)
}

type Portfolio struct {
	positions    map[string]Position
	eventStream EventStream
}

func InitPortfolio(eventStream EventStream) Portfolio {
	p := Portfolio{map[string]Position{}, eventStream}
	for _, event := range p.eventStream.Get() {
		p.reconstitute(event)
	}
	return p
}

func (portfolio *Portfolio) AddBuyOrder(ticker string, price float32, shares int) error {
	if shares <= 0 {
		return errors.New("number of shares must be greater than 0")
	}

	sharesAddedToPortfolioEvent := Event{
		"Portfolio.SharesAddedToPortfolio",
		map[string]interface{}{
			"ticker": ticker,
			"shares": shares,
			"price": price,
		},
	}

	portfolio.reconstitute(sharesAddedToPortfolioEvent)

	portfolio.eventStream.Add(sharesAddedToPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) AddSellOrder(ticker string, price float32, shares int) error {
	position := portfolio.positions[ticker]
	if position.shares < shares {
		return errors.New("not allowed to sell more shares than currently in portfolio")
	}

	sharesRemovedFromPortfolioEvent := Event{
		"Portfolio.SharesRemovedFromPortfolio",
		map[string]interface{}{
			"ticker": ticker,
			"shares": shares,
			"price": price,
		},
	}

	portfolio.reconstitute(sharesRemovedFromPortfolioEvent)

	portfolio.eventStream.Add(sharesRemovedFromPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) GetPositions() map[string]Position {
	return portfolio.positions
}

func (portfolio *Portfolio) GetTotalInvestedMoney() float32 {
	invested := float32(0.0)
	for _, event := range portfolio.eventStream.Get() {
		_, shares, price := extractEventData(event)

		if event.Name == "Portfolio.SharesAddedToPortfolio" {
			invested += price * float32(shares)
			continue
		}

		if event.Name == "Portfolio.SharesRemovedFromPortfolio" {
			invested -= price * float32(shares)
			continue
		}
	}

	return invested
}

func (portfolio *Portfolio) reconstitute(event Event) {
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

		if position.shares == shares {
			delete(portfolio.positions, ticker)
			return
		}

		position.shares -= shares
		portfolio.positions[ticker] = position
		return
	}
}

func extractEventData(event Event) (string, int, float32) {
	ticker := event.Payload["ticker"].(string)
	shares := event.Payload["shares"].(int)
	price, ok := event.Payload["price"].(float32)
	if !ok {
		price = float32(event.Payload["price"].(float64))
	}

	return ticker, shares, price
}
