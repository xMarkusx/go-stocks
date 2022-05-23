package portfolio

const SharesAddedToPortfolioEventName = "Portfolio.SharesAddedToPortfolio"
const SharesRemovedFromPortfolioEventName = "Portfolio.SharesRemovedFromPortfolio"
const TickerRenamedEventName = "Portfolio.TickerRenamed"

type SharesAddedToPortfolioEvent struct {
	ticker string
	shares int
	price  float32
	date   string
}

func NewSharesAddedToPortfolioEvent(ticker string, shares int, price float32, date string) SharesAddedToPortfolioEvent {
	return SharesAddedToPortfolioEvent{ticker: ticker, shares: shares, price: price, date: date}
}

func (event *SharesAddedToPortfolioEvent) Name() string {
	return SharesAddedToPortfolioEventName
}

func (event *SharesAddedToPortfolioEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"ticker": event.ticker,
		"shares": event.shares,
		"price":  event.price,
		"date":   event.date,
	}
}

type SharesRemovedFromPortfolioEvent struct {
	ticker string
	shares int
	price  float32
	date   string
}

func NewSharesRemovedFromPortfolioEvent(ticker string, shares int, price float32, date string) SharesRemovedFromPortfolioEvent {
	return SharesRemovedFromPortfolioEvent{ticker: ticker, shares: shares, price: price, date: date}
}

func (event *SharesRemovedFromPortfolioEvent) Name() string {
	return SharesRemovedFromPortfolioEventName
}

func (event *SharesRemovedFromPortfolioEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"ticker": event.ticker,
		"shares": event.shares,
		"price":  event.price,
		"date":   event.date,
	}
}

type TickerRenamedEvent struct {
	old  string
	new  string
	date string
}

func NewTickerRenamedEvent(old string, new string, date string) TickerRenamedEvent {
	return TickerRenamedEvent{old: old, new: new, date: date}
}

func (event *TickerRenamedEvent) Name() string {
	return TickerRenamedEventName
}

func (event *TickerRenamedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{
		"old":  event.old,
		"new":  event.new,
		"date": event.date,
	}
}
