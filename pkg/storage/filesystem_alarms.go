package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/alarm"
)

// SaveAlarm guarda una alarma en el archivo correspondiente
func (fs *FilesystemStorage) SaveAlarm(userID string, alarmTime time.Time, recurrence alarm.Recurrence, filename string, alm *alarm.Alarm) error {
	if err := alm.Validate(); err != nil {
		return fmt.Errorf("alarma inválida: %w", err)
	}

	ap := NewAlarmPaths(fs.dataDir, userID)
	if err := ap.EnsureAlarmDirs(); err != nil {
		return err
	}

	var filePath string
	if recurrence == alarm.RecurrenceOnce {
		filePath = ap.PendingFile(filename)
	} else {
		filePath = ap.RecurringFile(recurrence, filename)
	}

	// Leer alarmas existentes en el archivo (si existe)
	existingAlarms := []*alarm.Alarm{}
	if data, err := os.ReadFile(filePath); err == nil {
		if err := json.Unmarshal(data, &existingAlarms); err != nil {
			return fmt.Errorf("error leyendo alarmas existentes: %w", err)
		}
	}

	// Agregar la nueva alarma
	existingAlarms = append(existingAlarms, alm)

	// Guardar el archivo actualizado
	jsonData, err := json.MarshalIndent(existingAlarms, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando alarmas: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo alarmas: %w", err)
	}

	return nil
}

// GetAlarms lee todas las alarmas de un archivo
func (fs *FilesystemStorage) GetAlarms(userID string, recurrence alarm.Recurrence, filename string) ([]*alarm.Alarm, error) {
	ap := NewAlarmPaths(fs.dataDir, userID)

	var filePath string
	if recurrence == alarm.RecurrenceOnce {
		filePath = ap.PendingFile(filename)
	} else {
		filePath = ap.RecurringFile(recurrence, filename)
	}

	// Si el archivo no existe, retornar array vacío
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*alarm.Alarm{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo alarmas: %w", err)
	}

	var alarms []*alarm.Alarm
	if err := json.Unmarshal(data, &alarms); err != nil {
		return nil, fmt.Errorf("error deserializando alarmas: %w", err)
	}

	return alarms, nil
}

// DeleteAlarms elimina un archivo de alarmas
func (fs *FilesystemStorage) DeleteAlarms(userID string, recurrence alarm.Recurrence, filename string) error {
	ap := NewAlarmPaths(fs.dataDir, userID)

	var filePath string
	if recurrence == alarm.RecurrenceOnce {
		filePath = ap.PendingFile(filename)
	} else {
		filePath = ap.RecurringFile(recurrence, filename)
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error eliminando alarmas: %w", err)
	}

	return nil
}

