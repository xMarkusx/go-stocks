package add_stock

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/command_handler"
	"stock-monitor/application/shared"
)

type AddStockHandler struct {
	CommandHandler command_handler.PortfolioCommandHandlerInterface
}

type BuyOrder struct {
	Ticker string  `json:"ticker"`
	Shares int     `json:"shares"`
	Price  float32 `json:"price"`
	Date   string  `json:"date"`
}

func (handler *AddStockHandler) AddStock(c echo.Context) error {
	buyOrder := new(BuyOrder)
	if err := c.Bind(buyOrder); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	addSharesCommand := command.NewAddSharesToPortfolioCommand(buyOrder.Ticker, buyOrder.Shares, buyOrder.Price, shared.CommandDate(buyOrder.Date))

	err := handler.CommandHandler.HandleAddSharesToPortfolio(addSharesCommand)

	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}
