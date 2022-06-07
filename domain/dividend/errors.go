package dividend

type TickerUnknownError struct {
	ticker string
}

type DividendDateBeforeSharesWereAddedToPortfolioError struct {
	ticker string
	added  string
}

type DividendNetZeroOrNegativeError struct{}

type DividendGrossZeroOrNegativeError struct{}

func NewTickerUnknownError(ticker string) *TickerUnknownError {
	return &TickerUnknownError{ticker: ticker}
}

func NewDividendDateBeforeSharesWereAddedToPortfolioError(ticker string, added string) *DividendDateBeforeSharesWereAddedToPortfolioError {
	return &DividendDateBeforeSharesWereAddedToPortfolioError{ticker: ticker, added: added}
}

func (e *TickerUnknownError) Error() string {
	return "ticker not added to portfolio. ticker: " + e.ticker
}

func (e *DividendDateBeforeSharesWereAddedToPortfolioError) Error() string {
	return "dividend date is before shares were added to portfolio. ticker: " + e.ticker + " date: " + e.added
}

func (e *DividendNetZeroOrNegativeError) Error() string {
	return "dividend net must be greater than zero"
}

func (e *DividendGrossZeroOrNegativeError) Error() string {
	return "dividend gross must be greater than zero"
}
