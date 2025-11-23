package alarm

import (
	"fmt"
	"strings"
	"time"
)

// Recurrence define el tipo de recurrencia de una alarma
type Recurrence string

const (
	RecurrenceOnce    Recurrence = "once"
	RecurrenceDaily   Recurrence = "daily"
	RecurrenceWeekly  Recurrence = "weekly"
	RecurrenceMonthly Recurrence = "monthly"
	RecurrenceYearly  Recurrence = "yearly"
)

// Valid retorna true si la recurrencia es válida
func (r Recurrence) Valid() bool {
	switch r {
	case RecurrenceOnce, RecurrenceDaily, RecurrenceWeekly, RecurrenceMonthly, RecurrenceYearly:
		return true
	default:
		return false
	}
}

// String implementa Stringer
func (r Recurrence) String() string {
	return string(r)
}

// WeeklySchedule representa un horario semanal
type WeeklySchedule struct {
	Weekday time.Weekday
	Hour    int
	Minute  int
}

// Filename retorna el nombre de archivo para una alarma semanal
// Formato: monday_14-30-00.json
func (w WeeklySchedule) Filename() string {
	weekdayName := strings.ToLower(w.Weekday.String())
	return fmt.Sprintf("%s_%02d-%02d-00.json", weekdayName, w.Hour, w.Minute)
}

// MonthlySchedule representa un horario mensual
type MonthlySchedule struct {
	Day    int // 1-31
	Hour   int
	Minute int
}

// Filename retorna el nombre de archivo para una alarma mensual
// Formato: 15_14-30-00.json
func (m MonthlySchedule) Filename() string {
	return fmt.Sprintf("%02d_%02d-%02d-00.json", m.Day, m.Hour, m.Minute)
}

// YearlySchedule representa un horario anual
type YearlySchedule struct {
	Month  time.Month
	Day    int
	Hour   int
	Minute int
}

// Filename retorna el nombre de archivo para una alarma anual
// Formato: 11-21_14-30-00.json
func (y YearlySchedule) Filename() string {
	return fmt.Sprintf("%02d-%02d_%02d-%02d-00.json", y.Month, y.Day, y.Hour, y.Minute)
}

// DailySchedule representa un horario diario
type DailySchedule struct {
	Hour   int
	Minute int
}

// Filename retorna el nombre de archivo para una alarma diaria
// Formato: 14-30-00.json
func (d DailySchedule) Filename() string {
	return fmt.Sprintf("%02d-%02d-00.json", d.Hour, d.Minute)
}

// ParseWeekday parsea un string a time.Weekday
func ParseWeekday(s string) (time.Weekday, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "sunday":
		return time.Sunday, nil
	case "monday":
		return time.Monday, nil
	case "tuesday":
		return time.Tuesday, nil
	case "wednesday":
		return time.Wednesday, nil
	case "thursday":
		return time.Thursday, nil
	case "friday":
		return time.Friday, nil
	case "saturday":
		return time.Saturday, nil
	default:
		return 0, fmt.Errorf("invalid weekday: %s", s)
	}
}

// OneTimeFilename retorna el nombre de archivo para una alarma one-time
// Formato: 2025-11-23_14-30-00.json
func OneTimeFilename(t time.Time) string {
	return t.Format("2006-01-02_15-04-00") + ".json"
}

// CurrentDailyFilename retorna el filename para la hora actual (daily)
func CurrentDailyFilename(t time.Time) string {
	return DailySchedule{
		Hour:   t.Hour(),
		Minute: t.Minute(),
	}.Filename()
}

// CurrentWeeklyFilename retorna el filename para el día/hora actual (weekly)
func CurrentWeeklyFilename(t time.Time) string {
	return WeeklySchedule{
		Weekday: t.Weekday(),
		Hour:    t.Hour(),
		Minute:  t.Minute(),
	}.Filename()
}

// CurrentMonthlyFilename retorna el filename para el día/hora actual (monthly)
func CurrentMonthlyFilename(t time.Time) string {
	return MonthlySchedule{
		Day:    t.Day(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
	}.Filename()
}

// CurrentYearlyFilename retorna el filename para la fecha/hora actual (yearly)
func CurrentYearlyFilename(t time.Time) string {
	return YearlySchedule{
		Month:  t.Month(),
		Day:    t.Day(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
	}.Filename()
}

// RoundToMinute redondea un time.Time al minuto más cercano (segundos = 0)
func RoundToMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}
