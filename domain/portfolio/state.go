package portfolio

type PortfolioState struct {
	positions     map[string]Position
	lastOrderDate string
}

type Position struct {
	Ticker string
	Shares int
}

func NewPortfolioState() PortfolioState {
	return PortfolioState{map[string]Position{}, ""}
}

func (portfolioState *PortfolioState) GetNumberOfSharesForTicker(ticker string) int {
	return portfolioState.positions[ticker].Shares
}

func (portfolioState *PortfolioState) GetDateOfLastOrder() string {
	return portfolioState.lastOrderDate
}

func (portfolioState *PortfolioState) AddShares(ticker string, shares int, date string) {
	portfolioState.lastOrderDate = date

	p, found := portfolioState.positions[ticker]
	if !found {
		portfolioState.positions[ticker] = Position{ticker, shares}
	} else {
		p.Shares += shares
		portfolioState.positions[ticker] = p
	}
}

func (portfolioState *PortfolioState) RemoveShares(ticker string, shares int, date string) {
	portfolioState.lastOrderDate = date

	p := portfolioState.positions[ticker]

	p.Shares -= shares
	portfolioState.positions[ticker] = p
}
