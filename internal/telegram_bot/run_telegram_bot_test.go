package telegram_bot

import (
	"testing"
)

func TestDurationUntilNextExecution(t *testing.T) {
	duration := durationUntilNextExecution()
	if duration <= 0 {
		t.Errorf("Expected positive duration, got %v", duration)
	}
}
