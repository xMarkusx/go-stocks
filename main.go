package main

import(
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {
	orderStorage := &FileSystemOrderStorage{"./store/", "orders.gob"}
	p := initPortfolio(orderStorage)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "buy",
				Aliases: []string{},
				Usage: "add a buy order",
				Action: func(c *cli.Context) error {
					ticker, price, shares, error := prepareOrderArgs(c.Args().Slice())

					if error == nil {
						error = p.addBuyOrder(ticker, price, shares)
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
				Name: "sell",
				Aliases: []string{},
				Usage: "add a sell order",
				Action: func(c *cli.Context) error {
					ticker, price, shares, error := prepareOrderArgs(c.Args().Slice())

					if error == nil {
						error = p.addSellOrder(ticker, price, shares)
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
				Name: "show",
				Aliases: []string{"s"},
				Usage: "show positions in portfolio",
				Action: func(c *cli.Context) error {
					positions := p.getPositions()
					fmt.Printf("Number of positions: %d \n", len(positions))
					for _, position := range positions {
						fmt.Println(position.toString())
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
