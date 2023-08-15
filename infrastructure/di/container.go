package di

import (
	command_handler2 "stock-monitor/application/dividend/command_handler"
	persistence2 "stock-monitor/application/dividend/persistence"
	"stock-monitor/application/event"
	"stock-monitor/application/portfolio/command_handler"
	"stock-monitor/application/portfolio/persistence"
	"stock-monitor/infrastructure"
	"stock-monitor/query"
	dividend_history "stock-monitor/query/dividend-history"
	orderHistory "stock-monitor/query/order-history"
	positionList "stock-monitor/query/position_list"
)

func MakePortfolioEventStream() infrastructure.EventStream {
	return &infrastructure.FileSystemEventStream{"./store/", "portfolio_event_stream.gob"}
}

func MakeDividendEventStream() infrastructure.EventStream {
	return &infrastructure.FileSystemEventStream{"./store/", "dividend_event_stream.gob"}
}

func MakePositionListQuery() positionList.PositionListQuery {
	eventStream := MakePortfolioEventStream()
	return &positionList.EventStreamedPositionListQuery{eventStream, query.FinnHubValueTracker{}}
}

func MakeOrderHistoryQuery() orderHistory.OrderHistoryQueryInterface {
	eventStream := MakePortfolioEventStream()
	return &orderHistory.OrderHistoryQuery{eventStream}
}

func MakeDividendHistoryQuery() dividend_history.DividendHistoryQueryInterface {
	eventStream := MakeDividendEventStream()
	dividendQuery := dividend_history.NewDividendHistoryQuery(eventStream)
	return &dividendQuery
}

func MakePortfolioCommandHandler() command_handler.PortfolioCommandHandlerInterface {
	eventStream := MakePortfolioEventStream()
	publisher := event.NewEventPublisher(eventStream)
	repository := persistence.NewEventSourcedPortfolioRepository(eventStream)
	return command_handler.NewCommandHandler(&repository, publisher)
}

func MakeDividendCommandHandler() command_handler2.DividendCommandHandlerInterface {
	dividendEventStream := MakeDividendEventStream()
	portfolioEventStream := MakePortfolioEventStream()
	publisher := event.NewEventPublisher(dividendEventStream)
	repository := persistence2.NewEventSourcedDividendRepository(portfolioEventStream)
	return command_handler2.NewDividendCommandHandler(&repository, publisher)
}
