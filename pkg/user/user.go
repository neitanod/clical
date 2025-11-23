package user

import (
	"fmt"
	"time"
)

// User representa un usuario del sistema
type User struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Timezone string    `json:"timezone"`
	Config  UserConfig `json:"config"`
	Created time.Time  `json:"created"`
}

// UserConfig contiene la configuración personalizada del usuario
type UserConfig struct {
	DefaultDuration int    `json:"default_duration"` // minutos
	DateFormat      string `json:"date_format"`
	TimeFormat      string `json:"time_format"`
	FirstDayOfWeek  int    `json:"first_day_of_week"` // 0=Domingo, 1=Lunes
}

// NewUser crea un nuevo usuario con configuración por defecto
func NewUser(id, name, timezone string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Timezone: timezone,
		Config:   DefaultConfig(),
		Created:  time.Now(),
	}
}

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() UserConfig {
	return UserConfig{
		DefaultDuration: 60,              // 1 hora
		DateFormat:      "2006-01-02",    // YYYY-MM-DD
		TimeFormat:      "15:04",         // HH:MM
		FirstDayOfWeek:  1,               // Lunes
	}
}

// Validate valida que el usuario tenga todos los campos requeridos
func (u *User) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("id es requerido")
	}
	if u.Name == "" {
		return fmt.Errorf("name es requerido")
	}
	if u.Timezone == "" {
		return fmt.Errorf("timezone es requerido")
	}

	// Validar timezone
	_, err := time.LoadLocation(u.Timezone)
	if err != nil {
		return fmt.Errorf("timezone inválido: %v", err)
	}

	// Validar config
	if u.Config.DefaultDuration <= 0 {
		return fmt.Errorf("default_duration debe ser mayor a 0")
	}

	return nil
}

// Location retorna la zona horaria del usuario
func (u *User) Location() (*time.Location, error) {
	return time.LoadLocation(u.Timezone)
}

// FormatDate formatea una fecha según la configuración del usuario
func (u *User) FormatDate(t time.Time) string {
	return t.Format(u.Config.DateFormat)
}

// FormatTime formatea una hora según la configuración del usuario
func (u *User) FormatTime(t time.Time) string {
	return t.Format(u.Config.TimeFormat)
}

// FormatDateTime formatea una fecha y hora según la configuración del usuario
func (u *User) FormatDateTime(t time.Time) string {
	return u.FormatDate(t) + " " + u.FormatTime(t)
}
