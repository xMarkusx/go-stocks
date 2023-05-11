package dividend

const DividendRecordedEventName = "Dividend.DividendRecorded"

type DividendRecordedEvent struct {
	ticker string
	net    float32
	gross  float32
	date   string
}

func NewDividendRecordedEvent(ticker string, net float32, gross float32, date string) DividendRecordedEvent {
	return DividendRecordedEvent{ticker: ticker, net: net, gross: gross, date: date}
}

func (event *DividendRecordedEvent) Name() string {
	return DividendRecordedEventName
}

func (event *DividendRecordedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"ticker": event.ticker,
		"net":    event.net,
		"gross":  event.gross,
		"date":   event.date,
	}
}
