package rename_stock

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/command_handler"
	"stock-monitor/application/shared"
)

type RenameStockHandler struct {
	CommandHandler command_handler.PortfolioCommandHandlerInterface
}

type RenameOrder struct {
	OldTicker string `json:"old_ticker"`
	NewTicker string `json:"new_ticker"`
	Date      string `json:"date"`
}

func (handler *RenameStockHandler) RenameStock(c echo.Context) error {
	renameOrder := new(RenameOrder)
	if err := c.Bind(renameOrder); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	renameTickerCommand := command.NewRenameTickerCommand(renameOrder.OldTicker, renameOrder.NewTicker, shared.CommandDate(renameOrder.Date))

	err := handler.CommandHandler.HandleRenameTicker(renameTickerCommand)

	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}
