package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"stock-monitor/infrastructure"
	"stock-monitor/portfolio"
	"stock-monitor/query"
	"stock-monitor/query/order-history"
	"stock-monitor/query/position-list"
	"stock-monitor/query/total-invested-money"

	"github.com/urfave/cli/v2"
)

func main() {
	portfolioEventStream := &infrastructure.FileSystemEventStream{"./store/", "portfolio_event_stream.gob"}
	state := portfolio.NewEventBasedPortfolioState(portfolioEventStream)
	p := portfolio.NewPortfolio(&state)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "buy",
				Aliases: []string{},
				Usage:   "add a buy order",
				Action: func(c *cli.Context) error {
					ticker, price, shares, error := prepareOrderArgs(c.Args().Slice())
					command := portfolio.AddSharesToPortfolioCommand(ticker, shares, price)
					date, dateErr := getDate(c.Args().Slice())
					if dateErr == nil {
						command.Date = date
					}

					if error == nil {
						error = p.AddSharesToPortfolio(command)
					}

					if error != nil {
						fmt.Println(error.Error())

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
					ticker, price, shares, error := prepareOrderArgs(c.Args().Slice())
					command := portfolio.RemoveSharesFromPortfolioCommand(ticker, shares, price)
					date, dateErr := getDate(c.Args().Slice())
					if dateErr == nil {
						command.Date = date
					}

					if error == nil {
						error = p.RemoveSharesFromPortfolio(command)
					}

					if error != nil {
						fmt.Println(error.Error())

						return cli.Exit("Failed to add order", 1)
					}

					fmt.Println("sold")
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
					fmt.Print("Order hsitory: \n")
					for _, order := range orderHistoryQuery.GetOrders() {
						orderType, ticker, shares, price, date := order.Dto()
						fmt.Printf("%#v | %#v - Ticker: %#v, shares: %#v, price: %#v\n", date, orderType, ticker, shares, price)
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
