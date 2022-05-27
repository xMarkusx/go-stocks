package command_handler_test

import (
	"reflect"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/command_handler"
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
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	commandHandler.HandleAddSharesToPortfolio(addSharesCommand)

	expectedEvent := infrastructure.Event{
		portfolio.SharesAddedToPortfolioEventName,
		map[string]interface{}{
			"ticker": "MO",
			"shares": 10,
			"price":  float32(9.99),
		},
		map[string]interface{}{"occurred_at": "2000-01-01"},
	}
	got := eventStream.Events[0]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestItReturnsErrorWhenAddSharesToPortfolioCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	addSharesCommand := command.NewAddSharesToPortfolioCommand("MO", 0, 9.99)
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleAddSharesToPortfolio(addSharesCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestItReturnsErrorWhenPublishingEventAfterAddSharesCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	addSharesCommand := command.NewAddSharesToPortfolioCommand("MO", 1, 9.99)
	addSharesCommand.Date = "FOO"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleAddSharesToPortfolio(addSharesCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestSharesRemovedFromPortfolioEventCanBePublished(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
		},
	}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 10, 9.99)
	removeSharesCommand.Date = "2000-01-02"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	commandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)

	expectedEvent := infrastructure.Event{
		portfolio.SharesRemovedFromPortfolioEventName,
		map[string]interface{}{
			"ticker": "MO",
			"shares": 10,
			"price":  float32(9.99),
		},
		map[string]interface{}{"occurred_at": "2000-01-02"},
	}
	got := eventStream.Events[1]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestItReturnsErrorWhenPublishingEventAfterRemoveSharesCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
		},
	}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 10, 9.99)
	removeSharesCommand.Date = "FOO"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestItReturnsErrorWhenRemoveSharesFromPortfolioCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand("MO", 10, 9.99)
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestRenameTickerCommandIsHandled(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
		},
	}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	renameTickerCommand := command.NewRenameTickerCommand("MO", "FOO")
	renameTickerCommand.Date = "2000-01-02"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	commandHandler.HandleRenameTicker(renameTickerCommand)

	expectedEvent := infrastructure.Event{
		portfolio.TickerRenamedEventName,
		map[string]interface{}{
			"old": "MO",
			"new": "FOO",
		},
		map[string]interface{}{"occurred_at": "2000-01-02"},
	}
	got := eventStream.Events[1]

	if reflect.DeepEqual(got, expectedEvent) == false {
		t.Errorf("Unexpected event published. Expected:%#v Got:%#v", expectedEvent, got)
	}
}

func TestItReturnsErrorWhenRenameTickerCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
		},
	}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	renameTickerCommand := command.NewRenameTickerCommand("FOO", "BAR")
	renameTickerCommand.Date = "2000-01-02"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleRenameTicker(renameTickerCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}

func TestItReturnsErrorWhenPublishingEventAfterRenameTickerCommandFails(t *testing.T) {
	eventStream := infrastructure.InMemoryEventStream{
		[]infrastructure.Event{
			{
				portfolio.SharesAddedToPortfolioEventName,
				map[string]interface{}{
					"ticker": "MO",
					"shares": 20,
					"price":  10.00,
				},
				map[string]interface{}{"occurred_at": "2000-01-01"},
			},
		},
	}
	publisher := event.NewPortfolioEventPublisher(&eventStream)
	renameTickerCommand := command.NewRenameTickerCommand("MO", "BAR")
	renameTickerCommand.Date = "FOO"
	repository := persistence.NewEventSourcedPortfolioRepository(&eventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	err := commandHandler.HandleRenameTicker(renameTickerCommand)

	if err == nil {
		t.Errorf("Expected Error but got none")
	}
}
