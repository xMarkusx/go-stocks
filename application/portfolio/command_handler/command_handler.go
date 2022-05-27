package command_handler

import (
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/event"
	"stock-monitor/application/portfolio/persistence"
)

type CommandHandler struct {
	repository persistence.PortfolioRepository
	publisher  event.PortfolioEventPublisher
}

func NewCommandHandler(repository persistence.PortfolioRepository, publisher event.PortfolioEventPublisher) CommandHandler {
	return CommandHandler{repository: repository, publisher: publisher}
}

func (commandHandler *CommandHandler) HandleAddSharesToPortfolio(command command.AddSharesToPortfolioCommand) error {
	p := commandHandler.repository.Load()

	err := p.AddSharesToPortfolio(command.Ticker, command.NumberOfShares, command.Price)

	if err != nil {
		return err
	}

	err = commandHandler.publisher.PublishDomainEvents(p.GetRecordedEvents(), command.Date)

	if err != nil {
		return err
	}

	return nil
}

func (commandHandler *CommandHandler) HandleRemoveSharesFromPortfolio(command command.RemoveSharesFromPortfolioCommand) error {
	p := commandHandler.repository.Load()

	err := p.RemoveSharesFromPortfolio(command.Ticker, command.NumberOfShares, command.Price)

	if err != nil {
		return err
	}

	err = commandHandler.publisher.PublishDomainEvents(p.GetRecordedEvents(), command.Date)

	if err != nil {
		return err
	}

	return nil
}

func (commandHandler *CommandHandler) HandleRenameTicker(command command.RenameTickerCommand) error {
	p := commandHandler.repository.Load()

	err := p.RenameTicker(command.Old, command.New)

	if err != nil {
		return err
	}

	err = commandHandler.publisher.PublishDomainEvents(p.GetRecordedEvents(), command.Date)

	if err != nil {
		return err
	}

	return nil
}
