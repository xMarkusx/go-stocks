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
	aliases        []string
	numberOfShares int
	price          float32
	date           string
}

func (order *Order) Dto() (string, string, []string, int, float32, string) {
	return order.orderType, order.ticker, order.aliases, order.numberOfShares, order.price, order.date
}

type OrderHistoryQuery struct {
	EventStream infrastructure.EventStream
}

func (orderHistoryQuery *OrderHistoryQuery) GetOrders() []Order {
	orders := []Order{}
	renames := []infrastructure.Event{}
	for _, event := range orderHistoryQuery.EventStream.Get() {
		if event.Name == portfolio.SharesAddedToPortfolioEventName {
			ticker, shares, price, date := extractEventData(event)
			order := Order{"BUY", ticker, []string{}, shares, price, date}
			orders = append(orders, order)

			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEventName {
			ticker, shares, price, date := extractEventData(event)
			order := Order{"SELL", ticker, []string{}, shares, price, date}
			orders = append(orders, order)

			continue
		}

		if event.Name == portfolio.TickerRenamedEventName {
			renames = append(renames, event)
			continue
		}
	}

	for _, renameEvent := range renames {
		oldSymbol := renameEvent.Payload["old"].(string)
		newSymbol := renameEvent.Payload["new"].(string)
		for key, order := range orders {
			if order.ticker == oldSymbol {
				orders[key].ticker = newSymbol
				orders[key].aliases = append(orders[key].aliases, oldSymbol)
			}
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
	date, ok := event.MetaData["occurred_at"].(string)

	return ticker, shares, price, date
}