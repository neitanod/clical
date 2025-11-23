package user

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	id := "12345"
	name := "Test User"
	timezone := "America/Argentina/Buenos_Aires"

	user := NewUser(id, name, timezone)

	if user.ID != id {
		t.Errorf("Expected ID %s, got %s", id, user.ID)
	}

	if user.Name != name {
		t.Errorf("Expected Name %s, got %s", name, user.Name)
	}

	if user.Timezone != timezone {
		t.Errorf("Expected Timezone %s, got %s", timezone, user.Timezone)
	}

	// Verify default config
	if user.Config.DefaultDuration != 60 {
		t.Errorf("Expected DefaultDuration 60, got %d", user.Config.DefaultDuration)
	}

	if user.Config.FirstDayOfWeek != 1 {
		t.Errorf("Expected FirstDayOfWeek 1 (Monday), got %d", user.Config.FirstDayOfWeek)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				ID:       "12345",
				Name:     "Test",
				Timezone: "UTC",
				Config:   DefaultConfig(),
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			user: &User{
				ID:       "",
				Name:     "Test",
				Timezone: "UTC",
				Config:   DefaultConfig(),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			user: &User{
				ID:       "12345",
				Name:     "",
				Timezone: "UTC",
				Config:   DefaultConfig(),
			},
			wantErr: true,
		},
		{
			name: "empty timezone",
			user: &User{
				ID:       "12345",
				Name:     "Test",
				Timezone: "",
				Config:   DefaultConfig(),
			},
			wantErr: true,
		},
		{
			name: "invalid timezone",
			user: &User{
				ID:       "12345",
				Name:     "Test",
				Timezone: "Invalid/Timezone",
				Config:   DefaultConfig(),
			},
			wantErr: true,
		},
		{
			name: "invalid default duration",
			user: &User{
				ID:       "12345",
				Name:     "Test",
				Timezone: "UTC",
				Config: UserConfig{
					DefaultDuration: 0,
					DateFormat:      "2006-01-02",
					TimeFormat:      "15:04",
					FirstDayOfWeek:  1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocation(t *testing.T) {
	user := NewUser("12345", "Test", "America/Argentina/Buenos_Aires")

	loc, err := user.Location()
	if err != nil {
		t.Fatalf("Location() error = %v", err)
	}

	if loc == nil {
		t.Fatal("Expected non-nil location")
	}

	// Verify timezone name
	expectedName := "America/Argentina/Buenos_Aires"
	if loc.String() != expectedName {
		t.Errorf("Expected location %s, got %s", expectedName, loc.String())
	}
}

func TestFormatFunctions(t *testing.T) {
	user := NewUser("12345", "Test", "UTC")
	testTime := time.Date(2025, 11, 21, 14, 30, 0, 0, time.UTC)

	// Test FormatDate
	expectedDate := "2025-11-21"
	gotDate := user.FormatDate(testTime)
	if gotDate != expectedDate {
		t.Errorf("FormatDate() = %s, want %s", gotDate, expectedDate)
	}

	// Test FormatTime
	expectedTime := "14:30"
	gotTime := user.FormatTime(testTime)
	if gotTime != expectedTime {
		t.Errorf("FormatTime() = %s, want %s", gotTime, expectedTime)
	}

	// Test FormatDateTime
	expectedDateTime := "2025-11-21 14:30"
	gotDateTime := user.FormatDateTime(testTime)
	if gotDateTime != expectedDateTime {
		t.Errorf("FormatDateTime() = %s, want %s", gotDateTime, expectedDateTime)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.DefaultDuration != 60 {
		t.Errorf("Expected DefaultDuration 60, got %d", config.DefaultDuration)
	}

	if config.DateFormat != "2006-01-02" {
		t.Errorf("Expected DateFormat 2006-01-02, got %s", config.DateFormat)
	}

	if config.TimeFormat != "15:04" {
		t.Errorf("Expected TimeFormat 15:04, got %s", config.TimeFormat)
	}

	if config.FirstDayOfWeek != 1 {
		t.Errorf("Expected FirstDayOfWeek 1, got %d", config.FirstDayOfWeek)
	}
}
