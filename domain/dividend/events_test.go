package dividend_test

import (
	"reflect"
	"stock-monitor/domain/dividend"
	"testing"
)

func TestDividendRecordedEventCanBeCreated(t *testing.T) {
	event := dividend.NewDividendRecordedEvent("MO", 12.34, 23.45, "2000-01-01")

	if event.Name() != dividend.DividendRecordedEventName {
		t.Errorf("Unexpected events name. Expected:%#v Got:%#v", dividend.DividendRecordedEventName, event.Name())
	}

	expectedPayload := map[string]interface{}{
		"ticker": "MO",
		"net":    float32(12.34),
		"gross":  float32(23.45),
		"date":   "2000-01-01",
	}
	if reflect.DeepEqual(event.Payload(), expectedPayload) == false {
		t.Errorf("Unexpected event payload. Expected:%#v Got:%#v", expectedPayload, event.Payload())
	}
}
