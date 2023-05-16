package di

import (
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
