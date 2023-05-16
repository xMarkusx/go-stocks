package main

import (
	"github.com/labstack/echo/v4"
	"stock-monitor/infrastructure/di"
	"stock-monitor/infrastructure/handler/show_dividend_history"
	"stock-monitor/infrastructure/handler/show_order_history"
	"stock-monitor/infrastructure/handler/show_portfolio"
)

func main() {
	e := echo.New()

	positionListQuery := di.MakePositionListQuery()
	positionListHandler := show_portfolio.ShowPortfolioHandler{positionListQuery}
	e.GET("/portfolio", positionListHandler.ShowPortfolio)

	orderHistoryQuery := di.MakeOrderHistoryQuery()
	orderHistoryHandler := show_order_history.ShowOrderHistoryHandler{orderHistoryQuery}
	e.GET("/order-history", orderHistoryHandler.ShowOrderHistory)

	dividendHistoryQuery := di.MakeDividendHistoryQuery()
	dividendHistoryHandler := show_dividend_history.ShowDividendHistoryHandler{dividendHistoryQuery}
	e.GET("/dividend-history", dividendHistoryHandler.ShowDividendHistory)

	e.Logger.Fatal(e.Start(":8080"))
}
