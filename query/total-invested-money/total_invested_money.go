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
		_, shares, price := extractEventData(event)

		if event.Name == portfolio.SharesAddedToPortfolioEvent {
			invested += price * float32(shares)
			continue
		}

		if event.Name == portfolio.SharesRemovedFromPortfolioEvent {
			invested -= price * float32(shares)
			continue
		}
	}

	return invested
}

func extractEventData(event infrastructure.Event) (string, int, float32) {
	ticker := event.Payload["ticker"].(string)
	shares := event.Payload["shares"].(int)
	price, ok := event.Payload["price"].(float32)
	if !ok {
		price = float32(event.Payload["price"].(float64))
	}

	return ticker, shares, price
}
