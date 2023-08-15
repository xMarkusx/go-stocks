package command

import (
	"stock-monitor/application/shared"
)

type RecordDividendCommand struct {
	Ticker string
	Net    float32
	Gross  float32
	Date   string
}

func NewRecordDividendCommand(ticker string, net float32, gross float32, date shared.CommandDate) RecordDividendCommand {
	command := RecordDividendCommand{ticker, net, gross, date.Get()}

	return command
}
