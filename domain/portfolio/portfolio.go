package portfolio

import (
	"stock-monitor/domain"
	"time"
)

type Portfolio struct {
	state  State
	events []domain.DomainEvent
}

type State interface {
	AddShares(ticker string, shares int, date string)
	RemoveShares(ticker string, shares int, date string)
	GetNumberOfSharesForTicker(ticker string) int
	GetDateOfLastOrder() string
}

func NewPortfolio() Portfolio {
	state := NewPortfolioState()
	return Portfolio{&state, []domain.DomainEvent{}}
}

func (portfolio *Portfolio) AddSharesToPortfolio(ticker string, shares int, price float32, date string) error {
	if shares <= 0 {
		return &InvalidNumbersOfSharesError{"number of shares must be greater than 0"}
	}
	if !commandDateHasValidFormat(date) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + date}
	}

	if !dateIsInThePast(date) {
		return &InvalidDateError{"Date can't be in the future. Got: " + date}
	}

	if !commandDateIsLaterThanLastOrderDate(date, portfolio.state.GetDateOfLastOrder()) {
		return &InvalidDateError{"Date can't older than date of last order. Got: " + date}
	}

	sharesAddedToPortfolioEvent := NewSharesAddedToPortfolioEvent(ticker, shares, price, date)
	portfolio.events = append(portfolio.events, &sharesAddedToPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) RemoveSharesFromPortfolio(ticker string, shares int, price float32, date string) error {
	if portfolio.state.GetNumberOfSharesForTicker(ticker) < shares {
		return &CantSellMoreSharesThanExistingError{"not allowed to sell more shares than currently in portfolio"}
	}

	if !commandDateHasValidFormat(date) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + date}
	}

	if !dateIsInThePast(date) {
		return &InvalidDateError{"Date can't be in the future. Got: " + date}
	}

	if !commandDateIsLaterThanLastOrderDate(date, portfolio.state.GetDateOfLastOrder()) {
		return &InvalidDateError{"Date can't older than date of last order. Got: " + date}
	}

	sharesRemovedFromPortfolioEvent := NewSharesRemovedFromPortfolioEvent(ticker, shares, price, date)
	portfolio.events = append(portfolio.events, &sharesRemovedFromPortfolioEvent)

	return nil
}

func (portfolio *Portfolio) Apply(event domain.DomainEvent) {
	if event.Name() == SharesAddedToPortfolioEventName {
		sharesAddedToPortfolioEvent := event.(*SharesAddedToPortfolioEvent)
		portfolio.state.AddShares(
			sharesAddedToPortfolioEvent.ticker,
			sharesAddedToPortfolioEvent.shares,
			sharesAddedToPortfolioEvent.date,
		)
		return
	}
	if event.Name() == SharesRemovedFromPortfolioEventName {
		sharesRemovedFromPortfolioEvent := event.(*SharesRemovedFromPortfolioEvent)
		portfolio.state.RemoveShares(
			sharesRemovedFromPortfolioEvent.ticker,
			sharesRemovedFromPortfolioEvent.shares,
			sharesRemovedFromPortfolioEvent.date,
		)
	}
}

func (portfolio *Portfolio) GetRecordedEvents() []domain.DomainEvent {
	return portfolio.events
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
