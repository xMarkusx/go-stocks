package main

import (
	"github.com/labstack/echo/v4"
	"stock-monitor/infrastructure/di"
	"stock-monitor/infrastructure/handler/add_stock"
	"stock-monitor/infrastructure/handler/rename_stock"
	"stock-monitor/infrastructure/handler/sell_stock"
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

	commandHandler := di.MakePortfolioCommandHandler()
	addStockHandler := add_stock.AddStockHandler{commandHandler}
	e.POST("/add-stock", addStockHandler.AddStock)

	sellStockHandler := sell_stock.SellStockHandler{commandHandler}
	e.POST("/sell-stock", sellStockHandler.SellStock)

	renameStockHandler := rename_stock.RenameStockHandler{commandHandler}
	e.POST("/rename-stock", renameStockHandler.RenameStock)

	e.Logger.Fatal(e.Start(":8080"))
}
