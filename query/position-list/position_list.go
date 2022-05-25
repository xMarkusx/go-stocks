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
		if event.Name == portfolio.SharesAddedToPortfolioEventName {
			ticker := event.Payload["ticker"].(string)
			shares := event.Payload["shares"].(int)
			currentShares, found := positions[ticker]
			if !found {
				positions[ticker] = shares
			} else {
				positions[ticker] = currentShares + shares
			}

			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEventName {
			ticker := event.Payload["ticker"].(string)
			shares := event.Payload["shares"].(int)
			currentShares := positions[ticker]

			if currentShares == shares {
				delete(positions, ticker)
				continue
			}

			positions[ticker] = currentShares - shares
			continue
		}

		if event.Name == portfolio.TickerRenamedEventName {
			oldSymbol := event.Payload["old"].(string)
			newSymbol := event.Payload["new"].(string)

			currentShares := positions[oldSymbol]
			delete(positions, oldSymbol)

			positions[newSymbol] = currentShares

			continue
		}
	}

	return positions
}
