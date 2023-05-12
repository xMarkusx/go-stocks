package main

import (
	"stock-monitor/infrastructure/di"
	"stock-monitor/infrastructure/handler/portfolio"

	"github.com/labstack/echo/v4"
)

func main() {
	positionListQuery := di.MakePositionListQuery()
	positionListHandler := portfolio.ShowPortfolioHandler{positionListQuery}

	e := echo.New()

	e.GET("/portfolio", positionListHandler.ShowPortfolio)

	e.Logger.Fatal(e.Start(":8080"))
}
