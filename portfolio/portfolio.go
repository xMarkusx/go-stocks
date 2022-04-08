package portfolio

import (
	"stock-monitor/infrastructure"
	"time"
)

type Position struct {
	ticker string
	shares int
}

func (pos Position) Dto() (string, int) {
	return pos.ticker, pos.shares
}

type Portfolio struct {
	positions map[string]Position
	eventStream infrastructure.EventStream
	lastOrderDate string
}

func ReconstitueFromStream(eventStream infrastructure.EventStream) Portfolio {
	p := Portfolio{map[string]Position{}, eventStream, ""}
	for _, event := range p.eventStream.Get() {
		p.apply(event)
	}
	return p
}

func (portfolio *Portfolio) AddSharesToPortfolio(command addSharesToPortfolioCommand) error {
	if command.NumberOfShares <= 0 {
		return &InvalidNumbersOfSharesError{"number of shares must be greater than 0"}
	}
	if !commandDateHasValidFormat(command.Date) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + command.Date}
	}

	if !dateIsInThePast(command.Date) {
		return &InvalidDateError{"Date can't be in the future. Got: " + command.Date}
	}

	if !commandDateIsLaterThanLastOrderDate(command.Date, portfolio.lastOrderDate) {
		return &InvalidDateError{"Date can't older than date of last order. Got: " + command.Date}
	}

	sharesAddedToPortfolioEvent := infrastructure.Event{
		"Portfolio.SharesAddedToPortfolio",
		map[string]interface{}{
			"ticker": command.Ticker,
			"shares": command.NumberOfShares,
			"price": command.Price,
			"date": command.Date,
		},
	}

	portfolio.eventStream.Add(sharesAddedToPortfolioEvent)

	portfolio.apply(sharesAddedToPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) RemoveSharesFromPortfolio(command removeSharesFromPortfolioCommand) error {
	position := portfolio.positions[command.Ticker]
	if position.shares < command.NumberOfShares {
		return &CantSellMoreSharesThanExistingError{"not allowed to sell more shares than currently in portfolio"}
	}

	if !commandDateHasValidFormat(command.Date) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + command.Date}
	}

	if !dateIsInThePast(command.Date) {
		return &InvalidDateError{"Date can't be in the future. Got: " + command.Date}
	}

	if !commandDateIsLaterThanLastOrderDate(command.Date, portfolio.lastOrderDate) {
		return &InvalidDateError{"Date can't older than date of last order. Got: " + command.Date}
	}

	sharesRemovedFromPortfolioEvent := infrastructure.Event{
		"Portfolio.SharesRemovedFromPortfolio",
		map[string]interface{}{
			"ticker": command.Ticker,
			"shares": command.NumberOfShares,
			"price": command.Price,
			"date": command.Date,
		},
	}

	portfolio.eventStream.Add(sharesRemovedFromPortfolioEvent)

	portfolio.apply(sharesRemovedFromPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) apply(event infrastructure.Event) {
	ticker, shares, _ := extractEventData(event)

	portfolio.lastOrderDate, _ = event.Payload["date"].(string)

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

func commandDateHasValidFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	return true
}

func dateIsInThePast(date string) bool {
	commandDate, _ := time.Parse("2006-01-02", date)
	today := time.Now()
	diff := today.Sub(commandDate)

	if diff < 0 {
		return false
	}

	return true
}

func commandDateIsLaterThanLastOrderDate(commandDate string, lastOrderDate string) bool {
	cd, _ := time.Parse("2006-01-02", commandDate)
	lod, _ := time.Parse("2006-01-02", lastOrderDate)
	diff := cd.Sub(lod)
	if diff < 0 {
		return false
	}

	return true
}
