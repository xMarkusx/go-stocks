package portfolio

import (
	"reflect"
	"testing"
	"time"
)

type fakePortfolioState struct {
	lastOrderDate           string
	numberOfSharesForTicker map[string]int
	addOrders               []map[string]interface{}
	removeOrders            []map[string]interface{}
}

func (state *fakePortfolioState) GetNumberOfSharesForTicker(ticker string) int {
	return state.numberOfSharesForTicker[ticker]
}

func (state *fakePortfolioState) GetDateOfLastOrder() string {
	return state.lastOrderDate
}

func (state *fakePortfolioState) AddShares(ticker string, shares int, date string) {
	state.addOrders = append(state.addOrders, map[string]interface{}{
		"ticker": ticker,
		"shares": shares,
		"date":   date,
	})
}

func (state *fakePortfolioState) RemoveShares(ticker string, shares int, date string) {
	state.removeOrders = append(state.addOrders, map[string]interface{}{
		"ticker": ticker,
		"shares": shares,
		"date":   date,
	})
}

func initFakeState() fakePortfolioState {
	return fakePortfolioState{"", map[string]int{}, []map[string]interface{}{}, []map[string]interface{}{}}
}

func TestCanCreatePortfolio(t *testing.T) {
	state := initFakeState()

	portfolio := NewPortfolio(&state)

	expected := Portfolio{&state}

	if reflect.DeepEqual(portfolio, expected) == false {
		t.Errorf("State not updated. Expected:%#v Got:%#v", expected, portfolio)
	}
}

func TestCanAddShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}

	portfolio.AddSharesToPortfolio("MO", 10, "2000-01-02")

	expected := []map[string]interface{}{
		{
			"ticker": "MO",
			"shares": 10,
			"date":   "2000-01-02",
		},
	}
	got := state.addOrders

	if reflect.DeepEqual(got, expected) == false {
		t.Errorf("State not updated. Expected:%#v Got:%#v", expected, got)
	}
}

func TestCanRemoveShares(t *testing.T) {
	state := initFakeState()
	state.numberOfSharesForTicker = map[string]int{"MO": 11}
	portfolio := Portfolio{&state}

	portfolio.RemoveSharesFromPortfolio("MO", 10, "2000-01-02")

	expected := []map[string]interface{}{
		{
			"ticker": "MO",
			"shares": 10,
			"date":   "2000-01-02",
		},
	}
	got := state.removeOrders

	if reflect.DeepEqual(got, expected) == false {
		t.Errorf("State not updated. Expected:%#v Got:%#v", expected, got)
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}

	err := portfolio.AddSharesToPortfolio("MO", 0, "2000-01-02")

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}

	err := portfolio.AddSharesToPortfolio("MO", -10, "2000-01-02")

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}

	err := portfolio.RemoveSharesFromPortfolio("MO", 21, "2000-01-02")

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenAddingShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}

	err := portfolio.AddSharesToPortfolio("MO", 10, "Foo")

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenRemovingShares(t *testing.T) {
	state := initFakeState()
	state.numberOfSharesForTicker = map[string]int{"MO": 20}
	portfolio := Portfolio{&state}

	err := portfolio.RemoveSharesFromPortfolio("MO", 10, "Foo")

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenAddingShares(t *testing.T) {
	today := time.Now()
	state := initFakeState()
	portfolio := Portfolio{&state}

	err := portfolio.AddSharesToPortfolio("MO", 10, today.AddDate(0, 0, 1).Format("2006-01-02"))

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenRemovingShares(t *testing.T) {
	today := time.Now()
	state := initFakeState()
	state.numberOfSharesForTicker = map[string]int{"MO": 20}
	portfolio := Portfolio{&state}

	err := portfolio.RemoveSharesFromPortfolio("MO", 10, today.AddDate(0, 0, 1).Format("2006-01-02"))

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenAddingShares(t *testing.T) {
	state := initFakeState()
	state.lastOrderDate = "2020-01-02"
	portfolio := Portfolio{&state}

	err := portfolio.AddSharesToPortfolio("MO", 10, "2020-01-01")

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenRemovingShares(t *testing.T) {
	state := initFakeState()
	state.lastOrderDate = "2020-01-02"
	state.numberOfSharesForTicker = map[string]int{"MO": 20}
	portfolio := Portfolio{&state}

	err := portfolio.RemoveSharesFromPortfolio("MO", 10, "2020-01-01")

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}
