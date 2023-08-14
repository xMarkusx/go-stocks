package sell_stock

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/command_handler"
	"stock-monitor/application/shared"
)

type SellStockHandler struct {
	CommandHandler command_handler.PortfolioCommandHandlerInterface
}

type SellOrder struct {
	Ticker string  `json:"ticker"`
	Shares int     `json:"shares"`
	Price  float32 `json:"price"`
	Date   string  `json:"date"`
}

func (handler *SellStockHandler) SellStock(c echo.Context) error {
	sellOrder := new(SellOrder)
	if err := c.Bind(sellOrder); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand(sellOrder.Ticker, sellOrder.Shares, sellOrder.Price, shared.CommandDate(sellOrder.Date))

	err := handler.CommandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)

	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}
