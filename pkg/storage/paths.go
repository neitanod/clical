package storage

import (
	"fmt"
	"path/filepath"
	"time"
)

// getEntryDir retorna el directorio donde se almacena una entrada
// Formato: data/users/{userID}/events/{year}/{month}/{day}/
func getEntryDir(dataDir, userID string, datetime time.Time) string {
	year := datetime.Format("2006")
	month := datetime.Format("01")
	day := datetime.Format("02")

	return filepath.Join(dataDir, "users", userID, "events", year, month, day)
}

// getEntryPath retorna la ruta completa de un archivo de entrada
func getEntryPath(dataDir, userID string, datetime time.Time, filename string, ext string) string {
	dir := getEntryDir(dataDir, userID, datetime)
	return filepath.Join(dir, filename+ext)
}

// getUserDir retorna el directorio de un usuario
func getUserDir(dataDir, userID string) string {
	return filepath.Join(dataDir, "users", userID)
}

// getUserPath retorna la ruta del archivo de usuario
func getUserPath(dataDir, userID string, ext string) string {
	dir := getUserDir(dataDir, userID)
	return filepath.Join(dir, "user"+ext)
}

// getStateDir retorna el directorio de estado de un usuario
func getStateDir(dataDir, userID string) string {
	return filepath.Join(dataDir, "users", userID, ".state")
}

// getStatePath retorna la ruta de un archivo de estado
func getStatePath(dataDir, userID, filename string) string {
	dir := getStateDir(dataDir, userID)
	return filepath.Join(dir, filename)
}

// getYearMonthDayFromDate retorna año, mes, día como strings
func getYearMonthDayFromDate(t time.Time) (string, string, string) {
	return t.Format("2006"), t.Format("01"), t.Format("02")
}

// parseFilename extrae información del nombre de archivo
// Formato esperado: HH-MM-titulo-slug.md
func parseFilename(filename string) (hour, minute, title string, err error) {
	// Remover extensión
	name := filename
	if len(name) > 3 {
		name = name[:len(name)-3] // Remover .md o .json
	}

	// Parsear HH-MM-titulo-slug
	if len(name) < 5 {
		return "", "", "", fmt.Errorf("nombre de archivo inválido: %s", filename)
	}

	hour = name[0:2]
	minute = name[3:5]

	if len(name) > 6 {
		title = name[6:] // Resto es el slug del título
	}

	return hour, minute, title, nil
}
