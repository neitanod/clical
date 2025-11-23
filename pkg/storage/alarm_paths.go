package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sebasvalencia/clical/pkg/alarm"
)

// AlarmPaths maneja las rutas de archivos para alarmas
type AlarmPaths struct {
	baseDir string
	userID  string
}

// NewAlarmPaths crea un nuevo AlarmPaths
func NewAlarmPaths(baseDir, userID string) *AlarmPaths {
	return &AlarmPaths{
		baseDir: baseDir,
		userID:  userID,
	}
}

// UserAlarmsDir retorna el directorio base de alarmas del usuario
// Formato: <baseDir>/users/<userID>/alarms
func (ap *AlarmPaths) UserAlarmsDir() string {
	return filepath.Join(ap.baseDir, "users", ap.userID, "alarms")
}

// PendingDir retorna el directorio de alarmas pendientes (one-time)
func (ap *AlarmPaths) PendingDir() string {
	return filepath.Join(ap.UserAlarmsDir(), "pending")
}

// RecurringDir retorna el directorio de alarmas recurrentes
func (ap *AlarmPaths) RecurringDir(recurrence alarm.Recurrence) string {
	return filepath.Join(ap.UserAlarmsDir(), "recurring", string(recurrence))
}

// PastDir retorna el directorio de alarmas pasadas
func (ap *AlarmPaths) PastDir(recurrence alarm.Recurrence) string {
	if recurrence == alarm.RecurrenceOnce {
		return filepath.Join(ap.UserAlarmsDir(), "past", "one-time")
	}
	return filepath.Join(ap.UserAlarmsDir(), "past", "recurring", string(recurrence))
}

// PendingFile retorna la ruta completa para una alarma one-time
func (ap *AlarmPaths) PendingFile(filename string) string {
	return filepath.Join(ap.PendingDir(), filename)
}

// RecurringFile retorna la ruta completa para una alarma recurrente
func (ap *AlarmPaths) RecurringFile(recurrence alarm.Recurrence, filename string) string {
	return filepath.Join(ap.RecurringDir(recurrence), filename)
}

// PastFile retorna la ruta completa para una alarma pasada
func (ap *AlarmPaths) PastFile(recurrence alarm.Recurrence, filename string) string {
	return filepath.Join(ap.PastDir(recurrence), filename)
}

// EnsureAlarmDirs crea todos los directorios necesarios para alarmas
func (ap *AlarmPaths) EnsureAlarmDirs() error {
	dirs := []string{
		ap.PendingDir(),
		ap.RecurringDir(alarm.RecurrenceDaily),
		ap.RecurringDir(alarm.RecurrenceWeekly),
		ap.RecurringDir(alarm.RecurrenceMonthly),
		ap.RecurringDir(alarm.RecurrenceYearly),
		ap.PastDir(alarm.RecurrenceOnce),
		ap.PastDir(alarm.RecurrenceDaily),
		ap.PastDir(alarm.RecurrenceWeekly),
		ap.PastDir(alarm.RecurrenceMonthly),
		ap.PastDir(alarm.RecurrenceYearly),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create alarm directory %s: %w", dir, err)
		}
	}

	return nil
}
