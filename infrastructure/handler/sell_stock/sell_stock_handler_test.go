package sell_stock_test

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/infrastructure/handler/sell_stock"
	"strings"
	"testing"
)

type mockPortfolioCommandHandler struct {
	removeSharesCommand command.RemoveSharesFromPortfolioCommand
	expectedError       error
}

func (mockPortfolioCommandHandler *mockPortfolioCommandHandler) HandleAddSharesToPortfolio(command command.AddSharesToPortfolioCommand) error {
	return mockPortfolioCommandHandler.expectedError
}

func (mockPortfolioCommandHandler *mockPortfolioCommandHandler) HandleRemoveSharesFromPortfolio(command command.RemoveSharesFromPortfolioCommand) error {
	mockPortfolioCommandHandler.removeSharesCommand = command
	return mockPortfolioCommandHandler.expectedError
}

func (mockPortfolioCommandHandler *mockPortfolioCommandHandler) HandleRenameTicker(command command.RenameTickerCommand) error {
	return mockPortfolioCommandHandler.expectedError
}

func (mockPortfolioCommandHandler *mockPortfolioCommandHandler) expectError(err error) {
	mockPortfolioCommandHandler.expectedError = err
}

func TestSellStock(t *testing.T) {
	t.Run("it succeeds", func(t *testing.T) {
		mock := mockPortfolioCommandHandler{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"foo\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := sell_stock.SellStockHandler{&mock}
		handler.SellStock(c)

		if rec.Code != http.StatusCreated {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusCreated, rec.Code)
		}
	})

	t.Run("it fails with 422 when command failed", func(t *testing.T) {
		mock := mockPortfolioCommandHandler{}
		mock.expectError(errors.New("some error happened"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"foo\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := sell_stock.SellStockHandler{&mock}
		handler.SellStock(c)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("it fails with 400 when input is not accepted", func(t *testing.T) {
		mock := mockPortfolioCommandHandler{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"ticker\":1234}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := sell_stock.SellStockHandler{&mock}
		handler.SellStock(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusBadRequest, rec.Code)
		}
	})
}
