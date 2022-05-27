package infrastructure_test

import (
	"os"
	"reflect"
	"stock-monitor/infrastructure"
	"testing"
	"time"
)

var tmpStorePath = "./tmp/"
var tmpStoreFile = "test_events.gob"

func TestCanNotAddEventsWithInvalidOccurredAtDateFormat(t *testing.T) {
	inMemoryEventStream := infrastructure.InMemoryEventStream{}
	fileSystemEventStream := setUpFileSystemEventStream()

	eventStreams := map[string]infrastructure.EventStream{
		"InMemoryEventStream":   &inMemoryEventStream,
		"FileSystemEventStream": &fileSystemEventStream,
	}

	for name, eventStream := range eventStreams {
		t.Run(name, func(t *testing.T) {
			err := eventStream.Add(infrastructure.Event{
				"EventName",
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"occurred_at": "Foo",
				},
			})

			_, ok := err.(*infrastructure.UnsupportedDateFormatError)
			if !ok {
				t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
			}
		})
	}

	t.Cleanup(func() {
		cleanUpFileSystemEventStream()
	})
}

func TestCanNotAddEventsWithoutOccurredAt(t *testing.T) {
	inMemoryEventStream := infrastructure.InMemoryEventStream{}
	fileSystemEventStream := setUpFileSystemEventStream()

	eventStreams := map[string]infrastructure.EventStream{
		"InMemoryEventStream":   &inMemoryEventStream,
		"FileSystemEventStream": &fileSystemEventStream,
	}
	for name, eventStream := range eventStreams {
		t.Run(name, func(t *testing.T) {
			err := eventStream.Add(infrastructure.Event{
				"EventName",
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{},
			})

			_, ok := err.(*infrastructure.UnsupportedDateFormatError)
			if !ok {
				t.Errorf("Expected UnsupportedDateFormatError but got %#v", err)
			}
		})
	}

	t.Cleanup(func() {
		cleanUpFileSystemEventStream()
	})
}

func TestOccurredAtCanNotBeInTheFuture(t *testing.T) {
	inMemoryEventStream := infrastructure.InMemoryEventStream{}
	fileSystemEventStream := setUpFileSystemEventStream()
	today := time.Now()

	eventStreams := map[string]infrastructure.EventStream{
		"InMemoryEventStream":   &inMemoryEventStream,
		"FileSystemEventStream": &fileSystemEventStream,
	}
	for name, eventStream := range eventStreams {
		t.Run(name, func(t *testing.T) {
			err := eventStream.Add(infrastructure.Event{
				"EventName",
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"occurred_at": today.AddDate(0, 0, 1).Format("2006-01-02"),
				},
			})

			_, ok := err.(*infrastructure.InvalidDateError)
			if !ok {
				t.Errorf("Expected InvalidDateError but got %#v", err)
			}
		})
	}

	t.Cleanup(func() {
		cleanUpFileSystemEventStream()
	})
}

func TestOccurredAtCanNotBeOlderThanOccurredAtOfLastEvent(t *testing.T) {
	inMemoryEventStream := infrastructure.InMemoryEventStream{}
	fileSystemEventStream := setUpFileSystemEventStream()

	eventStreams := map[string]infrastructure.EventStream{
		"InMemoryEventStream":   &inMemoryEventStream,
		"FileSystemEventStream": &fileSystemEventStream,
	}
	for name, eventStream := range eventStreams {
		t.Run(name, func(t *testing.T) {
			eventStream.Add(infrastructure.Event{
				"EventName",
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"occurred_at": "2000-01-02",
				},
			})
			err := eventStream.Add(infrastructure.Event{
				"EventName",
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"occurred_at": "2000-01-01",
				},
			})

			_, ok := err.(*infrastructure.InvalidDateError)
			if !ok {
				t.Errorf("Expected InvalidDateError but got %#v", err)
			}
		})
	}

	t.Cleanup(func() {
		cleanUpFileSystemEventStream()
	})
}

func TestAddEvents(t *testing.T) {
	inMemoryEventStream := infrastructure.InMemoryEventStream{}
	fileSystemEventStream := setUpFileSystemEventStream()

	eventStreams := map[string]infrastructure.EventStream{
		"InMemoryEventStream":   &inMemoryEventStream,
		"FileSystemEventStream": &fileSystemEventStream,
	}

	for name, eventStream := range eventStreams {
		t.Run(name, func(t *testing.T) {
			eventStream.Add(infrastructure.Event{
				"EventName",
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"occurred_at": "2000-01-01",
				},
			})
			eventStream.Add(infrastructure.Event{
				"EventName2",
				map[string]interface{}{
					"foo":    "buz",
					"number": 3,
				},
				map[string]interface{}{
					"occurred_at": "2000-01-01",
				},
			})

			got := eventStream.Get()
			want := []infrastructure.Event{
				{
					"EventName", map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"occurred_at": "2000-01-01",
					},
				},
				{
					"EventName2", map[string]interface{}{
						"foo":    "buz",
						"number": 3,
					},
					map[string]interface{}{
						"occurred_at": "2000-01-01",
					},
				},
			}

			if reflect.DeepEqual(got, want) == false {
				t.Errorf("Event store state unequal. Expected:%#v Got:%#v", want, got)
			}
		})
	}

	t.Cleanup(func() {
		cleanUpFileSystemEventStream()
	})
}

func setUpFileSystemEventStream() infrastructure.FileSystemEventStream {
	os.Mkdir(tmpStorePath, 0777)
	return infrastructure.FileSystemEventStream{tmpStorePath, tmpStoreFile}
}

func cleanUpFileSystemEventStream() {
	os.Remove(tmpStorePath + tmpStoreFile)
	os.Remove(tmpStorePath)
}
