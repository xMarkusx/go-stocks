package portfolio

import (
	"testing"
)

func TestPositionProvidesCurrentValue(t *testing.T) {
	position := Position{"MO", 10}
	values := map[string]float32{"MO":10}
	valueTracker := FakeValueTracker{values}

	got := position.CurrentValue(valueTracker)
	expected := float32(100)

	if got != expected {
		t.Errorf("Expected current value: %v, got %v", expected, got)
	}
}
