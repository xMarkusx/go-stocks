package portfolio_test

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"reflect"
	showPortfolioHandler "stock-monitor/infrastructure/handler/portfolio"
	positionList "stock-monitor/query/position-list"
	"testing"
)

type MockPositionList struct {
	positions map[string]positionList.Position
}

func (mockPositionList *MockPositionList) GetPositions() map[string]positionList.Position {
	return mockPositionList.positions
}

func TestShowPortfolio(t *testing.T) {
	t.Run("should return 200 status ok", func(t *testing.T) {
		mock := MockPositionList{positions: map[string]positionList.Position{
			"MO": {
				Ticker:       "MO",
				Shares:       10,
				CurrentValue: 10,
			},
		}}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := showPortfolioHandler.ShowPortfolioHandler{&mock}
		handler.ShowPortfolio(c)

		if reflect.DeepEqual(rec.Code, http.StatusOK) == false {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusOK, rec.Code)
		}
	})
}
