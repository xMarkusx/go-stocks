package portfolio

import (
	"encoding/gob"
	"fmt"
	"os"
)

type OrderStorageDTO struct {
	OrderType orderType
	Ticker    string
	Price     float32
	Shares    int
}

type OrderStorage interface {
	Add(order Order)
	Get() []Order
}

type InMemoryOrderStorage struct {
	Orders []OrderStorageDTO
}

func (orderStorage *InMemoryOrderStorage) Add(order Order) {
	orderStorage.Orders = append(orderStorage.Orders, OrderStorageDTO{
		order.orderType,
		order.ticker,
		order.price,
		order.shares,
	})
}

func (orderStorage *InMemoryOrderStorage) Get() []Order {
	orders := []Order{}
	for _, orderDTO := range orderStorage.Orders {
		orders = append(orders, Order{
			orderDTO.OrderType,
			orderDTO.Ticker,
			orderDTO.Price,
			orderDTO.Shares,
		})
	}
	return orders
}

type FileSystemOrderStorage struct {
	StoragePath string
	FileName    string
}

func (orderStorage *FileSystemOrderStorage) Add(order Order) {
	orders := []OrderStorageDTO{}
	readGob(orderStorage.StoragePath+orderStorage.FileName, &orders)
	orders = append(orders, OrderStorageDTO{
		order.orderType,
		order.ticker,
		order.price,
		order.shares,
	})
	err := writeGob(orderStorage.StoragePath+orderStorage.FileName, orders)
	if err != nil {
		fmt.Println(err)
	}
}

func (orderStorage *FileSystemOrderStorage) Get() []Order {
	orders := []Order{}
	storedOrders := []OrderStorageDTO{}
	readGob(orderStorage.StoragePath+orderStorage.FileName, &storedOrders)

	for _, orderDTO := range storedOrders {
		orders = append(orders, Order{
			orderDTO.OrderType,
			orderDTO.Ticker,
			orderDTO.Price,
			orderDTO.Shares,
		})
	}

	return orders
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
