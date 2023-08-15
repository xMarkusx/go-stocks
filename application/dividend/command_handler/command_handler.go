package command_handler

import (
	"stock-monitor/application/dividend/command"
	"stock-monitor/application/dividend/persistence"
	"stock-monitor/application/event"
)

type DividendCommandHandlerInterface interface {
	HandleRecordDividend(command command.RecordDividendCommand) error
}

type DividendCommandHandler struct {
	repository persistence.DividendRepository
	publisher  event.EventPublisher
}

func NewDividendCommandHandler(repository persistence.DividendRepository, publisher event.EventPublisher) DividendCommandHandlerInterface {
	return &DividendCommandHandler{repository: repository, publisher: publisher}
}

func (commandHandler *DividendCommandHandler) HandleRecordDividend(command command.RecordDividendCommand) error {
	d := commandHandler.repository.Load()

	err := d.RecordDividend(command.Ticker, command.Net, command.Gross, command.Date)

	if err != nil {
		return err
	}

	err = commandHandler.publisher.PublishDomainEvents(d.GetRecordedEvents(), command.Date)

	if err != nil {
		return err
	}

	return nil
}
