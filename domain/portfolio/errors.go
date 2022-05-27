package portfolio

type InvalidNumbersOfSharesError struct{}

type CantSellMoreSharesThanExistingError struct{}

type TickerNotInPortfolioError struct {
	ticker string
}

type TickerAlreadyUsedError struct {
	ticker string
}

func NewTickerNotInPortfolioError(ticker string) *TickerNotInPortfolioError {
	return &TickerNotInPortfolioError{ticker: ticker}
}

func NewTickerAlreadyUsedError(ticker string) *TickerAlreadyUsedError {
	return &TickerAlreadyUsedError{ticker: ticker}
}

func (e *InvalidNumbersOfSharesError) Error() string {
	return "number of shares must be greater than 0"
}

func (e *CantSellMoreSharesThanExistingError) Error() string {
	return "not allowed to sell more shares than currently in portfolio"
}

func (e *TickerNotInPortfolioError) Error() string {
	return "Ticker to be renamed not found. Ticker: " + e.ticker
}

func (e *TickerAlreadyUsedError) Error() string {
	return "New ticker symbol already in use. Ticker: " + e.ticker
}
