package main

import "fmt"
import "errors"

type Position struct {
	Ticker string
	Shares int
}

type Order struct {
	OrderType orderType
	Ticker    string
	Price     float32
	Shares    int
}

type orderType string

const (
	BuyOrderType  orderType = "BUY"
	SellOrderType orderType = "SELL"
)

func (pos Position) toString() string {
	return fmt.Sprintf("Ticker: %q, shares: %d", pos.Ticker, pos.Shares)
}

type Portfolio struct {
	Positions    map[string]Position
	OrderStorage OrderStorage
}

func initPortfolio(orderStorage OrderStorage) Portfolio {
	p := Portfolio{map[string]Position{}, orderStorage}
	for _, order := range p.OrderStorage.Get() {
		p.apply(order)
	}
	return p
}

func (portfolio *Portfolio) addBuyOrder(ticker string, price float32, shares int) error {
	if shares <= 0 {
		return errors.New("number of shares must be greater than 0")
	}

	o := Order{BuyOrderType, ticker, price, shares}

	portfolio.apply(o)

	portfolio.OrderStorage.Add(o)

	return nil
}

func (portfolio *Portfolio) addSellOrder(ticker string, price float32, shares int) error {
	position := portfolio.Positions[ticker]
	if position.Shares < shares {
		return errors.New("not allowed to sell more shares than currently in portfolio")
	}

	o := Order{SellOrderType, ticker, price, shares}

	portfolio.apply(o)

	portfolio.OrderStorage.Add(o)

	return nil
}

func (portfolio *Portfolio) getPositions() map[string]Position {
	return portfolio.Positions
}

func (portfolio *Portfolio) apply(order Order) {

	if order.OrderType == BuyOrderType {
		position, found := portfolio.Positions[order.Ticker]
		if !found {
			portfolio.Positions[order.Ticker] = Position{order.Ticker, order.Shares}
		} else {
			position.Shares += order.Shares
			portfolio.Positions[order.Ticker] = position
		}

		return
	}

	if order.OrderType == SellOrderType {
		position := portfolio.Positions[order.Ticker]

		if position.Shares == order.Shares {
			delete(portfolio.Positions, order.Ticker)
			return
		}

		position.Shares -= order.Shares
		portfolio.Positions[order.Ticker] = position
		return
	}
}
