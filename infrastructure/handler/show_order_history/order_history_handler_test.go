package show_order_history_test

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"reflect"
	"stock-monitor/infrastructure/handler/show_order_history"
	orderHistory "stock-monitor/query/order-history"
	"testing"
)

type MockOrderHistoryQuery struct {
	orders []orderHistory.Order
}

func (mockOrderHistory *MockOrderHistoryQuery) GetOrders() []orderHistory.Order {
	return mockOrderHistory.orders
}

func TestShowOrderHistory(t *testing.T) {
	t.Run("should return 200 status ok", func(t *testing.T) {
		mock := MockOrderHistoryQuery{[]orderHistory.Order{{
			"BUY",
			"MO",
			[]string{"FOO"},
			10,
			10.00,
			"2001-01-01",
		}}}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := show_order_history.ShowOrderHistoryHandler{&mock}
		handler.ShowOrderHistory(c)

		if reflect.DeepEqual(rec.Code, http.StatusOK) == false {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusOK, rec.Code)
		}
	})
}
