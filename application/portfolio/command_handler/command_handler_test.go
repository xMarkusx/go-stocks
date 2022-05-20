package command_handler

import (
	"reflect"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/event"
	"stock-monitor/application/portfolio/persistence"
	"stock-monitor/domain/portfolio"
	"stock-monitor/infrastructure"
	"testing"
)

func TestItHandlesAddSharesToPortfolioCommand(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	addSharesCommand := command.NewAddSharesToPortfolioCommand("MO", 10, 9.99)
	addSharesCommand.Date = "2000-01-01"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := NewCommandHandler(&repository, publisher)

	commandHandler.HandleAddSharesToPortfolio(addSharesCommand)

	expectedEvent := infrastructure.Event{
		portfolio.SharesAddedToPortfolioEvent,
		map[string]interface{}{
			"ticker": "MO",
			"shares": 10,
			"price":  float32(9.99),
			"date":   "2000-01-01",
		},
	}
	got := eventStream.Events[0]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestItReturnsDomainErrorWhenAddSharesToPortfolioCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	addSharesCommand := command.NewAddSharesToPortfolioCommand("MO", 0, 9.99)
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleAddSharesToPortfolio(addSharesCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestSharesRemovedFromPortfolioEventCanBePublished(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{portfolio.SharesAddedToPortfolioEvent, map[string]interface{}{
				"ticker": "MO",
				"shares": 20,
				"price":  10.00,
				"date":   "2000-01-01",
			}},
		},
	}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 10, 9.99)
	removeSharesCommand.Date = "2000-01-02"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := NewCommandHandler(&repository, publisher)

	commandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)

	expectedEvent := infrastructure.Event{
		portfolio.SharesRemovedFromPortfolioEvent,
		map[string]interface{}{
			"ticker": "MO",
			"shares": 10,
			"price":  float32(9.99),
			"date":   "2000-01-02",
		},
	}
	got := eventStream.Events[1]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestItReturnsDomainErrorWhenRemoveSharesFromPortfolioCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 10, 9.99)
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}
