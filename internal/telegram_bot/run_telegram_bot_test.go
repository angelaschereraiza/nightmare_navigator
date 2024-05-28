package telegram_bot

import (
	"testing"
	"time"
)

func TestDurationUntilNextExecution(t *testing.T) {
	duration := durationUntilNextExecution()
	if duration <= 0 {
		t.Errorf("Expected positive duration, got %v", duration)
	}


	// Test edge case: current time is just before 02:59 AM
	now := time.Date(2024, 5, 27, 2, 59, 59, 0, time.UTC)
	expected := 1 * time.Second
	duration = timeUntilNextExecutionFrom(now)
	if duration != expected {
		t.Errorf("Expected %v, got %v", expected, duration)
	}

	// Test edge case: current time is exactly 03:00 AM
	now = time.Date(2024, 5, 27, 3, 0, 0, 0, time.UTC)
	expected = 0 * time.Hour
	duration = timeUntilNextExecutionFrom(now)
	if duration != expected {
		t.Errorf("Expected %v, got %v", expected, duration)
	}

	// Test edge case: current time is just before 03:01 AM
	now = time.Date(2024, 5, 27, 3, 01, 00, 0, time.UTC)
	expected = 23 * time.Hour + 59 * time.Minute + 0 * time.Second
	duration = timeUntilNextExecutionFrom(now)
	if duration != expected {
		t.Errorf("Expected %v, got %v", expected, duration)
	}
}

func timeUntilNextExecutionFrom(now time.Time) time.Duration {
	nextExecution := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	if now.After(nextExecution) {
		nextExecution = nextExecution.Add(24 * time.Hour)
	}
	return nextExecution.Sub(now)
}
