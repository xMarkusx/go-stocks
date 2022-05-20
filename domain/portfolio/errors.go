package portfolio

type InvalidNumbersOfSharesError struct {
	prob string
}

type CantSellMoreSharesThanExistingError struct {
	prob string
}

type UnsupportedDateFormatError struct {
	prob string
}

type InvalidDateError struct {
	prob string
}

func (e *InvalidNumbersOfSharesError) Error() string {
	return e.prob
}

func (e *CantSellMoreSharesThanExistingError) Error() string {
	return e.prob
}

func (e *UnsupportedDateFormatError) Error() string {
	return e.prob
}

func (e *InvalidDateError) Error() string {
	return e.prob
}
