package alarm

import (
	"strings"
	"testing"
	"time"
)

func TestNewAlarm(t *testing.T) {
	context := "Test alarm context"
	recurrence := RecurrenceOnce

	alarm := NewAlarm(context, recurrence)

	if alarm == nil {
		t.Fatal("NewAlarm returned nil")
	}

	if alarm.Context != context {
		t.Errorf("expected context %s, got %s", context, alarm.Context)
	}

	if alarm.Recurrence != recurrence {
		t.Errorf("expected recurrence %s, got %s", recurrence, alarm.Recurrence)
	}

	if alarm.ID == "" {
		t.Error("ID should not be empty")
	}

	if !strings.HasPrefix(alarm.ID, "alarm_once_") {
		t.Errorf("ID should have correct prefix, got: %s", alarm.ID)
	}

	if alarm.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		alarm     *Alarm
		wantError bool
	}{
		{
			name: "valid one-time alarm",
			alarm: &Alarm{
				ID:         "alarm_test_001",
				Context:    "Valid context",
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceOnce,
			},
			wantError: false,
		},
		{
			name: "valid recurring alarm with expiration",
			alarm: &Alarm{
				ID:         "alarm_test_002",
				Context:    "Valid recurring",
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceDaily,
				ExpiresAt:  ptrTime(time.Now().Add(24 * time.Hour)),
			},
			wantError: false,
		},
		{
			name: "missing ID",
			alarm: &Alarm{
				Context:    "No ID",
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceOnce,
			},
			wantError: true,
		},
		{
			name: "missing context",
			alarm: &Alarm{
				ID:         "alarm_test_003",
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceOnce,
			},
			wantError: true,
		},
		{
			name: "context too long",
			alarm: &Alarm{
				ID:         "alarm_test_004",
				Context:    strings.Repeat("a", 501),
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceOnce,
			},
			wantError: true,
		},
		{
			name: "invalid recurrence",
			alarm: &Alarm{
				ID:         "alarm_test_005",
				Context:    "Invalid recurrence",
				CreatedAt:  time.Now(),
				Recurrence: "invalid",
			},
			wantError: true,
		},
		{
			name: "missing created_at",
			alarm: &Alarm{
				ID:         "alarm_test_006",
				Context:    "No created at",
				Recurrence: RecurrenceOnce,
			},
			wantError: true,
		},
		{
			name: "expires_at on one-time alarm",
			alarm: &Alarm{
				ID:         "alarm_test_007",
				Context:    "Should not have expires",
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceOnce,
				ExpiresAt:  ptrTime(time.Now().Add(24 * time.Hour)),
			},
			wantError: true,
		},
		{
			name: "expires_at in the past",
			alarm: &Alarm{
				ID:         "alarm_test_008",
				Context:    "Expired already",
				CreatedAt:  time.Now(),
				Recurrence: RecurrenceDaily,
				ExpiresAt:  ptrTime(time.Now().Add(-24 * time.Hour)),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.alarm.Validate()
			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

func TestIsExpired(t *testing.T) {
	tests := []struct {
		name       string
		alarm      *Alarm
		wantExpired bool
	}{
		{
			name: "no expiration",
			alarm: &Alarm{
				ExpiresAt: nil,
			},
			wantExpired: false,
		},
		{
			name: "expires in future",
			alarm: &Alarm{
				ExpiresAt: ptrTime(time.Now().Add(24 * time.Hour)),
			},
			wantExpired: false,
		},
		{
			name: "expires in past",
			alarm: &Alarm{
				ExpiresAt: ptrTime(time.Now().Add(-24 * time.Hour)),
			},
			wantExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alarm.IsExpired()
			if got != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", got, tt.wantExpired)
			}
		})
	}
}

func TestShouldExecute(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		alarm       *Alarm
		shouldExec  bool
	}{
		{
			name: "one-time alarm not expired",
			alarm: &Alarm{
				Recurrence: RecurrenceOnce,
			},
			shouldExec: true,
		},
		{
			name: "recurring alarm not expired",
			alarm: &Alarm{
				Recurrence: RecurrenceDaily,
				ExpiresAt:  ptrTime(now.Add(24 * time.Hour)),
			},
			shouldExec: true,
		},
		{
			name: "recurring alarm expired",
			alarm: &Alarm{
				Recurrence: RecurrenceDaily,
				ExpiresAt:  ptrTime(now.Add(-24 * time.Hour)),
			},
			shouldExec: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alarm.ShouldExecute(now)
			if got != tt.shouldExec {
				t.Errorf("ShouldExecute() = %v, want %v", got, tt.shouldExec)
			}
		})
	}
}

func TestWithScheduledFor(t *testing.T) {
	alarm := NewAlarm("test", RecurrenceOnce)
	scheduledFor := time.Date(2025, 11, 24, 10, 0, 0, 0, time.UTC)

	result := alarm.WithScheduledFor(scheduledFor)

	if result != alarm {
		t.Error("WithScheduledFor should return the same alarm")
	}

	if !alarm.ScheduledFor.Equal(scheduledFor) {
		t.Errorf("ScheduledFor = %v, want %v", alarm.ScheduledFor, scheduledFor)
	}
}

func TestWithExecutedAt(t *testing.T) {
	alarm := NewAlarm("test", RecurrenceOnce)
	executedAt := time.Date(2025, 11, 24, 10, 0, 0, 0, time.UTC)

	result := alarm.WithExecutedAt(executedAt)

	if result != alarm {
		t.Error("WithExecutedAt should return the same alarm")
	}

	if alarm.ExecutedAt == nil {
		t.Fatal("ExecutedAt should not be nil")
	}

	if !alarm.ExecutedAt.Equal(executedAt) {
		t.Errorf("ExecutedAt = %v, want %v", *alarm.ExecutedAt, executedAt)
	}
}

func TestClone(t *testing.T) {
	expiresAt := time.Now().Add(24 * time.Hour)
	executedAt := time.Now()

	original := &Alarm{
		ID:           "test_id",
		Context:      "test context",
		CreatedAt:    time.Now(),
		Recurrence:   RecurrenceDaily,
		ExpiresAt:    &expiresAt,
		ScheduledFor: time.Now(),
		ExecutedAt:   &executedAt,
	}

	clone := original.Clone()

	// Verify values are equal
	if clone.ID != original.ID {
		t.Error("Clone ID mismatch")
	}
	if clone.Context != original.Context {
		t.Error("Clone Context mismatch")
	}
	if !clone.CreatedAt.Equal(original.CreatedAt) {
		t.Error("Clone CreatedAt mismatch")
	}
	if clone.Recurrence != original.Recurrence {
		t.Error("Clone Recurrence mismatch")
	}

	// Verify pointers are different (deep copy)
	if clone.ExpiresAt == original.ExpiresAt {
		t.Error("Clone should have different ExpiresAt pointer")
	}
	if clone.ExecutedAt == original.ExecutedAt {
		t.Error("Clone should have different ExecutedAt pointer")
	}

	// Verify pointer values are equal
	if !clone.ExpiresAt.Equal(*original.ExpiresAt) {
		t.Error("Clone ExpiresAt value mismatch")
	}
	if !clone.ExecutedAt.Equal(*original.ExecutedAt) {
		t.Error("Clone ExecutedAt value mismatch")
	}
}

// Helper function to create time pointer
func ptrTime(t time.Time) *time.Time {
	return &t
}
