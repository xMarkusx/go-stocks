package importer

import (
	"strconv"
)

type ImportItem struct {
	Date   string
	Ticker string
	Net    float32
	Gross  float32
}

func Parse(csv [][]string) []ImportItem {
	importItems := []ImportItem{}
	for _, record := range csv {
		net, _ := strconv.ParseFloat(record[2], 32)
		net32 := float32(net)
		gross, _ := strconv.ParseFloat(record[3], 32)
		gross32 := float32(gross)

		item := ImportItem{
			record[0],
			record[1],
			net32,
			gross32,
		}

		importItems = append(importItems, item)
	}

	return importItems
}
