package positionList

import (
	"reflect"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"stock-monitor/query"
	"testing"
)

func TestPositionListProvidesCompleteListOfPositions(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 20},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	valueTracker := query.FakeValueTracker{map[string]float32{"MO": 10.00}}

	positionListQuery := PositionListQuery{&infrastructure.InMemoryEventStream{events}, valueTracker}
	got := positionListQuery.GetPositions()
	want := map[string]Position{"MO": {"MO", 25, 250.00}}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}

func TestPositionIsRemovedWhenCompletelySold(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20},
			map[string]interface{}{"occurred_at": "2001-01-02"},
		},
	}

	positionListQuery := PositionListQuery{&infrastructure.InMemoryEventStream{events}, query.FakeValueTracker{}}
	_, found := positionListQuery.GetPositions()["MO"]

	if found {
		t.Errorf("Expected no position of MO in portfolio but found one")
	}
}

func TestPositionListHandlesTickerRenames(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2002-01-01"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 20},
			map[string]interface{}{"occurred_at": "2002-01-02"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "MO", "new": "FOO"},
			map[string]interface{}{"occurred_at": "2002-01-03"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "FOO", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2002-01-04"},
		},
	}

	valueTracker := query.FakeValueTracker{map[string]float32{"FOO": 10.00}}

	positionListQuery := PositionListQuery{&infrastructure.InMemoryEventStream{events}, valueTracker}
	got := positionListQuery.GetPositions()
	want := map[string]Position{"FOO": {"FOO", 25, 250.00}}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}

func TestPositionListHandlesMultipleTickerRenames(t *testing.T) {
	events := []infrastructure.Event{
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2002-01-01"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 20},
			map[string]interface{}{"occurred_at": "2002-01-02"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "MO", "new": "FOO"},
			map[string]interface{}{"occurred_at": "2002-01-03"},
		},
		{
			portfolio.SharesAddedToPortfolioEventName,
			map[string]interface{}{"ticker": "FOO", "price": 20.45, "shares": 10},
			map[string]interface{}{"occurred_at": "2002-01-04"},
		},
		{
			portfolio.TickerRenamedEventName,
			map[string]interface{}{"old": "FOO", "new": "BAR"},
			map[string]interface{}{"occurred_at": "2002-01-05"},
		},
		{
			portfolio.SharesRemovedFromPortfolioEventName,
			map[string]interface{}{"ticker": "BAR", "price": 40.00, "shares": 5},
			map[string]interface{}{"occurred_at": "2002-01-06"},
		},
	}

	valueTracker := query.FakeValueTracker{map[string]float32{"BAR": 10.00}}

	positionListQuery := PositionListQuery{&infrastructure.InMemoryEventStream{events}, valueTracker}
	got := positionListQuery.GetPositions()
	want := map[string]Position{"BAR": {"BAR", 35, 350.00}}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Positions unequal got: %#v, want: %#v", got, want)
	}
}
