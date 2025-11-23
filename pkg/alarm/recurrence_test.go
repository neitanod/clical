package alarm

import (
	"testing"
	"time"
)

func TestRecurrenceValid(t *testing.T) {
	tests := []struct {
		name       string
		recurrence Recurrence
		wantValid  bool
	}{
		{"once is valid", RecurrenceOnce, true},
		{"daily is valid", RecurrenceDaily, true},
		{"weekly is valid", RecurrenceWeekly, true},
		{"monthly is valid", RecurrenceMonthly, true},
		{"yearly is valid", RecurrenceYearly, true},
		{"invalid recurrence", Recurrence("invalid"), false},
		{"empty recurrence", Recurrence(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.recurrence.Valid()
			if got != tt.wantValid {
				t.Errorf("Valid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}

func TestRecurrenceString(t *testing.T) {
	tests := []struct {
		name       string
		recurrence Recurrence
		want       string
	}{
		{"once", RecurrenceOnce, "once"},
		{"daily", RecurrenceDaily, "daily"},
		{"weekly", RecurrenceWeekly, "weekly"},
		{"monthly", RecurrenceMonthly, "monthly"},
		{"yearly", RecurrenceYearly, "yearly"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.recurrence.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeeklyScheduleFilename(t *testing.T) {
	tests := []struct {
		name     string
		schedule WeeklySchedule
		want     string
	}{
		{
			name:     "monday morning",
			schedule: WeeklySchedule{Weekday: time.Monday, Hour: 9, Minute: 30},
			want:     "monday_09-30-00.json",
		},
		{
			name:     "friday afternoon",
			schedule: WeeklySchedule{Weekday: time.Friday, Hour: 17, Minute: 0},
			want:     "friday_17-00-00.json",
		},
		{
			name:     "sunday midnight",
			schedule: WeeklySchedule{Weekday: time.Sunday, Hour: 0, Minute: 0},
			want:     "sunday_00-00-00.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.Filename()
			if got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonthlyScheduleFilename(t *testing.T) {
	tests := []struct {
		name     string
		schedule MonthlySchedule
		want     string
	}{
		{
			name:     "first day of month",
			schedule: MonthlySchedule{Day: 1, Hour: 9, Minute: 0},
			want:     "01_09-00-00.json",
		},
		{
			name:     "middle of month",
			schedule: MonthlySchedule{Day: 15, Hour: 14, Minute: 30},
			want:     "15_14-30-00.json",
		},
		{
			name:     "end of month",
			schedule: MonthlySchedule{Day: 31, Hour: 23, Minute: 59},
			want:     "31_23-59-00.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.Filename()
			if got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYearlyScheduleFilename(t *testing.T) {
	tests := []struct {
		name     string
		schedule YearlySchedule
		want     string
	}{
		{
			name:     "new year",
			schedule: YearlySchedule{Month: time.January, Day: 1, Hour: 0, Minute: 0},
			want:     "01-01_00-00-00.json",
		},
		{
			name:     "project anniversary",
			schedule: YearlySchedule{Month: time.November, Day: 21, Hour: 10, Minute: 0},
			want:     "11-21_10-00-00.json",
		},
		{
			name:     "christmas",
			schedule: YearlySchedule{Month: time.December, Day: 25, Hour: 8, Minute: 0},
			want:     "12-25_08-00-00.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.Filename()
			if got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDailyScheduleFilename(t *testing.T) {
	tests := []struct {
		name     string
		schedule DailySchedule
		want     string
	}{
		{
			name:     "morning",
			schedule: DailySchedule{Hour: 9, Minute: 0},
			want:     "09-00-00.json",
		},
		{
			name:     "afternoon",
			schedule: DailySchedule{Hour: 14, Minute: 30},
			want:     "14-30-00.json",
		},
		{
			name:     "midnight",
			schedule: DailySchedule{Hour: 0, Minute: 0},
			want:     "00-00-00.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schedule.Filename()
			if got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWeekday(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Weekday
		wantErr bool
	}{
		{"sunday", "sunday", time.Sunday, false},
		{"monday", "monday", time.Monday, false},
		{"tuesday", "tuesday", time.Tuesday, false},
		{"wednesday", "wednesday", time.Wednesday, false},
		{"thursday", "thursday", time.Thursday, false},
		{"friday", "friday", time.Friday, false},
		{"saturday", "saturday", time.Saturday, false},
		{"uppercase", "MONDAY", time.Monday, false},
		{"mixed case", "FrIdAy", time.Friday, false},
		{"with spaces", "  monday  ", time.Monday, false},
		{"invalid", "invalid", 0, true},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWeekday(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWeekday() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOneTimeFilename(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "specific datetime",
			time: time.Date(2025, 11, 23, 14, 30, 0, 0, time.UTC),
			want: "2025-11-23_14-30-00.json",
		},
		{
			name: "midnight",
			time: time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
			want: "2025-12-25_00-00-00.json",
		},
		{
			name: "end of day",
			time: time.Date(2025, 1, 1, 23, 59, 0, 0, time.UTC),
			want: "2025-01-01_23-59-00.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OneTimeFilename(tt.time)
			if got != tt.want {
				t.Errorf("OneTimeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrentFilenames(t *testing.T) {
	// Use Friday 2025-11-21 which is universally Friday in UTC
	testTime := time.Date(2025, 11, 21, 14, 30, 45, 0, time.UTC) // Friday in UTC

	t.Run("CurrentDailyFilename", func(t *testing.T) {
		got := CurrentDailyFilename(testTime)
		want := "14-30-00.json"
		if got != want {
			t.Errorf("CurrentDailyFilename() = %v, want %v", got, want)
		}
	})

	t.Run("CurrentWeeklyFilename", func(t *testing.T) {
		got := CurrentWeeklyFilename(testTime)
		want := "friday_14-30-00.json"
		if got != want {
			t.Errorf("CurrentWeeklyFilename() = %v, want %v", got, want)
		}
	})

	t.Run("CurrentMonthlyFilename", func(t *testing.T) {
		got := CurrentMonthlyFilename(testTime)
		want := "21_14-30-00.json"
		if got != want {
			t.Errorf("CurrentMonthlyFilename() = %v, want %v", got, want)
		}
	})

	t.Run("CurrentYearlyFilename", func(t *testing.T) {
		got := CurrentYearlyFilename(testTime)
		want := "11-21_14-30-00.json"
		if got != want {
			t.Errorf("CurrentYearlyFilename() = %v, want %v", got, want)
		}
	})
}

func TestRoundToMinute(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
		want  time.Time
	}{
		{
			name:  "with seconds",
			input: time.Date(2025, 11, 23, 14, 30, 45, 0, time.UTC),
			want:  time.Date(2025, 11, 23, 14, 30, 0, 0, time.UTC),
		},
		{
			name:  "already rounded",
			input: time.Date(2025, 11, 23, 14, 30, 0, 0, time.UTC),
			want:  time.Date(2025, 11, 23, 14, 30, 0, 0, time.UTC),
		},
		{
			name:  "with nanoseconds",
			input: time.Date(2025, 11, 23, 14, 30, 15, 123456789, time.UTC),
			want:  time.Date(2025, 11, 23, 14, 30, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundToMinute(tt.input)
			if !got.Equal(tt.want) {
				t.Errorf("RoundToMinute() = %v, want %v", got, tt.want)
			}
		})
	}
}
