package portfolio

const SharesAddedToPortfolioEventName = "Portfolio.SharesAddedToPortfolio"
const SharesRemovedFromPortfolioEventName = "Portfolio.SharesRemovedFromPortfolio"
const TickerRenamedEventName = "Portfolio.TickerRenamed"

type SharesAddedToPortfolioEvent struct {
	ticker string
	shares int
	price  float32
}

func NewSharesAddedToPortfolioEvent(ticker string, shares int, price float32) SharesAddedToPortfolioEvent {
	return SharesAddedToPortfolioEvent{ticker: ticker, shares: shares, price: price}
}

func (event *SharesAddedToPortfolioEvent) Name() string {
	return SharesAddedToPortfolioEventName
}

func (event *SharesAddedToPortfolioEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"ticker": event.ticker,
		"shares": event.shares,
		"price":  event.price,
	}
}

type SharesRemovedFromPortfolioEvent struct {
	ticker string
	shares int
	price  float32
}

func NewSharesRemovedFromPortfolioEvent(ticker string, shares int, price float32) SharesRemovedFromPortfolioEvent {
	return SharesRemovedFromPortfolioEvent{ticker: ticker, shares: shares, price: price}
}

func (event *SharesRemovedFromPortfolioEvent) Name() string {
	return SharesRemovedFromPortfolioEventName
}

func (event *SharesRemovedFromPortfolioEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"ticker": event.ticker,
		"shares": event.shares,
		"price":  event.price,
	}
}

type TickerRenamedEvent struct {
	old string
	new string
}

func NewTickerRenamedEvent(old string, new string) TickerRenamedEvent {
	return TickerRenamedEvent{old: old, new: new}
}

func (event *TickerRenamedEvent) Name() string {
	return TickerRenamedEventName
}

func (event *TickerRenamedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"old": event.old,
		"new": event.new,
	}
}
