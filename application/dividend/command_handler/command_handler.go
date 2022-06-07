package command_handler

import (
	"stock-monitor/application/dividend/command"
	"stock-monitor/application/dividend/persistence"
	"stock-monitor/application/event"
)

type CommandHandler struct {
	repository persistence.DividendRepository
	publisher  event.EventPublisher
}

func NewCommandHandler(repository persistence.DividendRepository, publisher event.EventPublisher) CommandHandler {
	return CommandHandler{repository: repository, publisher: publisher}
}

func (commandHandler *CommandHandler) HandleRecordDividend(command command.RecordDividendCommand) error {
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
