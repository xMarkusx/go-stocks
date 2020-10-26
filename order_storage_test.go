package main

import (
	"testing"
	"reflect"
	"os"
)

func TestInMemoryOrderStorage(t *testing.T) {
	os := InMemoryOrderStorage{}
	os.Add(Order{BuyOrderType, "MO", 20.45, 20})
	os.Add(Order{SellOrderType, "MO", 20.45, 10})
	os.Add(Order{SellOrderType, "MO", 20.45, 5})

    got := os.Get()
    want := []Order{
		Order{BuyOrderType, "MO", 20.45, 20},
		Order{SellOrderType, "MO", 20.45, 10},
		Order{SellOrderType, "MO", 20.45, 5},
	}

    if reflect.DeepEqual(got, want) == false {
        t.Errorf("Order storage unequal. Expected:%#v Got:%#v", want, got)
    }
}

func TestFileSystemOrderStorage(t *testing.T) {
	tmpStorePath := "./tmp/"
	tmpStoreFile := "test_orders.gob"

	os.Mkdir(tmpStorePath, 0777)
	orderStorage := FileSystemOrderStorage{tmpStorePath, tmpStoreFile}
	orderStorage.Add(Order{BuyOrderType, "MO", 20.45, 20})
	orderStorage.Add(Order{SellOrderType, "MO", 20.45, 10})
	orderStorage.Add(Order{SellOrderType, "MO", 20.45, 5})

    got := orderStorage.Get()
    want := []Order{
		Order{BuyOrderType, "MO", 20.45, 20},
		Order{SellOrderType, "MO", 20.45, 10},
		Order{SellOrderType, "MO", 20.45, 5},
	}

    if reflect.DeepEqual(got, want) == false {
        t.Errorf("Order storage unequal. Expected:%#v Got:%#v", want, got)
	}

	t.Cleanup(func() {
		os.Remove(tmpStorePath + tmpStoreFile)
		os.Remove(tmpStorePath)
	})
}
