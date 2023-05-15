package show_portfolio

import (
	"github.com/labstack/echo/v4"
	"net/http"
	positionList "stock-monitor/query/position_list"
)

type ShowPortfolioHandler struct {
	Query positionList.PositionListQuery
}

type PositionResponse struct {
	Ticker       string
	Shares       int
	CurrentValue float32
}

func (handler *ShowPortfolioHandler) ShowPortfolio(c echo.Context) error {
	positionsResponse := map[string]PositionResponse{}

	for _, position := range handler.Query.GetPositions() {
		positionsResponse[position.Ticker] = PositionResponse{position.Ticker, position.Shares, position.CurrentValue}
	}

	return c.JSON(http.StatusOK, positionsResponse)
}
