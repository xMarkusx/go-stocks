package command_handler_test

import (
	"reflect"
	"stock-monitor/application/dividend/command"
	"stock-monitor/application/dividend/command_handler"
	"stock-monitor/application/dividend/persistence"
	"stock-monitor/application/event"
	"stock-monitor/domain/dividend"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestItHandlesRecordDividendCommand(t *testing.T) {
	portfolioEventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
					"date":   "2000-01-01",
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
		},
	}
	dividendEventStream := infrastructure.InMemoryEventStream{}

	publisher := event.NewEventPublisher(&dividendEventStream)
	recordDividendCommand := command.NewRecordDividendCommand("MO", 20.00, 21.00, "2000-01-02")
	repository := persistence.NewEventSourcedDividendRepository(&portfolioEventStream)
	commandHandler := command_handler.NewDividendCommandHandler(&repository, publisher)

	commandHandler.HandleRecordDividend(recordDividendCommand)

	expectedEvent := infrastructure.Event{
		dividend.DividendRecordedEventName,
		map[string]interface{}{
			"ticker": "MO",
			"net":    float32(20.00),
			"gross":  float32(21.00),
			"date":   "2000-01-02",
		},
		map[string]interface{}{"occurred_at": "2000-01-02"},
	}
	got := dividendEventStream.Events[0]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestItReturnsErrorWhenRecordDividendCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewEventPublisher(&eventStream)
	recordDividendCommand := command.NewRecordDividendCommand("MO", 0, 9.99, "2000-01-02")
	repository := persistence.NewEventSourcedDividendRepository(&eventStream)
	commandHandler := command_handler.NewDividendCommandHandler(&repository, publisher)

	err := commandHandler.HandleRecordDividend(recordDividendCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestItReturnsErrorWhenPublishingEventAfterRecordDividendCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewEventPublisher(&eventStream)
	recordDividendCommand := command.NewRecordDividendCommand("MO", 20.00, 21.00, "FOO")
	repository := persistence.NewEventSourcedDividendRepository(&eventStream)
	commandHandler := command_handler.NewDividendCommandHandler(&repository, publisher)

	err := commandHandler.HandleRecordDividend(recordDividendCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}
