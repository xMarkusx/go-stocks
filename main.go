package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"stock-monitor/application/portfolio/command"
	"stock-monitor/application/portfolio/command_handler"
	"stock-monitor/application/portfolio/event"
	"stock-monitor/application/portfolio/persistence"
	"stock-monitor/infrastructure"
	"stock-monitor/query"
	"stock-monitor/query/order-history"
	"stock-monitor/query/position-list"
	"stock-monitor/query/total-invested-money"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	portfolioEventStream := &infrastructure.FileSystemEventStream{"./store/", "portfolio_event_stream.gob"}
	publisher := event.NewPortfolioEventPublisher(portfolioEventStream)
	repository := persistence.NewEventSourcedPortfolioRepository(portfolioEventStream)
	commandHandler := command_handler.NewCommandHandler(&repository, publisher)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "buy",
				Aliases: []string{},
				Usage:   "add a buy order",
				Action: func(c *cli.Context) error {
					ticker, price, shares, err := prepareOrderArgs(c.Args().Slice())
					addSharesCommand := command.NewAddSharesToPortfolioCommand(ticker, shares, price)
					date, dateErr := getDate(c.Args().Slice())
					if dateErr == nil {
						addSharesCommand.Date = date
					}

					if err == nil {
						err = commandHandler.HandleAddSharesToPortfolio(addSharesCommand)
					}

					if err != nil {
						fmt.Println(err.Error())

						return cli.Exit("Failed to add order", 1)
					}

					fmt.Println("added")
					return nil
				},
			},
			{
				Name:    "sell",
				Aliases: []string{},
				Usage:   "add a sell order",
				Action: func(c *cli.Context) error {
					ticker, price, shares, err := prepareOrderArgs(c.Args().Slice())
					removeSharesCommand := command.NewRemoveSharesFromPortfolioCommand(ticker, shares, price)
					date, dateErr := getDate(c.Args().Slice())
					if dateErr == nil {
						removeSharesCommand.Date = date
					}

					if err == nil {
						err = commandHandler.HandleRemoveSharesFromPortfolio(removeSharesCommand)
					}

					if err != nil {
						fmt.Println(err.Error())

						return cli.Exit("Failed to add order", 1)
					}

					fmt.Println("sold")
					return nil
				},
			},
			{
				Name:    "rename",
				Aliases: []string{},
				Usage:   "rename a ticker in portfolio",
				Action: func(c *cli.Context) error {
					oldSymbol := c.Args().Slice()[0]
					newSymbol := c.Args().Slice()[1]
					renameTickerCommand := command.NewRenameTickerCommand(oldSymbol, newSymbol)
					date, dateErr := getDate(c.Args().Slice())
					if dateErr == nil {
						renameTickerCommand.Date = date
					}

					err := commandHandler.HandleRenameTicker(renameTickerCommand)

					if err != nil {
						fmt.Println(err.Error())

						return cli.Exit("Failed to rename ticker", 1)
					}

					fmt.Println("ticker renamed")
					return nil
				},
			},
			{
				Name:    "show",
				Aliases: []string{"s"},
				Usage:   "show positions in portfolio",
				Action: func(c *cli.Context) error {
					positionListQuery := positionList.PositionListQuery{portfolioEventStream, query.FinnHubValueTracker{}}
					totalInvestedMoneyQuery := totalInvestedMoney.TotalInvestedMoneyQuery{portfolioEventStream}
					fmt.Printf("Total invested: %v \n", totalInvestedMoneyQuery.GetTotalInvestedMoney())
					for _, position := range positionListQuery.GetPositions() {
						fmt.Printf("Ticker: %q, shares: %d, value: %#v\n", position.Ticker, position.Shares, position.CurrentValue)
					}

					return nil
				},
			},
			{
				Name:    "order-history",
				Aliases: []string{"oh"},
				Usage:   "Shows history of all orders",
				Action: func(c *cli.Context) error {
					orderHistoryQuery := orderHistory.OrderHistoryQuery{portfolioEventStream}
					fmt.Print("Order history: \n")
					for _, order := range orderHistoryQuery.GetOrders() {
						fmt.Printf("%#v | %#v - Ticker: %#v, Aliases: %s, shares: %#v, price: %#v\n", order.Date, order.OrderType, order.Ticker, strings.Join(order.Aliases, ", "), order.NumberOfShares, order.Price)
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func prepareOrderArgs(args []string) (string, float32, int, error) {
	ticker := args[0]
	price, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return "", 0.0, 0, fmt.Errorf("input for price must be float32")
	}
	shares, err := strconv.Atoi(args[2])
	if err != nil {
		return "", 0.0, 0, fmt.Errorf("input for shares must be int")
	}

	return ticker, float32(price), shares, nil
}

func getDate(args []string) (string, error) {
	if len(args) < 4 {
		return "", errors.New("No date given")
	}

	return args[3], nil
}
