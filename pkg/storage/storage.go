package storage

import (
	"time"

	"github.com/sebasvalencia/clical/pkg/alarm"
	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/sebasvalencia/clical/pkg/user"
)

// Storage define la interfaz para almacenamiento de datos
type Storage interface {
	// Entries
	SaveEntry(userID string, entry *calendar.Entry) error
	GetEntry(userID, entryID string) (*calendar.Entry, error)
	ListEntries(userID string, filter *calendar.Filter) ([]*calendar.Entry, error)
	DeleteEntry(userID, entryID string) error
	UpdateEntry(userID string, entry *calendar.Entry) error

	// Users
	SaveUser(user *user.User) error
	GetUser(userID string) (*user.User, error)
	ListUsers() ([]*user.User, error)
	DeleteUser(userID string) error

	// State (para reportes)
	GetReportState(userID string) (*ReportState, error)
	SaveReportState(userID string, state *ReportState) error

	// Alarms
	SaveAlarm(userID string, alarmTime time.Time, recurrence alarm.Recurrence, filename string, alm *alarm.Alarm) error
	GetAlarms(userID string, recurrence alarm.Recurrence, filename string) ([]*alarm.Alarm, error)
	DeleteAlarms(userID string, recurrence alarm.Recurrence, filename string) error
	CheckAlarms(userID string, at time.Time) ([]*alarm.Alarm, error)
	ListActiveAlarms(userID string) ([]*alarm.Alarm, error)
	ListPastAlarms(userID string) ([]*alarm.Alarm, error)
	CancelAlarm(userID string, alarmID string) error
	MoveAlarmsToPast(userID string, recurrence alarm.Recurrence, filename string) error
}

// ReportState almacena el estado de los reportes generados
type ReportState struct {
	LastDailyReport    *string           `json:"last_daily_report,omitempty"`
	LastTomorrowReport *string           `json:"last_tomorrow_report,omitempty"`
	LastUpcomingReport *string           `json:"last_upcoming_report,omitempty"`
	LastWeeklyReport   *string           `json:"last_weekly_report,omitempty"`
	ReportedEvents     map[string]string `json:"reported_events"` // eventID -> timestamp
}

// NewReportState crea un nuevo estado de reportes
func NewReportState() *ReportState {
	return &ReportState{
		ReportedEvents: make(map[string]string),
	}
}
