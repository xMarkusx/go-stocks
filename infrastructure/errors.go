package infrastructure

type UnsupportedDateFormatError struct {
	prob string
}

type InvalidDateError struct {
	prob string
}

func NewUnsupportedDateFormatError(prob string) *UnsupportedDateFormatError {
	return &UnsupportedDateFormatError{prob: prob}
}

func NewInvalidDateError(prob string) *InvalidDateError {
	return &InvalidDateError{prob: prob}
}

func (e *UnsupportedDateFormatError) Error() string {
	return e.prob
}

func (e *InvalidDateError) Error() string {
	return e.prob
}
