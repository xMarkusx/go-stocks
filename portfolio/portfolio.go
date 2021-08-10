package portfolio

import "errors"

type Order struct {
	orderType orderType
	ticker    string
	price     float32
	shares    int
}

type orderType string

const (
	BuyOrderType  orderType = "BUY"
	SellOrderType orderType = "SELL"
)

type Position struct {
	ticker string
	shares int
}

func (pos Position) Dto() (string, int) {
	return pos.ticker, pos.shares
}

func (pos Position) CurrentValue(valueTracker ValueTracker) float32 {
	return valueTracker.Current(pos.ticker) * float32(pos.shares)
}

type Portfolio struct {
	positions    map[string]Position
	orderStorage OrderStorage
}

func InitPortfolio(orderStorage OrderStorage) Portfolio {
	p := Portfolio{map[string]Position{}, orderStorage}
	for _, order := range p.orderStorage.Get() {
		p.applyOrder(order)
	}
	return p
}

func (portfolio *Portfolio) AddBuyOrder(ticker string, price float32, shares int) error {
	if shares <= 0 {
		return errors.New("number of shares must be greater than 0")
	}

	o := Order{BuyOrderType, ticker, price, shares}

	portfolio.applyOrder(o)

	portfolio.orderStorage.Add(o)

	return nil
}

func (portfolio *Portfolio) AddSellOrder(ticker string, price float32, shares int) error {
	position := portfolio.positions[ticker]
	if position.shares < shares {
		return errors.New("not allowed to sell more shares than currently in portfolio")
	}

	o := Order{SellOrderType, ticker, price, shares}

	portfolio.applyOrder(o)

	portfolio.orderStorage.Add(o)

	return nil
}

func (portfolio *Portfolio) GetPositions() map[string]Position {
	return portfolio.positions
}

func (portfolio *Portfolio) GetTotalInvestedMoney() float32 {
	invested := float32(0.0)
	for _, order := range portfolio.orderStorage.Get() {
		if order.orderType == BuyOrderType {
			invested += order.price * float32(order.shares)
			continue
		}

		if order.orderType == SellOrderType {
			invested -= order.price * float32(order.shares)
			continue
		}
	}

	return invested
}

func (portfolio *Portfolio) applyOrder(order Order) {

	if order.orderType == BuyOrderType {
		position, found := portfolio.positions[order.ticker]
		if !found {
			portfolio.positions[order.ticker] = Position{order.ticker, order.shares}
		} else {
			position.shares += order.shares
			portfolio.positions[order.ticker] = position
		}

		return
	}

	if order.orderType == SellOrderType {
		position := portfolio.positions[order.ticker]

		if position.shares == order.shares {
			delete(portfolio.positions, order.ticker)
			return
		}

		position.shares -= order.shares
		portfolio.positions[order.ticker] = position
		return
	}
}
