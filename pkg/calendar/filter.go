package calendar

import (
	"strings"
	"time"
)

// Filter representa criterios de búsqueda/filtrado de entradas
type Filter struct {
	From         *time.Time
	To           *time.Time
	Query        string // búsqueda de texto libre
	Title        string // filtro por título
	Location     string // filtro por ubicación
	Tags         []string
	MinDuration  int
	MaxDuration  int
	HasLocation  bool
	OnlyPast     bool
	OnlyFuture   bool
	OnlyCurrent  bool
}

// NewFilter crea un filtro vacío
func NewFilter() *Filter {
	return &Filter{
		Tags: []string{},
	}
}

// Matches verifica si una entrada cumple con el filtro
func (f *Filter) Matches(entry *Entry) bool {
	// Filtro por rango de fechas
	if f.From != nil && entry.DateTime.Before(*f.From) {
		return false
	}
	if f.To != nil && entry.DateTime.After(*f.To) {
		return false
	}

	// Filtro por query (búsqueda en título y notas)
	if f.Query != "" {
		query := strings.ToLower(f.Query)
		title := strings.ToLower(entry.Title)
		notes := strings.ToLower(entry.Notes)
		if !strings.Contains(title, query) && !strings.Contains(notes, query) {
			return false
		}
	}

	// Filtro por título
	if f.Title != "" {
		if !strings.Contains(strings.ToLower(entry.Title), strings.ToLower(f.Title)) {
			return false
		}
	}

	// Filtro por ubicación
	if f.Location != "" {
		if !strings.Contains(strings.ToLower(entry.Location), strings.ToLower(f.Location)) {
			return false
		}
	}

	// Filtro por tags
	if len(f.Tags) > 0 {
		hasAllTags := true
		for _, tag := range f.Tags {
			if !entry.HasTag(tag) {
				hasAllTags = false
				break
			}
		}
		if !hasAllTags {
			return false
		}
	}

	// Filtro por duración
	if f.MinDuration > 0 && entry.Duration < f.MinDuration {
		return false
	}
	if f.MaxDuration > 0 && entry.Duration > f.MaxDuration {
		return false
	}

	// Filtro por presencia de ubicación
	if f.HasLocation && entry.Location == "" {
		return false
	}

	// Filtros temporales
	if f.OnlyPast && !entry.IsPast() {
		return false
	}
	if f.OnlyFuture && !entry.IsFuture() {
		return false
	}
	if f.OnlyCurrent && !entry.IsCurrent() {
		return false
	}

	return true
}

// WithDateRange establece el rango de fechas
func (f *Filter) WithDateRange(from, to time.Time) *Filter {
	f.From = &from
	f.To = &to
	return f
}

// WithQuery establece la búsqueda de texto
func (f *Filter) WithQuery(query string) *Filter {
	f.Query = query
	return f
}

// WithTags establece los tags a filtrar
func (f *Filter) WithTags(tags ...string) *Filter {
	f.Tags = tags
	return f
}
