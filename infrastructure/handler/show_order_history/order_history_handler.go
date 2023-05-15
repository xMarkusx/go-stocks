package show_order_history

import (
	"github.com/labstack/echo/v4"
	"net/http"
	orderHistory "stock-monitor/query/order-history"
)

type ShowOrderHistoryHandler struct {
	Query orderHistory.OrderHistoryQueryInterface
}

type OrderResponse struct {
	OrderType      string
	Ticker         string
	Aliases        []string
	NumberOfShares int
	Price          float32
	Date           string
}

func (handler *ShowOrderHistoryHandler) ShowOrderHistory(c echo.Context) error {
	orderResponse := []OrderResponse{}

	for _, order := range handler.Query.GetOrders() {
		orderResponse = append(orderResponse, OrderResponse{
			OrderType:      order.OrderType,
			Ticker:         order.Ticker,
			Aliases:        order.Aliases,
			NumberOfShares: order.NumberOfShares,
			Price:          order.Price,
			Date:           order.Date,
		})
	}

	return c.JSON(http.StatusOK, orderResponse)
}
