package add_dividends

import (
	"github.com/labstack/echo/v4"
	"net/http"
	dividend_command "stock-monitor/application/dividend/command"
	dividend_command_handler "stock-monitor/application/dividend/command_handler"
	"stock-monitor/application/shared"
)

type AddDividendsHandler struct {
	CommandHandler dividend_command_handler.DividendCommandHandlerInterface
}

type Dividend struct {
	Ticker string  `json:"ticker"`
	Net    float32 `json:"net"`
	Gross  float32 `json:"gross"`
	Date   string  `json:"date"`
}

type Dividends struct {
	Dividends []Dividend `json:"dividends"`
}

func (handler *AddDividendsHandler) AddDividends(c echo.Context) error {
	dividends := new(Dividends)

	if err := c.Bind(dividends); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	for _, dividend := range dividends.Dividends {
		recordDividendCommand := dividend_command.NewRecordDividendCommand(dividend.Ticker, dividend.Net, dividend.Gross, shared.CommandDate(dividend.Date))

		err := handler.CommandHandler.HandleRecordDividend(recordDividendCommand)

		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
	}

	return c.NoContent(http.StatusCreated)
}
