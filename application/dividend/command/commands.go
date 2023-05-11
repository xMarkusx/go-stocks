package command

import (
	"time"
)

type RecordDividendCommand struct {
	Ticker string
	Net    float32
	Gross  float32
	Date   string
}

func NewRecordDividendCommand(ticker string, net float32, gross float32) RecordDividendCommand {
	today := time.Now().Format("2006-01-02")
	command := RecordDividendCommand{ticker, net, gross, today}

	return command
}
