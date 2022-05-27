package infrastructure

type UnsupportedDateFormatError struct {
	prob string
}

type InvalidDateError struct {
	prob string
}

func (e *UnsupportedDateFormatError) Error() string {
	return e.prob
}

func (e *InvalidDateError) Error() string {
	return e.prob
}
