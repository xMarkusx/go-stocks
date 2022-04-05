package portfolio

import (
	"time"
)

type addSharesToPortfolioCommand struct {
	Ticker string
	NumberOfShares int
	Price float32
	Date string
}

func AddSharesToPortfolioCommand(ticker string, numberOfShares int, price float32) addSharesToPortfolioCommand {
	today := time.Now().Format("2006-01-02")
	command := addSharesToPortfolioCommand{ticker, numberOfShares, price, today}

	return command
}

type removeSharesFromPortfolioCommand struct {
	Ticker string
	NumberOfShares int
	Price float32
	Date string
}

func RemoveSharesFromPortfolioCommand(ticker string, numberOfShares int, price float32) removeSharesFromPortfolioCommand {
	today := time.Now().Format("2006-01-02")
	command := removeSharesFromPortfolioCommand{ticker, numberOfShares, price, today}

	return command
}
