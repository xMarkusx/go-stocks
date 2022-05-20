package positionList

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"stock-monitor/query"
)

type Position struct {
	Ticker       string
	Shares       int
	CurrentValue float32
}

type PositionListQuery struct {
	EventStream  infrastructure.EventStream
	ValueTracker query.ValueTracker
}

func (positionListQuery *PositionListQuery) GetPositions() map[string]Position {
	positions := map[string]Position{}
	for ticker, shares := range runPositionListProjection(positionListQuery.EventStream) {
		currentValue := positionListQuery.ValueTracker.Current(ticker) * float32(shares)
		positions[ticker] = Position{ticker, shares, currentValue}
	}
	return positions
}

func runPositionListProjection(eventStream infrastructure.EventStream) map[string]int {
	positions := map[string]int{}

	for _, event := range eventStream.Get() {
		ticker, shares, _ := extractEventData(event)

		if event.Name == portfolio.SharesAddedToPortfolioEvent {
			currentShares, found := positions[ticker]
			if !found {
				positions[ticker] = shares
			} else {
				positions[ticker] = currentShares + shares
			}

			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEvent {
			currentShares := positions[ticker]

			if currentShares == shares {
				delete(positions, ticker)
				continue
			}

			positions[ticker] = currentShares - shares
			continue
		}
	}

	return positions
}

func extractEventData(event infrastructure.Event) (string, int, float32) {
	ticker := event.Payload["ticker"].(string)
	shares := event.Payload["shares"].(int)
	price, ok := event.Payload["price"].(float32)
	if !ok {
		price = float32(event.Payload["price"].(float64))
	}

	return ticker, shares, price
}
