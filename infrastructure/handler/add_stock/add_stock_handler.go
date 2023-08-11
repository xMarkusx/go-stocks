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

type BuyOrders struct {
	Orders []BuyOrder `json:"orders"`
}

func (handler *AddStockHandler) AddStock(c echo.Context) error {
	buyOrders := new(BuyOrders)
	if err := c.Bind(buyOrders); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if len(buyOrders.Orders) == 0 {
		return c.String(http.StatusBadRequest, "no orders provided")
	}

	for _, buyOrder := range buyOrders.Orders {

		addSharesCommand := command.NewAddSharesToPortfolioCommand(buyOrder.Ticker, buyOrder.Shares, buyOrder.Price, shared.CommandDate(buyOrder.Date))

		err := handler.CommandHandler.HandleAddSharesToPortfolio(addSharesCommand)

		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
	}

	return c.NoContent(http.StatusCreated)
}
