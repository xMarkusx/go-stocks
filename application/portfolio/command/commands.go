package command

import (
	"time"
)

type AddSharesToPortfolioCommand struct {
	Ticker         string
	NumberOfShares int
	Price          float32
	Date           string
}

func NewAddSharesToPortfolioCommand(ticker string, numberOfShares int, price float32) AddSharesToPortfolioCommand {
	today := time.Now().Format("2006-01-02")
	command := AddSharesToPortfolioCommand{ticker, numberOfShares, price, today}

	return command
}

type RemoveSharesFromPortfolioCommand struct {
	Ticker         string
	NumberOfShares int
	Price          float32
	Date           string
}

func NewRemoveSharesFromPortfolioCommand(ticker string, numberOfShares int, price float32) RemoveSharesFromPortfolioCommand {
	today := time.Now().Format("2006-01-02")
	command := RemoveSharesFromPortfolioCommand{ticker, numberOfShares, price, today}

	return command
}

type RenameTickerCommand struct {
	Old  string
	New  string
	Date string
}

func NewRenameTickerCommand(old string, new string) RenameTickerCommand {
	today := time.Now().Format("2006-01-02")
	command := RenameTickerCommand{old, new, today}

	return command
}
