package positionList

import (
	"reflect"
	"stock-monitor/query"
	"stock-monitor/infrastructure"
	"testing"
)

func TestPositionListProvidesCompleteListOfPositions(t *testing.T) {
	events := []infrastructure.Event{
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 10}},
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 40.00, "shares": 5}},
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
		{"Portfolio.SharesAddedToPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20}},
		{"Portfolio.SharesRemovedFromPortfolio", map[string]interface{}{"ticker": "MO", "price": 20.45, "shares": 20}},
	}

	positionListQuery := PositionListQuery{&infrastructure.InMemoryEventStream{events}, query.FakeValueTracker{}}
	_, found := positionListQuery.GetPositions()["MO"]

	if found {
		t.Errorf("Expected no position of MO in portfolio but found one")
	}
}
