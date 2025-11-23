package calendar

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Entry representa un evento en el calendario
type Entry struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	DateTime  time.Time         `json:"datetime"`
	Title     string            `json:"title"`
	Duration  int               `json:"duration"` // minutos
	Location  string            `json:"location,omitempty"`
	Notes     string            `json:"notes,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// NewEntry crea una nueva entrada con valores por defecto
func NewEntry(userID, title string, datetime time.Time, duration int) *Entry {
	now := time.Now()
	return &Entry{
		ID:        GenerateID(),
		UserID:    userID,
		DateTime:  datetime,
		Title:     title,
		Duration:  duration,
		Tags:      []string{},
		Metadata:  make(map[string]string),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GenerateID genera un ID único para una entrada
func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Validate valida que la entrada tenga todos los campos requeridos
func (e *Entry) Validate() error {
	if e.UserID == "" {
		return fmt.Errorf("user_id es requerido")
	}
	if e.Title == "" {
		return fmt.Errorf("title es requerido")
	}
	if e.DateTime.IsZero() {
		return fmt.Errorf("datetime es requerido")
	}
	if e.Duration <= 0 {
		return fmt.Errorf("duration debe ser mayor a 0")
	}
	return nil
}

// EndTime calcula la hora de finalización del evento
func (e *Entry) EndTime() time.Time {
	return e.DateTime.Add(time.Duration(e.Duration) * time.Minute)
}

// IsPast verifica si el evento ya pasó
func (e *Entry) IsPast() bool {
	return time.Now().After(e.EndTime())
}

// IsFuture verifica si el evento es futuro
func (e *Entry) IsFuture() bool {
	return time.Now().Before(e.DateTime)
}

// IsCurrent verifica si el evento está en curso
func (e *Entry) IsCurrent() bool {
	now := time.Now()
	return now.After(e.DateTime) && now.Before(e.EndTime())
}

// GenerateFilename genera el nombre de archivo para esta entrada
// Formato: HH-MM-titulo-slug.md
func (e *Entry) GenerateFilename() string {
	timeStr := e.DateTime.Format("15-04")
	titleSlug := slugify(e.Title)
	return fmt.Sprintf("%s-%s", timeStr, titleSlug)
}

// slugify convierte un título en un slug válido para nombre de archivo
func slugify(s string) string {
	// Convertir a minúsculas
	s = strings.ToLower(s)

	// Reemplazar espacios y caracteres especiales
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "á", "a")
	s = strings.ReplaceAll(s, "é", "e")
	s = strings.ReplaceAll(s, "í", "i")
	s = strings.ReplaceAll(s, "ó", "o")
	s = strings.ReplaceAll(s, "ú", "u")
	s = strings.ReplaceAll(s, "ñ", "n")

	// Remover caracteres no alfanuméricos excepto guiones
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug := result.String()

	// Limitar longitud
	if len(slug) > 50 {
		slug = slug[:50]
	}

	// Remover guiones al inicio/fin
	slug = strings.Trim(slug, "-")

	return slug
}

// HasTag verifica si la entrada tiene un tag específico
func (e *Entry) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag agrega un tag a la entrada (evita duplicados)
func (e *Entry) AddTag(tag string) {
	if !e.HasTag(tag) {
		e.Tags = append(e.Tags, tag)
	}
}

// RemoveTag elimina un tag de la entrada
func (e *Entry) RemoveTag(tag string) {
	for i, t := range e.Tags {
		if t == tag {
			e.Tags = append(e.Tags[:i], e.Tags[i+1:]...)
			return
		}
	}
}
