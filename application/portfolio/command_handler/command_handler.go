package command_handler

import (
	"stock-monitor/application/event"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/persistence"
)

type PortfolioCommandHandlerInterface interface {
	HandleAddSharesToPortfolio(command command.AddSharesToPortfolioCommand) error
	HandleRemoveSharesFromPortfolio(command command.RemoveSharesFromPortfolioCommand) error
	HandleRenameTicker(command command.RenameTickerCommand) error
}

type CommandHandler struct {
	repository persistence.PortfolioRepository
	publisher  event.EventPublisher
}

func NewCommandHandler(repository persistence.PortfolioRepository, publisher event.EventPublisher) PortfolioCommandHandlerInterface {
	return &CommandHandler{repository: repository, publisher: publisher}
}

func (commandHandler *CommandHandler) HandleAddSharesToPortfolio(command command.AddSharesToPortfolioCommand) error {
	p := commandHandler.repository.Load()

	err := p.AddSharesToPortfolio(command.Ticker, command.NumberOfShares, command.Price, command.Date)

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
