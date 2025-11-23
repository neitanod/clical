package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/sebasvalencia/clical/pkg/user"
)

// FilesystemStorage implementa Storage usando sistema de archivos
type FilesystemStorage struct {
	dataDir string
}

// NewFilesystemStorage crea un nuevo storage en filesystem
func NewFilesystemStorage(dataDir string) (*FilesystemStorage, error) {
	// Crear directorio si no existe
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio: %w", err)
	}

	return &FilesystemStorage{
		dataDir: dataDir,
	}, nil
}

// SaveEntry guarda una entrada en filesystem
func (fs *FilesystemStorage) SaveEntry(userID string, entry *calendar.Entry) error {
	if err := entry.Validate(); err != nil {
		return fmt.Errorf("entrada inválida: %w", err)
	}

	// Crear directorio para la fecha
	dir := getEntryDir(fs.dataDir, userID, entry.DateTime)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	filename := entry.GenerateFilename()

	// Guardar Markdown
	mdPath := getEntryPath(fs.dataDir, userID, entry.DateTime, filename, ".md")
	mdContent := entryToMarkdown(entry)
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		return fmt.Errorf("error escribiendo markdown: %w", err)
	}

	// Guardar JSON
	jsonPath := getEntryPath(fs.dataDir, userID, entry.DateTime, filename, ".json")
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo JSON: %w", err)
	}

	return nil
}

// GetEntry obtiene una entrada por ID
func (fs *FilesystemStorage) GetEntry(userID, entryID string) (*calendar.Entry, error) {
	// Buscar en todos los directorios de eventos
	entries, err := fs.ListEntries(userID, calendar.NewFilter())
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.ID == entryID {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("entrada no encontrada: %s", entryID)
}

// ListEntries lista todas las entradas de un usuario con filtro opcional
func (fs *FilesystemStorage) ListEntries(userID string, filter *calendar.Filter) ([]*calendar.Entry, error) {
	if filter == nil {
		filter = calendar.NewFilter()
	}

	var entries []*calendar.Entry

	eventsDir := filepath.Join(fs.dataDir, "users", userID, "events")

	// Verificar que el directorio existe
	if _, err := os.Stat(eventsDir); os.IsNotExist(err) {
		return entries, nil // Retornar lista vacía si no existe
	}

	// Recorrer todos los archivos JSON
	err := filepath.Walk(eventsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Solo procesar archivos .json
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Leer y parsear JSON
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %w", path, err)
		}

		var entry calendar.Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			return fmt.Errorf("error parseando %s: %w", path, err)
		}

		// Aplicar filtro
		if filter.Matches(&entry) {
			entries = append(entries, &entry)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Ordenar por fecha
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].DateTime.Before(entries[j].DateTime)
	})

	return entries, nil
}

// DeleteEntry elimina una entrada
func (fs *FilesystemStorage) DeleteEntry(userID, entryID string) error {
	// Primero obtener la entrada para saber su ubicación
	entry, err := fs.GetEntry(userID, entryID)
	if err != nil {
		return err
	}

	filename := entry.GenerateFilename()

	// Eliminar archivos .md y .json
	mdPath := getEntryPath(fs.dataDir, userID, entry.DateTime, filename, ".md")
	jsonPath := getEntryPath(fs.dataDir, userID, entry.DateTime, filename, ".json")

	if err := os.Remove(mdPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error eliminando markdown: %w", err)
	}

	if err := os.Remove(jsonPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error eliminando JSON: %w", err)
	}

	return nil
}

// UpdateEntry actualiza una entrada existente
func (fs *FilesystemStorage) UpdateEntry(userID string, entry *calendar.Entry) error {
	// Eliminar la entrada antigua
	if err := fs.DeleteEntry(userID, entry.ID); err != nil {
		return err
	}

	// Actualizar timestamp
	entry.UpdatedAt = time.Now()

	// Guardar la nueva versión
	return fs.SaveEntry(userID, entry)
}

// SaveUser guarda un usuario
func (fs *FilesystemStorage) SaveUser(u *user.User) error {
	if err := u.Validate(); err != nil {
		return fmt.Errorf("usuario inválido: %w", err)
	}

	userDir := getUserDir(fs.dataDir, u.ID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	// Guardar Markdown
	mdPath := getUserPath(fs.dataDir, u.ID, ".md")
	mdContent := userToMarkdown(u)
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		return fmt.Errorf("error escribiendo markdown: %w", err)
	}

	// Guardar JSON
	jsonPath := getUserPath(fs.dataDir, u.ID, ".json")
	jsonData, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo JSON: %w", err)
	}

	return nil
}

// GetUser obtiene un usuario por ID
func (fs *FilesystemStorage) GetUser(userID string) (*user.User, error) {
	jsonPath := getUserPath(fs.dataDir, userID, ".json")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("usuario no encontrado: %s", userID)
		}
		return nil, fmt.Errorf("error leyendo usuario: %w", err)
	}

	var u user.User
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, fmt.Errorf("error parseando usuario: %w", err)
	}

	return &u, nil
}

// ListUsers lista todos los usuarios
func (fs *FilesystemStorage) ListUsers() ([]*user.User, error) {
	var users []*user.User

	usersDir := filepath.Join(fs.dataDir, "users")
	if _, err := os.Stat(usersDir); os.IsNotExist(err) {
		return users, nil
	}

	entries, err := os.ReadDir(usersDir)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		userID := entry.Name()
		u, err := fs.GetUser(userID)
		if err != nil {
			continue // Skip usuarios con errores
		}

		users = append(users, u)
	}

	return users, nil
}

// DeleteUser elimina un usuario y todos sus datos
func (fs *FilesystemStorage) DeleteUser(userID string) error {
	userDir := getUserDir(fs.dataDir, userID)
	return os.RemoveAll(userDir)
}

// GetReportState obtiene el estado de reportes
func (fs *FilesystemStorage) GetReportState(userID string) (*ReportState, error) {
	statePath := getStatePath(fs.dataDir, userID, "report-state.json")

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewReportState(), nil
		}
		return nil, fmt.Errorf("error leyendo estado: %w", err)
	}

	var state ReportState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("error parseando estado: %w", err)
	}

	if state.ReportedEvents == nil {
		state.ReportedEvents = make(map[string]string)
	}

	return &state, nil
}

// SaveReportState guarda el estado de reportes
func (fs *FilesystemStorage) SaveReportState(userID string, state *ReportState) error {
	stateDir := getStateDir(fs.dataDir, userID)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	statePath := getStatePath(fs.dataDir, userID, "report-state.json")
	jsonData, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando estado: %w", err)
	}

	if err := os.WriteFile(statePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo estado: %w", err)
	}

	return nil
}
