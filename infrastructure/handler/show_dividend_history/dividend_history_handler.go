package show_dividend_history

import (
	"github.com/labstack/echo/v4"
	"net/http"
	dividend_history "stock-monitor/query/dividend-history"
	"strconv"
)

type ShowDividendHistoryHandler struct {
	Query dividend_history.DividendHistoryQueryInterface
}

type DividendResponse struct {
	Ticker string
	Net    float32
	Gross  float32
	Date   string
}

type DividendHistoryResponse struct {
	Dividends      []DividendResponse
	TotalDividends float32
}

func (handler *ShowDividendHistoryHandler) ShowDividendHistory(c echo.Context) error {
	dividends := []DividendResponse{}

	filter := dividend_history.NewFilter()
	yearParam := c.QueryParam("year")
	if yearParam != "" {
		year, err := strconv.Atoi(yearParam)
		if err != nil {
			return err
		}
		filter.ByYear(year)
	}

	ticker := c.QueryParam("ticker")
	if ticker != "" {
		filter.ByTicker(ticker)
	}

	for _, dividend := range handler.Query.GetDividends(filter) {
		dividends = append(dividends, DividendResponse{
			Ticker: dividend.Ticker,
			Net:    dividend.Net,
			Gross:  dividend.Gross,
			Date:   dividend.Date,
		})
	}

	dividendHistoryResponse := DividendHistoryResponse{
		Dividends:      dividends,
		TotalDividends: handler.Query.GetSum(filter),
	}

	return c.JSON(http.StatusOK, dividendHistoryResponse)
}
