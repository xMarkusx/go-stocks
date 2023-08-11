package command

import (
	"stock-monitor/application/shared"
)

type AddSharesToPortfolioCommand struct {
	Ticker         string
	NumberOfShares int
	Price          float32
	Date           string
}

func NewAddSharesToPortfolioCommand(ticker string, numberOfShares int, price float32, date shared.CommandDate) AddSharesToPortfolioCommand {
	command := AddSharesToPortfolioCommand{ticker, numberOfShares, price, date.Get()}

	return command
}

type RemoveSharesFromPortfolioCommand struct {
	Ticker         string
	NumberOfShares int
	Price          float32
	Date           string
}

func NewRemoveSharesFromPortfolioCommand(ticker string, numberOfShares int, price float32, date shared.CommandDate) RemoveSharesFromPortfolioCommand {
	command := RemoveSharesFromPortfolioCommand{ticker, numberOfShares, price, date.Get()}

	return command
}

type RenameTickerCommand struct {
	Old  string
	New  string
	Date string
}

func NewRenameTickerCommand(old string, new string, date shared.CommandDate) RenameTickerCommand {
	command := RenameTickerCommand{old, new, date.Get()}

	return command
}
