package positionList

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"stock-monitor/query"
)

type PositionListQuery interface {
	GetPositions() map[string]Position
}

type Position struct {
	Ticker       string
	Shares       int
	CurrentValue float32
}

type EventStreamedPositionListQuery struct {
	EventStream  infrastructure.EventStream
	ValueTracker query.ValueTracker
}

func (positionListQuery *EventStreamedPositionListQuery) GetPositions() map[string]Position {
	positions := map[string]Position{}
	positionChannel := make(chan Position)

	positionProjection := runPositionListProjection(positionListQuery.EventStream)

	for ticker, shares := range positionProjection {
		go func(ticker string, shares int) {
			currentValue := positionListQuery.ValueTracker.Current(ticker) * float32(shares)
			positionChannel <- Position{ticker, shares, currentValue}
			positions[ticker] = Position{ticker, shares, currentValue}
		}(ticker, shares)
	}

	for i := 0; i < len(positionProjection); i++ {
		position := <-positionChannel
		positions[position.Ticker] = position
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
