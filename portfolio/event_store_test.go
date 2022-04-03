package portfolio

import (
	"os"
	"reflect"
	"testing"
)

func TestInMemoryEventStream(t *testing.T) {
	os := InMemoryEventStream{}
	os.Add(Event{
		"EventName",
		map[string]interface{}{
			"foo": "bar",
		},
	})
	os.Add(Event{
		"EventName2",
		map[string]interface{}{
			"foo": "buz",
			"number": 3,
		},
	})

	got := os.Get()
	want := []Event{
		{
			"EventName", map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			"EventName2", map[string]interface{}{
				"foo": "buz",
				"number": 3,
			},
		},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Event store state unequal. Expected:%#v Got:%#v", want, got)
	}
}

func TestFileSystemEventStream(t *testing.T) {
	tmpStorePath := "./tmp/"
	tmpStoreFile := "test_events.gob"

	os.Mkdir(tmpStorePath, 0777)
	eventStream := FileSystemEventStream{tmpStorePath, tmpStoreFile}
	eventStream.Add(Event{
		"EventName",
		map[string]interface{}{
			"foo": "bar",
		},
	})
	eventStream.Add(Event{
		"EventName2",
		map[string]interface{}{
			"foo": "buz",
			"number": 3,
		},
	})

	got := eventStream.Get()
	want := []Event{
		{
			"EventName", map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			"EventName2", map[string]interface{}{
				"foo": "buz",
				"number": 3,
			},
		},
	}

	if reflect.DeepEqual(got, want) == false {
		t.Errorf("Event stream unequal. Expected:%#v Got:%#v", want, got)
	}

	t.Cleanup(func() {
		os.Remove(tmpStorePath + tmpStoreFile)
		os.Remove(tmpStorePath)
	})
}
