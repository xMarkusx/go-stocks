package totalInvestedMoney

import (
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
)

type TotalInvestedMoneyQuery struct {
	EventStream infrastructure.EventStream
}

func (totalInvestedMoneyQuery *TotalInvestedMoneyQuery) GetTotalInvestedMoney() float32 {
	invested := float32(0.0)
	for _, event := range totalInvestedMoneyQuery.EventStream.Get() {

		if event.Name == portfolio.SharesAddedToPortfolioEventName {
			shares, price := extractEventData(event)
			invested += price * float32(shares)
			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEventName {
			shares, price := extractEventData(event)
			invested -= price * float32(shares)
			continue
		}
	}

	return invested
}

func extractEventData(event infrastructure.Event) (int, float32) {
	shares := event.Payload["shares"].(int)
	price, ok := event.Payload["price"].(float32)
	if !ok {
		price = float32(event.Payload["price"].(float64))
	}

	return shares, price
}
