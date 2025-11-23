package calendar

import (
	"testing"
	"time"
)

func TestNewEntry(t *testing.T) {
	userID := "12345"
	title := "Test Event"
	datetime := time.Now()
	duration := 60

	entry := NewEntry(userID, title, datetime, duration)

	if entry.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, entry.UserID)
	}

	if entry.Title != title {
		t.Errorf("Expected Title %s, got %s", title, entry.Title)
	}

	if entry.Duration != duration {
		t.Errorf("Expected Duration %d, got %d", duration, entry.Duration)
	}

	if entry.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if len(entry.ID) != 16 {
		t.Errorf("Expected ID length 16, got %d", len(entry.ID))
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		entry   *Entry
		wantErr bool
	}{
		{
			name: "valid entry",
			entry: &Entry{
				UserID:   "12345",
				Title:    "Test",
				DateTime: time.Now(),
				Duration: 60,
			},
			wantErr: false,
		},
		{
			name: "empty user_id",
			entry: &Entry{
				UserID:   "",
				Title:    "Test",
				DateTime: time.Now(),
				Duration: 60,
			},
			wantErr: true,
		},
		{
			name: "empty title",
			entry: &Entry{
				UserID:   "12345",
				Title:    "",
				DateTime: time.Now(),
				Duration: 60,
			},
			wantErr: true,
		},
		{
			name: "zero datetime",
			entry: &Entry{
				UserID:   "12345",
				Title:    "Test",
				DateTime: time.Time{},
				Duration: 60,
			},
			wantErr: true,
		},
		{
			name: "zero duration",
			entry: &Entry{
				UserID:   "12345",
				Title:    "Test",
				DateTime: time.Now(),
				Duration: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndTime(t *testing.T) {
	now := time.Now()
	entry := &Entry{
		DateTime: now,
		Duration: 60,
	}

	expected := now.Add(60 * time.Minute)
	got := entry.EndTime()

	if !got.Equal(expected) {
		t.Errorf("EndTime() = %v, want %v", got, expected)
	}
}

func TestIsPastFutureCurrent(t *testing.T) {
	now := time.Now()

	pastEntry := &Entry{
		DateTime: now.Add(-2 * time.Hour),
		Duration: 60,
	}

	futureEntry := &Entry{
		DateTime: now.Add(2 * time.Hour),
		Duration: 60,
	}

	currentEntry := &Entry{
		DateTime: now.Add(-30 * time.Minute),
		Duration: 60,
	}

	if !pastEntry.IsPast() {
		t.Error("Expected past entry to be in the past")
	}

	if !futureEntry.IsFuture() {
		t.Error("Expected future entry to be in the future")
	}

	if !currentEntry.IsCurrent() {
		t.Error("Expected current entry to be current")
	}
}

func TestGenerateFilename(t *testing.T) {
	datetime, _ := time.Parse("2006-01-02 15:04", "2025-11-21 14:30")

	tests := []struct {
		name     string
		entry    *Entry
		expected string
	}{
		{
			name: "simple title",
			entry: &Entry{
				DateTime: datetime,
				Title:    "Meeting",
			},
			expected: "14-30-meeting",
		},
		{
			name: "title with spaces",
			entry: &Entry{
				DateTime: datetime,
				Title:    "Team Stand Up",
			},
			expected: "14-30-team-stand-up",
		},
		{
			name: "title with accents",
			entry: &Entry{
				DateTime: datetime,
				Title:    "Reunión Técnica",
			},
			expected: "14-30-reunion-tecnica",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entry.GenerateFilename()
			if got != tt.expected {
				t.Errorf("GenerateFilename() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTags(t *testing.T) {
	entry := NewEntry("12345", "Test", time.Now(), 60)

	// Test HasTag
	if entry.HasTag("trabajo") {
		t.Error("Expected entry to not have tag 'trabajo'")
	}

	// Test AddTag
	entry.AddTag("trabajo")
	if !entry.HasTag("trabajo") {
		t.Error("Expected entry to have tag 'trabajo'")
	}

	// Test duplicate AddTag
	entry.AddTag("trabajo")
	if len(entry.Tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(entry.Tags))
	}

	// Test RemoveTag
	entry.RemoveTag("trabajo")
	if entry.HasTag("trabajo") {
		t.Error("Expected entry to not have tag 'trabajo' after removal")
	}
}
