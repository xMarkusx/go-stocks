package infrastructure

import (
	"encoding/gob"
	"fmt"
	"os"
)

type Event struct {
	Name    string
	Payload map[string]interface{}
}

type EventStream interface {
	Add(event Event)
	Get() []Event
}

type InMemoryEventStream struct {
	Events []Event
}

func (eventStream *InMemoryEventStream) Add(event Event) {
	eventStream.Events = append(eventStream.Events, event)
}

func (eventStream *InMemoryEventStream) Get() []Event {
	return eventStream.Events
}

type FileSystemEventStream struct {
	StoragePath string
	FileName    string
}

func (eventStream *FileSystemEventStream) Add(event Event) {
	events := []Event{}
	read(eventStream.StoragePath+eventStream.FileName, &events)
	events = append(events, event)

	err := write(eventStream.StoragePath+eventStream.FileName, events)
	if err != nil {
		fmt.Println(err)
	}
}

func (eventStream *FileSystemEventStream) Get() []Event {
	storedEvents := []Event{}
	read(eventStream.StoragePath+eventStream.FileName, &storedEvents)

	return storedEvents
}

func write(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func read(filePath string, object interface{}) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
