package add_dividends_test

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"stock-monitor/application/dividend/command"
	"stock-monitor/infrastructure/handler/add_dividends"
	"strings"
	"testing"
)

type mockDividendCommandHandler struct {
	recordDividendCommand command.RecordDividendCommand
	expectedError         error
}

func (mockDividendCommandHandler *mockDividendCommandHandler) HandleRecordDividend(command command.RecordDividendCommand) error {
	mockDividendCommandHandler.recordDividendCommand = command
	return mockDividendCommandHandler.expectedError
}

func (mockDividendCommandHandler *mockDividendCommandHandler) expectError(err error) {
	mockDividendCommandHandler.expectedError = err
}

func TestAddDividend(t *testing.T) {
	t.Run("it succeeds", func(t *testing.T) {
		mock := mockDividendCommandHandler{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"foo\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := add_dividends.AddDividendsHandler{&mock}
		handler.AddDividends(c)

		if rec.Code != http.StatusCreated {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusCreated, rec.Code)
		}
	})

	t.Run("it fails with 422 when command failed", func(t *testing.T) {
		mock := mockDividendCommandHandler{}
		mock.expectError(errors.New("some error happened"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"dividends\":[{\"ticker\":\"FOO\"}]}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := add_dividends.AddDividendsHandler{&mock}
		handler.AddDividends(c)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("it fails with 400 when input is not accepted", func(t *testing.T) {
		mock := mockDividendCommandHandler{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"dividends\":1234}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := add_dividends.AddDividendsHandler{&mock}
		handler.AddDividends(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Unexpected status code. Expected:%#v Got:%#v", http.StatusBadRequest, rec.Code)
		}
	})
}
