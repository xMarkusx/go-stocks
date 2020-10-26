package main

import (
	"fmt"
	"os"
    "encoding/gob"
)

type OrderStorage interface {
    Add(order Order)
    Get() []Order
}

type InMemoryOrderStorage struct {
    Orders []Order
}

func (orderStorage *InMemoryOrderStorage) Add(order Order) {
    orderStorage.Orders = append(orderStorage.Orders, order)
}

func (orderStorage *InMemoryOrderStorage) Get() []Order {
    return orderStorage.Orders
}

type FileSystemOrderStorage struct {
	storagePath string
	fileName string
}

func (orderStorage *FileSystemOrderStorage) Add(order Order) {
	orders := append(orderStorage.Get(), order)
	err := writeGob(orderStorage.storagePath + orderStorage.fileName, orders)
	if err != nil{
			fmt.Println(err)
	}
}

func (orderStorage *FileSystemOrderStorage) Get() []Order {
    orders := []Order{}
	readGob(orderStorage.storagePath + orderStorage.fileName, &orders)

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
