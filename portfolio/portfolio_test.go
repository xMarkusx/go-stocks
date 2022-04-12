package portfolio

import (
	"reflect"
	"testing"
	"time"
)

type fakePortfolioState struct {
	lastOrderDate string
	numberOfSharesForTicker map[string]int
	addOrders []addSharesToPortfolioCommand
	removeOrders []removeSharesFromPortfolioCommand
}

func (state *fakePortfolioState) GetNumberOfSharesForTicker(ticker string) int {
	return state.numberOfSharesForTicker[ticker]
}

func (state *fakePortfolioState) GetDateOfLastOrder() string {
	return state.lastOrderDate
}

func (state *fakePortfolioState) AddShares(command addSharesToPortfolioCommand) {
	state.addOrders = append(state.addOrders, command)
}

func (state *fakePortfolioState) RemoveShares(command removeSharesFromPortfolioCommand) {
	state.removeOrders = append(state.removeOrders, command)
}

func initFakeState() fakePortfolioState {
	return fakePortfolioState{"", map[string]int{}, []addSharesToPortfolioCommand{}, []removeSharesFromPortfolioCommand{}}
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
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = "2000-01-02"

	portfolio.AddSharesToPortfolio(command)

	expected := []addSharesToPortfolioCommand{{"MO", 10, 20.45, "2000-01-02"}}
	got := state.addOrders
	
	if reflect.DeepEqual(got, expected) == false {
		t.Errorf("State not updated. Expected:%#v Got:%#v", expected, got)
	}
}

func TestCanRemoveShares(t *testing.T) {
	state := initFakeState()
	state.numberOfSharesForTicker = map[string]int{"MO": 11}
	portfolio := Portfolio{&state}
	command := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	command.Date = "2000-01-02"

	portfolio.RemoveSharesFromPortfolio(command)

	expected := []removeSharesFromPortfolioCommand{{"MO", 10, 20.45, "2000-01-02"}}
	got := state.removeOrders
	
	if reflect.DeepEqual(got, expected) == false {
		t.Errorf("State not updated. Expected:%#v Got:%#v", expected, got)
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}
	command := AddSharesToPortfolioCommand("MO", 0, 20.45)

	err := portfolio.AddSharesToPortfolio(command)

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}
	command := AddSharesToPortfolioCommand("MO", -10, 20.45)

	err := portfolio.AddSharesToPortfolio(command)

	_, ok := err.(*InvalidNumbersOfSharesError)
	if !ok {
		t.Errorf("Expected InvalidNumbersOfSharesError but got %#v", err)
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}

	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 21, 20.45)

	err := portfolio.RemoveSharesFromPortfolio(removeSharesCommand)

	_, ok := err.(*CantSellMoreSharesThanExistingError)
	if !ok {
		t.Errorf("Expected CantSellMoreSharesThanExistingError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenAddingShares(t *testing.T) {
	state := initFakeState()
	portfolio := Portfolio{&state}
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = "Foo"

	err := portfolio.AddSharesToPortfolio(command)

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateHasToBeInValidFormatWhenRemovingShares(t *testing.T) {
	state := initFakeState()
	state.numberOfSharesForTicker = map[string]int{"MO": 20}
	portfolio := Portfolio{&state}
	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand.Date = "Foo"

	err := portfolio.RemoveSharesFromPortfolio(removeSharesCommand)

	_, ok := err.(*UnsupportedDateFormatError)
	if !ok {
		t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
	}
}

func TestDateCanNotBeInTheFutureWhenAddingShares(t *testing.T) {
	today := time.Now()
	state := initFakeState()
	portfolio := Portfolio{&state}
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = today.AddDate(0,0,1).Format("2006-01-02")

	err := portfolio.AddSharesToPortfolio(command)

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
	removeSharesCommand := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	removeSharesCommand.Date = today.AddDate(0,0,1).Format("2006-01-02")

	err := portfolio.RemoveSharesFromPortfolio(removeSharesCommand)

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}

func TestDateCanNotBeOlderThanDateOfLastOrderWhenAddingShares(t *testing.T) {
	state := initFakeState()
	state.lastOrderDate = "2020-01-02"
	portfolio := Portfolio{&state}
	command := AddSharesToPortfolioCommand("MO", 10, 20.45)
	command.Date = "2020-01-01"

	err := portfolio.AddSharesToPortfolio(command)

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
	command := RemoveSharesFromPortfolioCommand("MO", 10, 20.45)
	command.Date = "2020-01-01"

	err := portfolio.RemoveSharesFromPortfolio(command)

	_, ok := err.(*InvalidDateError)
	if !ok {
		t.Errorf("Expected InvalidDateError but got %#v", err)
	}
}
