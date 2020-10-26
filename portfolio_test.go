package main

import "testing"
import "reflect"

func TestBuyOrder(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	p.addBuyOrder("MO", 20.45, 20)
	got := p.getPositions()["MO"]
	want := Position{"MO", 20}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %q, want: %q", got.toString(), want.toString())
	}
}

func TestMultipleBuyOrderForSameTicker(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	p.addBuyOrder("MO", 20.45, 20)
	p.addBuyOrder("MO", 30.45, 20)
	got := p.getPositions()["MO"]
	want := Position{"MO", 40}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %q, want: %q", got.toString(), want.toString())
	}
}

func TestCanNotBuyZeroShares(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	err := p.addBuyOrder("MO", 20.45, 0)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestCanNotBuyNegativeNumberOfShares(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	err := p.addBuyOrder("MO", 20.45, -10)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestSellOrder(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	p.addBuyOrder("MO", 20.45, 20)
	p.addSellOrder("MO", 20.45, 10)
	got := p.getPositions()["MO"]
	want := Position{"MO", 10}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %q, want: %q", got.toString(), want.toString())
	}
}

func TestMultipleSellOrdersOnSamePosition(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	p.addBuyOrder("MO", 20.45, 20)
	p.addSellOrder("MO", 20.45, 10)
	p.addSellOrder("MO", 20.45, 5)
	got := p.getPositions()["MO"]
	want := Position{"MO", 5}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %q, want: %q", got.toString(), want.toString())
	}
}

func TestCanNotSellMoreSharesThenCurrentlyInPortfolio(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	p.addBuyOrder("MO", 20.45, 20)
	err := p.addSellOrder("MO", 20.45, 21)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func TestPositionIsRemovedWhenCompletelySold(t *testing.T) {
	p := initPortfolio(&InMemoryOrderStorage{})
	p.addBuyOrder("MO", 20.45, 20)
	p.addSellOrder("MO", 20.45, 20)
	_, found := p.getPositions()["MO"]

	if found {
		t.Errorf("Expected no position of MO in portfolio but found one")
	}
}

func TestPortfolioCanBeInitializedWithOrders(t *testing.T) {
	orders := []Order{
		Order{BuyOrderType, "MO", 20.45, 10},
		Order{BuyOrderType, "PG", 40.00, 20},
		Order{SellOrderType, "MO", 24.00, 5},
	}
	p := initPortfolio(&InMemoryOrderStorage{orders})
	got := p.getPositions()
	mo := got["MO"]
	pg := got["PG"]

	expectedNumberOfPositions := 2
	expectedMo := Position{"MO", 5}
	expectedPg := Position{"PG", 20}

	if len(got) != expectedNumberOfPositions {
		t.Errorf("Unexpected number of positions in portfolio. Expected %d, got %d", expectedNumberOfPositions, len(got))
	}
	if reflect.DeepEqual(mo, expectedMo) == false {
		t.Errorf("Positions unequal got: %q, want: %q", mo.toString(), expectedMo.toString())
	}
	if reflect.DeepEqual(pg, expectedPg) == false {
		t.Errorf("Positions unequal got: %q, want: %q", pg.toString(), expectedPg.toString())
	}
}

func TestOrdersWillBeAddedToStorage(t *testing.T) {
	os := &InMemoryOrderStorage{}
	p := initPortfolio(os)
	p.addBuyOrder("MO", 20.45, 20)
	p.addSellOrder("MO", 20.45, 10)
	p.addSellOrder("MO", 20.45, 5)

	got := os.Get()
	want := []Order{
		Order{BuyOrderType, "MO", 20.45, 20},
		Order{SellOrderType, "MO", 20.45, 10},
		Order{SellOrderType, "MO", 20.45, 5},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Order storage unequal")
	}
}
