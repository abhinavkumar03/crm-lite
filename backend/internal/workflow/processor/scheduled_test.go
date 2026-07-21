package processor

import (
	"testing"
	"time"
)

func TestScheduledDue(t *testing.T) {
	now := time.Date(2026, 7, 20, 9, 30, 0, 0, time.UTC) // Monday
	if !scheduledDue(map[string]any{"hour": 9.0, "minute": 30.0}, now) {
		t.Fatal("expected due at 09:30")
	}
	if scheduledDue(map[string]any{"hour": 10.0, "minute": 30.0}, now) {
		t.Fatal("expected not due at different hour")
	}
	// days_of_week: Monday = 1
	if !scheduledDue(map[string]any{
		"hour": 9.0, "minute": 30.0, "days_of_week": []any{float64(1)},
	}, now) {
		t.Fatal("expected due on Monday")
	}
	if scheduledDue(map[string]any{
		"hour": 9.0, "minute": 30.0, "days_of_week": []any{float64(0)},
	}, now) {
		t.Fatal("expected not due on Sunday-only schedule")
	}
}