// CheckAlarms verifica alarmas que deben ejecutarse en el momento dado
func (fs *FilesystemStorage) CheckAlarms(userID string, at time.Time) ([]*alarm.Alarm, error) {
	ap := NewAlarmPaths(fs.dataDir, userID)
	if err := ap.EnsureAlarmDirs(); err != nil {
		return nil, err
	}

	roundedTime := alarm.RoundToMinute(at)
	result := []*alarm.Alarm{}

	// 1. Chequear alarmas one-time (pending/) - incluir recovery (últimos 60 minutos)
	for i := 0; i <= 60; i++ {
		checkTime := roundedTime.Add(-time.Duration(i) * time.Minute)
		filename := alarm.OneTimeFilename(checkTime)
		filePath := ap.PendingFile(filename)

		if _, err := os.Stat(filePath); err == nil {
			// Archivo existe
			alarms, err := fs.GetAlarms(userID, alarm.RecurrenceOnce, filename)
			if err != nil {
				return nil, err
			}

			// Agregar ScheduledFor a cada alarma
			for _, alm := range alarms {
				alm.WithScheduledFor(checkTime)
				result = append(result, alm)
			}

			// Mover archivo a past/
			if err := fs.MoveAlarmsToPast(userID, alarm.RecurrenceOnce, filename); err != nil {
				return nil, fmt.Errorf("error moviendo alarma a past: %w", err)
			}
		}
	}

	// 2. Chequear alarmas recurrentes (daily, weekly, monthly, yearly) con recovery de 60 minutos
	if err := fs.checkRecurringAlarmsWithRecovery(userID, roundedTime, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// checkRecurringAlarmsWithRecovery chequea alarmas recurrentes con recovery de 60 minutos
func (fs *FilesystemStorage) checkRecurringAlarmsWithRecovery(userID string, roundedTime time.Time, result *[]*alarm.Alarm) error {
	recurrenceTypes := []alarm.Recurrence{
		alarm.RecurrenceDaily,
		alarm.RecurrenceWeekly,
		alarm.RecurrenceMonthly,
		alarm.RecurrenceYearly,
	}

	// Para cada tipo de recurrencia, chequear últimos 60 minutos
	for _, recurrence := range recurrenceTypes {
		for i := 0; i <= 60; i++ {
			checkTime := roundedTime.Add(-time.Duration(i) * time.Minute)

			// Verificar si ya fue ejecutada en este momento
			wasExecuted, err := fs.WasRecurringAlarmExecuted(userID, recurrence, checkTime)
			if err != nil {
				return fmt.Errorf("error checking execution record: %w", err)
			}

			if wasExecuted {
				// Ya fue ejecutada, omitir
				continue
			}

			// Obtener el filename correspondiente al checkTime
			var filename string
			switch recurrence {
			case alarm.RecurrenceDaily:
				filename = alarm.CurrentDailyFilename(checkTime)
			case alarm.RecurrenceWeekly:
				filename = alarm.CurrentWeeklyFilename(checkTime)
			case alarm.RecurrenceMonthly:
				filename = alarm.CurrentMonthlyFilename(checkTime)
			case alarm.RecurrenceYearly:
				filename = alarm.CurrentYearlyFilename(checkTime)
			}

			// Chequear si existe el archivo
			if err := fs.checkAndExecuteRecurringAlarm(userID, recurrence, filename, checkTime, result); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkAndExecuteRecurringAlarm chequea y ejecuta una alarma recurrente
func (fs *FilesystemStorage) checkAndExecuteRecurringAlarm(userID string, recurrence alarm.Recurrence, filename string, at time.Time, result *[]*alarm.Alarm) error {
	ap := NewAlarmPaths(fs.dataDir, userID)
	filePath := ap.RecurringFile(recurrence, filename)

	if _, err := os.Stat(filePath); err == nil {
		alarms, err := fs.GetAlarms(userID, recurrence, filename)
		if err != nil {
			return err
		}

		if len(alarms) == 0 {
			return nil
		}

		hasExpired := false
		activeAlarms := []*alarm.Alarm{}

		for _, alm := range alarms {
			alm.WithScheduledFor(at)

			// Si está expirada, agregar y marcar para mover a past
			if alm.IsExpired() {
				*result = append(*result, alm)
				hasExpired = true
			} else {
				// No expirada, ejecutar
				*result = append(*result, alm)
				activeAlarms = append(activeAlarms, alm)
			}
		}

		// Copiar registro de ejecución a past/ para evitar duplicados
		if len(activeAlarms) > 0 {
			if err := fs.CopyRecurringAlarmExecution(userID, recurrence, activeAlarms, at); err != nil {
				return fmt.Errorf("error copying execution record: %w", err)
			}
		}

		// Si alguna alarma expiró, mover archivo completo a past/
		if hasExpired {
			if err := fs.MoveAlarmsToPast(userID, recurrence, filename); err != nil {
				return fmt.Errorf("error moving expired alarm to past: %w", err)
			}
		}
	}

	return nil
}

// ListActiveAlarms lista todas las alarmas activas
func (fs *FilesystemStorage) ListActiveAlarms(userID string) ([]*alarm.Alarm, error) {
	ap := NewAlarmPaths(fs.dataDir, userID)
	if err := ap.EnsureAlarmDirs(); err != nil {
		return nil, err
	}

	result := []*alarm.Alarm{}

	// 1. Listar alarmas pending (one-time)
	pendingFiles, err := filepath.Glob(filepath.Join(ap.PendingDir(), "*.json"))
	if err != nil {
		return nil, fmt.Errorf("error listando alarmas pending: %w", err)
	}

	for _, file := range pendingFiles {
		filename := filepath.Base(file)
		alarms, err := fs.GetAlarms(userID, alarm.RecurrenceOnce, filename)
		if err != nil {
			return nil, err
		}
		result = append(result, alarms...)
	}

	// 2. Listar alarmas recurrentes
	recurrences := []alarm.Recurrence{
		alarm.RecurrenceDaily,
		alarm.RecurrenceWeekly,
		alarm.RecurrenceMonthly,
		alarm.RecurrenceYearly,
	}

	for _, rec := range recurrences {
		files, err := filepath.Glob(filepath.Join(ap.RecurringDir(rec), "*.json"))
		if err != nil {
			return nil, fmt.Errorf("error listando alarmas %s: %w", rec, err)
		}

		for _, file := range files {
			filename := filepath.Base(file)
			alarms, err := fs.GetAlarms(userID, rec, filename)
			if err != nil {
				return nil, err
			}
			result = append(result, alarms...)
		}
	}

	return result, nil
}

// ListPastAlarms lista todas las alarmas pasadas
func (fs *FilesystemStorage) ListPastAlarms(userID string) ([]*alarm.Alarm, error) {
	ap := NewAlarmPaths(fs.dataDir, userID)
	result := []*alarm.Alarm{}

	// Listar todas las carpetas de past/
	recurrences := []alarm.Recurrence{
		alarm.RecurrenceOnce,
		alarm.RecurrenceDaily,
		alarm.RecurrenceWeekly,
		alarm.RecurrenceMonthly,
		alarm.RecurrenceYearly,
	}

	for _, rec := range recurrences {
		pastDir := ap.PastDir(rec)
		files, err := filepath.Glob(filepath.Join(pastDir, "*.json"))
		if err != nil {
			continue // Puede no existir el directorio
		}

		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			var alarms []*alarm.Alarm
			if err := json.Unmarshal(data, &alarms); err != nil {
				continue
			}

			result = append(result, alarms...)
		}
	}

	return result, nil
}

// CancelAlarm cancela (elimina) una alarma por ID
func (fs *FilesystemStorage) CancelAlarm(userID string, alarmID string) error {
	ap := NewAlarmPaths(fs.dataDir, userID)

	// Buscar en pending/
	if err := fs.cancelAlarmInDir(ap.PendingDir(), alarmID, alarm.RecurrenceOnce); err == nil {
		return nil // Encontrada y eliminada
	}

	// Buscar en recurring/
	recurrences := []alarm.Recurrence{
		alarm.RecurrenceDaily,
		alarm.RecurrenceWeekly,
		alarm.RecurrenceMonthly,
		alarm.RecurrenceYearly,
	}

	for _, rec := range recurrences {
		dir := ap.RecurringDir(rec)
		if err := fs.cancelAlarmInDir(dir, alarmID, rec); err == nil {
			return nil // Encontrada y eliminada
		}
	}

	return fmt.Errorf("alarm not found: %s", alarmID)
}

// cancelAlarmInDir busca y cancela una alarma en un directorio específico
func (fs *FilesystemStorage) cancelAlarmInDir(dir string, alarmID string, recurrence alarm.Recurrence) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var alarms []*alarm.Alarm
		if err := json.Unmarshal(data, &alarms); err != nil {
			continue
		}

		// Buscar la alarma por ID
		found := false
		newAlarms := []*alarm.Alarm{}
		for _, alm := range alarms {
			if alm.ID == alarmID {
				found = true
				// No agregar (eliminar)
			} else {
				newAlarms = append(newAlarms, alm)
			}
		}

		if found {
			filename := filepath.Base(file)

			// Si quedan alarmas, reescribir el archivo
			if len(newAlarms) > 0 {
				jsonData, err := json.MarshalIndent(newAlarms, "", "  ")
				if err != nil {
					return err
				}
				if err := os.WriteFile(file, jsonData, 0644); err != nil {
					return err
				}
			} else {
				// Si no quedan alarmas, eliminar el archivo
				userID := extractUserIDFromPath(file)
				if err := fs.DeleteAlarms(userID, recurrence, filename); err != nil {
					return err
				}
			}

			return nil // Éxito
		}
	}

	return fmt.Errorf("alarm not found in directory")
}

// MoveAlarmsToPast mueve un archivo de alarmas a la carpeta past/
func (fs *FilesystemStorage) MoveAlarmsToPast(userID string, recurrence alarm.Recurrence, filename string) error {
	ap := NewAlarmPaths(fs.dataDir, userID)

	var srcPath, dstPath string
	if recurrence == alarm.RecurrenceOnce {
		srcPath = ap.PendingFile(filename)
	} else {
		srcPath = ap.RecurringFile(recurrence, filename)
	}

	dstPath = ap.PastFile(recurrence, filename)

	// Asegurar que el directorio de destino existe
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio past: %w", err)
	}

	// Mover el archivo
	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("error moviendo alarma a past: %w", err)
	}

	return nil
}

