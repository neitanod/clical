package alarm

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// Alarm representa una alarma
type Alarm struct {
	ID          string     `json:"id"`
	Context     string     `json:"context"`
	CreatedAt   time.Time  `json:"created_at"`
	Recurrence  Recurrence `json:"recurrence"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ScheduledFor time.Time `json:"scheduled_for,omitempty"` // Solo para output
	ExecutedAt  *time.Time `json:"executed_at,omitempty"`   // Solo para past alarms
}

// NewAlarm crea una nueva alarma
func NewAlarm(context string, recurrence Recurrence) *Alarm {
	return &Alarm{
		ID:         generateID(recurrence),
		Context:    context,
		CreatedAt:  time.Now(),
		Recurrence: recurrence,
	}
}

// generateID genera un ID único para la alarma
func generateID(recurrence Recurrence) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("alarm_%s_%d_%s", recurrence, timestamp, randomHex)
}

// Validate valida los campos de la alarma
func (a *Alarm) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("ID is required")
	}

	if a.Context == "" {
		return fmt.Errorf("context is required")
	}

	if len(a.Context) > 500 {
		return fmt.Errorf("context too long (max 500 characters)")
	}

	if !a.Recurrence.Valid() {
		return fmt.Errorf("invalid recurrence: %s", a.Recurrence)
	}

	if a.CreatedAt.IsZero() {
		return fmt.Errorf("created_at is required")
	}

	// ExpiresAt solo válido para alarmas recurrentes
	if a.ExpiresAt != nil && a.Recurrence == RecurrenceOnce {
		return fmt.Errorf("expires_at not allowed for one-time alarms")
	}

	// ExpiresAt debe ser futura
	if a.ExpiresAt != nil && a.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("expires_at must be in the future")
	}

	return nil
}

// IsExpired retorna true si la alarma ha expirado
func (a *Alarm) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// ShouldExecute retorna true si la alarma debería ejecutarse en el momento dado
func (a *Alarm) ShouldExecute(at time.Time) bool {
	// Si está expirada, no ejecutar
	if a.IsExpired() {
		return false
	}

	// Alarmas one-time se ejecutan una sola vez
	if a.Recurrence == RecurrenceOnce {
		return true
	}

	// Alarmas recurrentes siempre se ejecutan (si no están expiradas)
	return true
}

// WithScheduledFor establece el campo ScheduledFor (para output)
func (a *Alarm) WithScheduledFor(t time.Time) *Alarm {
	a.ScheduledFor = t
	return a
}

// WithExecutedAt establece el campo ExecutedAt (para historial)
func (a *Alarm) WithExecutedAt(t time.Time) *Alarm {
	a.ExecutedAt = &t
	return a
}

// Clone crea una copia de la alarma
func (a *Alarm) Clone() *Alarm {
	clone := &Alarm{
		ID:          a.ID,
		Context:     a.Context,
		CreatedAt:   a.CreatedAt,
		Recurrence:  a.Recurrence,
		ScheduledFor: a.ScheduledFor,
	}

	if a.ExpiresAt != nil {
		expiresAt := *a.ExpiresAt
		clone.ExpiresAt = &expiresAt
	}

	if a.ExecutedAt != nil {
		executedAt := *a.ExecutedAt
		clone.ExecutedAt = &executedAt
	}

	return clone
}
