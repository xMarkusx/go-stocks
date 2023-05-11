package dividend_history

import (
	"stock-monitor/domain/dividend"
	"stock-monitor/infrastructure"
	"time"
)

type DividendHistoryQueryInterface interface {
	GetDividends() []Dividend
	GetSum() float32
	SetYearFilter(year int)
	SetTickerFilter(ticker string)
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

func NewDividendHistoryQuery(eventStream infrastructure.EventStream) DividendHistoryQuery {
	return DividendHistoryQuery{eventStream, 0, ""}
}

func (dividendHistoryQuery *DividendHistoryQuery) GetDividends() []Dividend {
	dividends := []Dividend{}

	for _, event := range dividendHistoryQuery.EventStream.Get() {
		if event.Name == dividend.DividendRecordedEventName {
			if !dividendHistoryQuery.dividendMatchesYearFilter(event.Payload["date"].(string)) {
				continue
			}
			ticker := event.Payload["ticker"].(string)
			if !dividendHistoryQuery.dividendMatchesTickerFilter(ticker) {
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

func (dividendHistoryQuery *DividendHistoryQuery) GetSum() float32 {
	dividends := float32(0.0)

	for _, event := range dividendHistoryQuery.EventStream.Get() {
		if event.Name == dividend.DividendRecordedEventName {
			if !dividendHistoryQuery.dividendMatchesYearFilter(event.Payload["date"].(string)) {
				continue
			}
			ticker := event.Payload["ticker"].(string)
			if !dividendHistoryQuery.dividendMatchesTickerFilter(ticker) {
				continue
			}
			dividends += getFloatValue(event.Payload["net"])

			continue
		}
	}

	return dividends
}

func (dividendHistoryQuery *DividendHistoryQuery) SetYearFilter(year int) {
	dividendHistoryQuery.yearFilter = year
}

func (dividendHistoryQuery *DividendHistoryQuery) SetTickerFilter(ticker string) {
	dividendHistoryQuery.tickerFilter = ticker
}

func getFloatValue(value interface{}) float32 {
	floatValue, ok := value.(float32)
	if !ok {
		floatValue = float32(value.(float64))
	}

	return floatValue
}

func (dividendHistoryQuery *DividendHistoryQuery) dividendMatchesYearFilter(date string) bool {
	dividendDate, _ := time.Parse("2006-01-02", date)
	if dividendHistoryQuery.yearFilter != 0 && dividendDate.Year() != dividendHistoryQuery.yearFilter {
		return false
	}
	return true
}

func (dividendHistoryQuery *DividendHistoryQuery) dividendMatchesTickerFilter(ticker string) bool {
	if dividendHistoryQuery.tickerFilter != "" && ticker != dividendHistoryQuery.tickerFilter {
		return false
	}
	return true
}
