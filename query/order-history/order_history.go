package orderHistory

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
)

type OrderHistoryQueryInterface interface {
	GetOrders() []Order
}

type Order struct {
	orderType      string
	ticker         string
	numberOfShares int
	price          float32
	date           string
}

func (order *Order) Dto() (string, string, int, float32, string) {
	return order.orderType, order.ticker, order.numberOfShares, order.price, order.date
}

type OrderHistoryQuery struct {
	EventStream infrastructure.EventStream
}

func (orderHistoryQuery *OrderHistoryQuery) GetOrders() []Order {
	orders := []Order{}
	for _, event := range orderHistoryQuery.EventStream.Get() {
		if event.Name == portfolio.SharesAddedToPortfolioEvent {
			ticker, shares, price, date := extractEventData(event)
			order := Order{"BUY", ticker, shares, price, date}
			orders = append(orders, order)

			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEvent {
			ticker, shares, price, date := extractEventData(event)
			order := Order{"SELL", ticker, shares, price, date}
			orders = append(orders, order)

			continue
		}
	}

	return orders
}

func extractEventData(event infrastructure.Event) (string, int, float32, string) {
	ticker := event.Payload["ticker"].(string)
	shares := event.Payload["shares"].(int)
	price, ok := event.Payload["price"].(float32)
	if !ok {
		price = float32(event.Payload["price"].(float64))
	}
	date, ok := event.Payload["date"].(string)

	return ticker, shares, price, date
}
