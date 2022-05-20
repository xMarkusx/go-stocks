package event

import (
	"stock-monitor/application/portfolio/command"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
)

type PortfolioEventPublisher struct {
	eventStream infrastructure.EventStream
}

func NewPortfolioEventPublisher(eventStream infrastructure.EventStream) PortfolioEventPublisher {
	return PortfolioEventPublisher{eventStream: eventStream}
}

func (publisher *PortfolioEventPublisher) PublishSharesAddedToPortfolioEvent(command command.AddSharesToPortfolioCommand) {
	sharesAddedToPortfolioEvent := infrastructure.Event{
		portfolio.SharesAddedToPortfolioEvent,
		map[string]interface{}{
			"ticker": command.Ticker,
			"shares": command.NumberOfShares,
			"price":  command.Price,
			"date":   command.Date,
		},
	}

	publisher.eventStream.Add(sharesAddedToPortfolioEvent)
}
func (publisher *PortfolioEventPublisher) PublishSharesRemovedFromPortfolioEvent(command command.RemoveSharesFromPortfolioCommand) {
	sharesRemovedFromPortfolioEvent := infrastructure.Event{
		portfolio.SharesRemovedFromPortfolioEvent,
		map[string]interface{}{
			"ticker": command.Ticker,
			"shares": command.NumberOfShares,
			"price":  command.Price,
			"date":   command.Date,
		},
	}

	publisher.eventStream.Add(sharesRemovedFromPortfolioEvent)
}
