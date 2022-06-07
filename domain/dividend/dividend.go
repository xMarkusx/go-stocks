package dividend

import (
	"stock-monitor/domain"
	"stock-monitor/domain/portfolio"
	"time"
)

type Dividend struct {
	Positions map[string]string
	events    []domain.DomainEvent
}

func NewDividend() Dividend {
	return Dividend{map[string]string{}, []domain.DomainEvent{}}
}

func (d *Dividend) RecordDividend(ticker string, net float32, gross float32, date string) error {
	stockAddedDate, found := d.Positions[ticker]
	if !found {
		return NewTickerUnknownError(ticker)
	}
	if !dividendRecordedAfterStockWasAdded(date, stockAddedDate) {
		return NewDividendDateBeforeSharesWereAddedToPortfolioError(ticker, date)
	}
	if net <= 0 {
		return &DividendNetZeroOrNegativeError{}
	}
	if gross <= 0 {
		return &DividendGrossZeroOrNegativeError{}
	}

	dividendRecordedEvent := NewDividendRecordedEvent(ticker, net, gross, date)
	d.events = append(d.events, &dividendRecordedEvent)

	return nil
}

func (d *Dividend) GetRecordedEvents() []domain.DomainEvent {
	return d.events
}

func (d *Dividend) Apply(event domain.DomainEvent) {
	if event.Name() == portfolio.SharesAddedToPortfolioEventName {
		ticker := event.Payload()["ticker"].(string)
		date := event.Payload()["date"].(string)
		_, found := d.Positions[ticker]
		if found {
			return
		}
		d.Positions[ticker] = date
	}

	if event.Name() == portfolio.TickerRenamedEventName {
		oldTicker := event.Payload()["old"].(string)
		newTicker := event.Payload()["new"].(string)
		addedDate := d.Positions[oldTicker]
		d.Positions[newTicker] = addedDate
		delete(d.Positions, oldTicker)
	}
}

func dividendRecordedAfterStockWasAdded(dividendRecorded string, stockAdded string) bool {
	dividendDate, _ := time.Parse("2006-01-02", dividendRecorded)
	stockDate, _ := time.Parse("2006-01-02", stockAdded)
	diff := dividendDate.Sub(stockDate)
	if diff > 0 {
		return true
	}

	return false
}
