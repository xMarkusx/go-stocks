package portfolio

type InvalidNumbersOfSharesError struct {
	prob string
}

type CantSellMoreSharesThanExistingError struct {
	prob string
}

type TickerNotInPortfolioError struct {
	prob string
}

type TickerAlreadyUsedError struct {
	prob string
}

func (e *InvalidNumbersOfSharesError) Error() string {
	return e.prob
}

func (e *CantSellMoreSharesThanExistingError) Error() string {
	return e.prob
}

func (e *TickerNotInPortfolioError) Error() string {
	return e.prob
}

func (e *TickerAlreadyUsedError) Error() string {
	return e.prob
}
