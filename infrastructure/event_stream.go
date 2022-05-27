package infrastructure

import (
	"encoding/gob"
	"os"
	"time"
)

type Event struct {
	Name     string
	Payload  map[string]interface{}
	MetaData map[string]interface{}
}

type EventStream interface {
	Add(event Event) error
	Get() []Event
}

type InMemoryEventStream struct {
	Events []Event
}

func (eventStream *InMemoryEventStream) Add(event Event) error {
	occurredAt, ok := event.MetaData["occurred_at"].(string)
	if !ok || !commandDateHasValidFormat(occurredAt) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + occurredAt}
	}
	if !occurredAtIsInThePast(occurredAt) {
		return &InvalidDateError{"OccurredAt can't be in the future. Got: " + occurredAt}
	}

	if len(eventStream.Events) > 0 {
		lastOccurredAt := eventStream.Events[len(eventStream.Events)-1].MetaData["occurred_at"].(string)

		if !occurredAtIsLaterThanLastOccurredAt(occurredAt, lastOccurredAt) {
			return &InvalidDateError{"OccurredAt can't be older than occurredAt of last event. Got: " + occurredAt}
		}
	}

	eventStream.Events = append(eventStream.Events, event)

	return nil
}

func (eventStream *InMemoryEventStream) Get() []Event {
	return eventStream.Events
}

type FileSystemEventStream struct {
	StoragePath string
	FileName    string
}

func (eventStream *FileSystemEventStream) Add(event Event) error {
	occurredAt, ok := event.MetaData["occurred_at"].(string)
	if !ok || !commandDateHasValidFormat(occurredAt) {
		return &UnsupportedDateFormatError{"Unsupported date time format. Must be YYYY-MM-DD. Got: " + occurredAt}
	}
	if !occurredAtIsInThePast(occurredAt) {
		return &InvalidDateError{"OccurredAt can't be in the future. Got: " + occurredAt}
	}

	events := []Event{}
	read(eventStream.StoragePath+eventStream.FileName, &events)

	if len(events) > 0 {
		lastOccurredAt := events[len(events)-1].MetaData["occurred_at"].(string)

		if !occurredAtIsLaterThanLastOccurredAt(occurredAt, lastOccurredAt) {
			return &InvalidDateError{"OccurredAt can't be older than occurredAt of last event. Got: " + occurredAt}
		}
	}

	events = append(events, event)

	err := write(eventStream.StoragePath+eventStream.FileName, events)
	if err != nil {
		return err
	}

	return nil
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

func commandDateHasValidFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	return true
}

func occurredAtIsInThePast(occurredAt string) bool {
	occurredAtDate, _ := time.Parse("2006-01-02", occurredAt)
	today := time.Now()
	diff := today.Sub(occurredAtDate)

	if diff < 0 {
		return false
	}

	return true
}

func occurredAtIsLaterThanLastOccurredAt(occurredAt string, lastOccurredAt string) bool {
	oa, _ := time.Parse("2006-01-02", occurredAt)
	loa, _ := time.Parse("2006-01-02", lastOccurredAt)
	diff := oa.Sub(loa)
	if diff < 0 {
		return false
	}

	return true
}
