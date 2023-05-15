package main

import (
	"github.com/labstack/echo/v4"
	"stock-monitor/infrastructure/di"
	"stock-monitor/infrastructure/handler/show_portfolio"
)

func main() {
	positionListQuery := di.MakePositionListQuery()
	positionListHandler := show_portfolio.ShowPortfolioHandler{positionListQuery}

	e := echo.New()

	e.GET("/portfolio", positionListHandler.ShowPortfolio)

	e.Logger.Fatal(e.Start(":8080"))
}
