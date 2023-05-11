package importer

import (
	"strconv"
)

type ImportItem struct {
	Type   string
	Date   string
	Ticker string
	Alias  string
	Price  float32
	Shares int
}

func Parse(csv [][]string) []ImportItem {
	importItems := []ImportItem{}
	for _, record := range csv {
		shares, _ := strconv.Atoi(record[5])
		price, _ := strconv.ParseFloat(record[4], 32)
		pricef32 := float32(price)

		item := ImportItem{
			record[0],
			record[1],
			record[2],
			record[3],
			pricef32,
			shares,
		}

		importItems = append(importItems, item)
	}

	return importItems
}
