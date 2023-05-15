package di

import (
	"stock-monitor/infrastructure"
	"stock-monitor/query"
	orderHistory "stock-monitor/query/order-history"
	positionList "stock-monitor/query/position_list"
)

func MakePortfolioEventStream() infrastructure.EventStream {
	return &infrastructure.FileSystemEventStream{"./store/", "portfolio_event_stream.gob"}
}

func MakePositionListQuery() positionList.PositionListQuery {
	eventStream := MakePortfolioEventStream()
	return &positionList.EventStreamedPositionListQuery{eventStream, query.FinnHubValueTracker{}}
}

func MakeOrderHistoryQuery() orderHistory.OrderHistoryQueryInterface {
	eventStream := MakePortfolioEventStream()
	return &orderHistory.OrderHistoryQuery{eventStream}
}
