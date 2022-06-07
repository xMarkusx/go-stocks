package dividend_history

import (
	"stock-monitor/domain/dividend"
	"stock-monitor/infrastructure"
)

type DividendHistoryQueryInterface interface {
	GetDividends() []Dividend
}

type Dividend struct {
	Ticker string
	Net    float32
	Gross  float32
	Date   string
}

type DividendHistoryQuery struct {
	EventStream infrastructure.EventStream
}

func (dividendHistoryQuery *DividendHistoryQuery) GetDividends() []Dividend {
	dividends := []Dividend{}

	for _, event := range dividendHistoryQuery.EventStream.Get() {
		if event.Name == dividend.DividendRecordedEventName {
			ticker := event.Payload["ticker"].(string)
			net := getFloatValue(event.Payload["net"])
			gross := getFloatValue(event.Payload["gross"])
			date := event.Payload["date"].(string)
			d := Dividend{ticker, net, gross, date}
			dividends = append(dividends, d)

			continue
		}
	}

	return dividends
}

func getFloatValue(value interface{}) float32 {
	floatValue, ok := value.(float32)
	if !ok {
		floatValue = float32(value.(float64))
	}

	return floatValue
}
