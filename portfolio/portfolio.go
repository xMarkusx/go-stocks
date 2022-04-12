package portfolio

import (
	"time"
)

type Portfolio struct {
	state State
}

type State interface {
	AddShares(command addSharesToPortfolioCommand)
	RemoveShares(command removeSharesFromPortfolioCommand)
	GetNumberOfSharesForTicker(ticker string) int
	GetDateOfLastOrder() string
}

func NewPortfolio(state State) Portfolio {
	return Portfolio{state}
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

	if !commandDateIsLaterThanLastOrderDate(command.Date, portfolio.state.GetDateOfLastOrder()) {
		return &InvalidDateError{"Date can't older than date of last order. Got: " + command.Date}
	}

	portfolio.state.AddShares(command)

	return nil
}

func (portfolio *Portfolio) RemoveSharesFromPortfolio(command removeSharesFromPortfolioCommand) error {
	if portfolio.state.GetNumberOfSharesForTicker(command.Ticker) < command.NumberOfShares {
		return &CantSellMoreSharesThanExistingError{"not allowed to sell more shares than currently in portfolio"}
	}

	if !commandDateHasValidFormat(command.Date) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + command.Date}
	}

	if !dateIsInThePast(command.Date) {
		return &InvalidDateError{"Date can't be in the future. Got: " + command.Date}
	}

	if !commandDateIsLaterThanLastOrderDate(command.Date, portfolio.state.GetDateOfLastOrder()) {
		return &InvalidDateError{"Date can't older than date of last order. Got: " + command.Date}
	}

	portfolio.state.RemoveShares(command)

	return nil
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
