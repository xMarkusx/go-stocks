package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"stock-monitor/portfolio"
)

func main() {
	orderStorage := &portfolio.FileSystemEventStream{"./store/", "portfolio_event_stream.gob"}
	p := portfolio.InitPortfolio(orderStorage)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "buy",
				Aliases: []string{},
				Usage:   "add a buy order",
				Action: func(c *cli.Context) error {
					ticker, price, shares, error := prepareOrderArgs(c.Args().Slice())

					if error == nil {
						error = p.AddBuyOrder(ticker, price, shares)
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

					if error == nil {
						error = p.AddSellOrder(ticker, price, shares)
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
					positions := p.GetPositions()
					fmt.Printf("Number of positions: %d \n", len(positions))
					fmt.Printf("Total invested: %v \n", p.GetTotalInvestedMoney())
					for _, position := range positions {
						ticker, shares := position.Dto()
						valueTracker := portfolio.FinnHubValueTracker{}
						fmt.Printf("Ticker: %q, shares: %d, value: %#v\n", ticker, shares, position.CurrentValue(valueTracker))
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
