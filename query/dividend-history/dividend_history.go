package dividend_history

import (
	"stock-monitor/domain/dividend"
	"stock-monitor/infrastructure"
	"time"
)

type DividendHistoryQueryInterface interface {
	GetDividends(filter Filter) []Dividend
	GetSum(filter Filter) float32
}

type Dividend struct {
	Ticker string
	Net    float32
	Gross  float32
	Date   string
}

type DividendHistoryQuery struct {
	EventStream  infrastructure.EventStream
	yearFilter   int
	tickerFilter string
}

type Filter struct {
	year   int
	ticker string
}

func NewFilter() Filter {
	return Filter{0, ""}
}

func (filter *Filter) ByYear(year int) {
	filter.year = year
}

func (filter *Filter) ByTicker(ticker string) {
	filter.ticker = ticker
}

func NewDividendHistoryQuery(eventStream infrastructure.EventStream) DividendHistoryQuery {
	return DividendHistoryQuery{eventStream, 0, ""}
}

func (dividendHistoryQuery *DividendHistoryQuery) GetDividends(filter Filter) []Dividend {
	dividends := []Dividend{}

	for _, event := range dividendHistoryQuery.EventStream.Get() {
		if event.Name == dividend.DividendRecordedEventName {
			if !dividendMatchesYearFilter(event.Payload["date"].(string), filter) {
				continue
			}
			ticker := event.Payload["ticker"].(string)
			if !dividendMatchesTickerFilter(ticker, filter) {
				continue
			}
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

func (dividendHistoryQuery *DividendHistoryQuery) GetSum(filter Filter) float32 {
	dividends := float32(0.0)

	for _, event := range dividendHistoryQuery.EventStream.Get() {
		if event.Name == dividend.DividendRecordedEventName {
			if !dividendMatchesYearFilter(event.Payload["date"].(string), filter) {
				continue
			}
			ticker := event.Payload["ticker"].(string)
			if !dividendMatchesTickerFilter(ticker, filter) {
				continue
			}
			dividends += getFloatValue(event.Payload["net"])

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

func dividendMatchesYearFilter(date string, filter Filter) bool {
	dividendDate, _ := time.Parse("2006-01-02", date)
	if filter.year != 0 && dividendDate.Year() != filter.year {
		return false
	}
	return true
}

func dividendMatchesTickerFilter(ticker string, filter Filter) bool {
	if filter.ticker != "" && ticker != filter.ticker {
		return false
	}
	return true
}
