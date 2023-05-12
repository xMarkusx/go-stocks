package di

import (
	"stock-monitor/infrastructure"
	"stock-monitor/query"
	positionList "stock-monitor/query/position-list"
)

func MakePortfolioEventStream() infrastructure.EventStream {
	return &infrastructure.FileSystemEventStream{"./store/", "portfolio_event_stream.gob"}
}

func MakePositionListQuery() positionList.PositionListQuery {
	eventStream := MakePortfolioEventStream()
	return &positionList.EventStreamedPositionListQuery{eventStream, query.FinnHubValueTracker{}}
}