// CopyRecurringAlarmExecution copia las alarmas recurrentes a past/ con timestamp de ejecución
// Esto permite rastrear cuándo se disparó cada alarma recurrente y evitar duplicados
func (fs *FilesystemStorage) CopyRecurringAlarmExecution(userID string, recurrence alarm.Recurrence, alarms []*alarm.Alarm, executedAt time.Time) error {
	if len(alarms) == 0 {
		return nil
	}

	ap := NewAlarmPaths(fs.dataDir, userID)

	// Generar nombre de archivo con timestamp de ejecución
	executionFilename := alarm.ExecutionFilename(executedAt)
	dstPath := ap.PastFile(recurrence, executionFilename)

	// Asegurar que el directorio de destino existe
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("error creating past directory: %w", err)
	}

	// Serializar alarmas
	jsonData, err := json.MarshalIndent(alarms, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing alarms: %w", err)
	}

	// Escribir archivo
	if err := os.WriteFile(dstPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing execution record: %w", err)
	}

	return nil
}

// WasRecurringAlarmExecuted verifica si una alarma recurrente ya fue ejecutada en un momento dado
func (fs *FilesystemStorage) WasRecurringAlarmExecuted(userID string, recurrence alarm.Recurrence, executedAt time.Time) (bool, error) {
	ap := NewAlarmPaths(fs.dataDir, userID)
	executionFilename := alarm.ExecutionFilename(executedAt)
	executionPath := ap.PastFile(recurrence, executionFilename)

	// Si el archivo existe, la alarma ya fue ejecutada
	if _, err := os.Stat(executionPath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("error checking execution record: %w", err)
	}
}

// extractUserIDFromPath extrae el userID de una ruta de archivo
// Ejemplo: /data/users/alice/alarms/pending/file.json -> alice
func extractUserIDFromPath(filePath string) string {
	parts := strings.Split(filePath, string(filepath.Separator))
	for i, part := range parts {
		if part == "users" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
