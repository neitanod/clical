package reporter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/sebasvalencia/clical/pkg/storage"
)

// DailyReport contiene el reporte diario completo
type DailyReport struct {
	Date           time.Time
	UserID         string
	Events         []*calendar.Entry
	Summary        Summary
	FreetimeBlocks []FreetimeBlock
	Tomorrow       []*calendar.Entry
	Suggestions    []string
}

// Summary contiene estadísticas del día
type Summary struct {
	TotalEvents   int
	TotalHours    float64
	FirstEvent    *time.Time
	LastEvent     *time.Time
	FreeHours     float64
	NextEvent     *calendar.Entry
	MinutesToNext int
}

// FreetimeBlock representa un bloque de tiempo libre
type FreetimeBlock struct {
	Start    time.Time
	End      time.Time
	Duration int // minutos
}

// GenerateDailyReport genera el reporte diario para un usuario
func GenerateDailyReport(store storage.Storage, userID string, date time.Time) (*DailyReport, error) {
	// Normalizar fecha a inicio del día
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	// Obtener eventos del día
	filter := calendar.NewFilter()
	filter.WithDateRange(dayStart, dayEnd)
	events, err := store.ListEntries(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo eventos: %w", err)
	}

	// Obtener eventos de mañana
	tomorrowStart := dayEnd
	tomorrowEnd := tomorrowStart.Add(24 * time.Hour)
	filterTomorrow := calendar.NewFilter()
	filterTomorrow.WithDateRange(tomorrowStart, tomorrowEnd)
	tomorrow, err := store.ListEntries(userID, filterTomorrow)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo eventos de mañana: %w", err)
	}

	// Calcular resumen
	summary := calculateSummary(events)

	// Calcular bloques libres
	freetime := calculateFreetime(events, dayStart, dayEnd)

	// Generar sugerencias
	suggestions := generateSuggestions(events, tomorrow)

	report := &DailyReport{
		Date:           dayStart,
		UserID:         userID,
		Events:         events,
		Summary:        summary,
		FreetimeBlocks: freetime,
		Tomorrow:       tomorrow,
		Suggestions:    suggestions,
	}

	return report, nil
}

// calculateSummary calcula estadísticas del día
func calculateSummary(events []*calendar.Entry) Summary {
	summary := Summary{
		TotalEvents: len(events),
	}

	if len(events) == 0 {
		return summary
	}

	// Calcular horas totales
	totalMinutes := 0
	for _, e := range events {
		totalMinutes += e.Duration
	}
	summary.TotalHours = float64(totalMinutes) / 60.0

	// Primer y último evento
	first := events[0].DateTime
	summary.FirstEvent = &first

	last := events[len(events)-1].EndTime()
	summary.LastEvent = &last

	// Calcular tiempo libre (asumiendo día laboral de 8am a 6pm)
	workDayMinutes := 10 * 60 // 10 horas
	summary.FreeHours = float64(workDayMinutes-totalMinutes) / 60.0

	// Próximo evento
	now := time.Now()
	for _, e := range events {
		if e.DateTime.After(now) {
			summary.NextEvent = e
			summary.MinutesToNext = int(e.DateTime.Sub(now).Minutes())
			break
		}
	}

	return summary
}

// calculateFreetime calcula bloques de tiempo libre
func calculateFreetime(events []*calendar.Entry, dayStart, dayEnd time.Time) []FreetimeBlock {
	var blocks []FreetimeBlock

	if len(events) == 0 {
		// Todo el día libre
		return []FreetimeBlock{
			{
				Start:    dayStart.Add(8 * time.Hour), // 8am
				End:      dayStart.Add(18 * time.Hour), // 6pm
				Duration: 10 * 60, // 10 horas
			},
		}
	}

	// Bloques entre eventos
	workStart := dayStart.Add(8 * time.Hour) // 8am
	workEnd := dayStart.Add(18 * time.Hour) // 6pm

	lastEnd := workStart

	for _, event := range events {
		// Si el evento está dentro del horario laboral
		if event.DateTime.After(workStart) && event.DateTime.Before(workEnd) {
			// Si hay gap desde último evento
			if event.DateTime.After(lastEnd) {
				gap := int(event.DateTime.Sub(lastEnd).Minutes())
				if gap >= 15 { // Solo bloques de al menos 15 min
					blocks = append(blocks, FreetimeBlock{
						Start:    lastEnd,
						End:      event.DateTime,
						Duration: gap,
					})
				}
			}
			lastEnd = event.EndTime()
		}
	}

	// Bloque final hasta fin de jornada
	if lastEnd.Before(workEnd) {
		gap := int(workEnd.Sub(lastEnd).Minutes())
		if gap >= 15 {
			blocks = append(blocks, FreetimeBlock{
				Start:    lastEnd,
				End:      workEnd,
				Duration: gap,
			})
		}
	}

	return blocks
}

// generateSuggestions genera sugerencias para el usuario
func generateSuggestions(today, tomorrow []*calendar.Entry) []string {
	var suggestions []string

	// Sugerencia si mañana hay muchos eventos
	if len(tomorrow) > 4 {
		suggestions = append(suggestions, fmt.Sprintf("⚠️ Mañana tienes %d eventos programados. Considera revisar preparación hoy.", len(tomorrow)))
	}

	// Sugerencia si hay eventos hoy con notas
	for _, e := range today {
		if e.Notes != "" && strings.Contains(strings.ToLower(e.Notes), "preparar") {
			suggestions = append(suggestions, fmt.Sprintf("Revisar preparación para: %s (%s)", e.Title, e.DateTime.Format("15:04")))
		}
	}

	return suggestions
}

// FormatDailyReport formatea el reporte diario como Markdown
func FormatDailyReport(report *DailyReport) string {
	var md strings.Builder

	// Fecha en español
	weekdays := []string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}
	months := []string{"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

	weekday := weekdays[report.Date.Weekday()]
	month := months[report.Date.Month()]

	md.WriteString(fmt.Sprintf("# Reporte Diario: %s %d de %s, %d\n\n",
		weekday, report.Date.Day(), month, report.Date.Year()))

	// Resumen
	md.WriteString("## Resumen del Día\n\n")
	md.WriteString(fmt.Sprintf("- **Eventos totales:** %d\n", report.Summary.TotalEvents))
	md.WriteString(fmt.Sprintf("- **Horas ocupadas:** %.1f horas\n", report.Summary.TotalHours))

	if report.Summary.FirstEvent != nil {
		md.WriteString(fmt.Sprintf("- **Primer evento:** %s\n", report.Summary.FirstEvent.Format("15:04")))
	}
	if report.Summary.LastEvent != nil {
		md.WriteString(fmt.Sprintf("- **Último evento:** %s\n", report.Summary.LastEvent.Format("15:04")))
	}
	md.WriteString(fmt.Sprintf("- **Tiempo libre:** %.1f horas\n", report.Summary.FreeHours))
	md.WriteString("\n")

	// Próximo evento
	if report.Summary.NextEvent != nil {
		md.WriteString(fmt.Sprintf("### 🔴 PRÓXIMO (en %d minutos)\n\n", report.Summary.MinutesToNext))
		md.WriteString(fmt.Sprintf("**[%s - %s] %s**\n",
			report.Summary.NextEvent.DateTime.Format("15:04"),
			report.Summary.NextEvent.EndTime().Format("15:04"),
			report.Summary.NextEvent.Title))
		md.WriteString(fmt.Sprintf("- ID: %s\n", report.Summary.NextEvent.ID))
		md.WriteString(fmt.Sprintf("- Duración: %d min\n", report.Summary.NextEvent.Duration))
		if report.Summary.NextEvent.Location != "" {
			md.WriteString(fmt.Sprintf("- Ubicación: %s\n", report.Summary.NextEvent.Location))
		}
		if len(report.Summary.NextEvent.Tags) > 0 {
			md.WriteString(fmt.Sprintf("- Tags: #%s\n", strings.Join(report.Summary.NextEvent.Tags, " #")))
		}
		if report.Summary.NextEvent.Notes != "" {
			md.WriteString(fmt.Sprintf("- Notas: %s\n", report.Summary.NextEvent.Notes))
		}
		md.WriteString("\n")
	}

	// Agenda del día
	if len(report.Events) > 0 {
		md.WriteString("## Agenda de Hoy\n\n")

		for _, event := range report.Events {
			md.WriteString(fmt.Sprintf("**[%s - %s] %s**\n",
				event.DateTime.Format("15:04"),
				event.EndTime().Format("15:04"),
				event.Title))

			md.WriteString(fmt.Sprintf("- ID: %s\n", event.ID))
			md.WriteString(fmt.Sprintf("- Duración: %d min\n", event.Duration))

			if event.Location != "" {
				md.WriteString(fmt.Sprintf("- Ubicación: %s\n", event.Location))
			}

			if len(event.Tags) > 0 {
				md.WriteString(fmt.Sprintf("- Tags: #%s\n", strings.Join(event.Tags, " #")))
			}

			if event.Notes != "" {
				md.WriteString(fmt.Sprintf("- Notas: %s\n", event.Notes))
			}

			md.WriteString("\n")
		}
	} else {
		md.WriteString("## Agenda de Hoy\n\n")
		md.WriteString("No hay eventos programados para hoy.\n\n")
	}

	// Bloques libres
	if len(report.FreetimeBlocks) > 0 {
		md.WriteString("## Bloques de Tiempo Libre\n\n")
		for _, block := range report.FreetimeBlocks {
			md.WriteString(fmt.Sprintf("- %s - %s (%d min) - ",
				block.Start.Format("15:04"),
				block.End.Format("15:04"),
				block.Duration))

			// Sugerencia de uso
			if block.Duration >= 120 {
				md.WriteString("Ideal para: trabajo profundo, reuniones largas\n")
			} else if block.Duration >= 60 {
				md.WriteString("Ideal para: reuniones, tareas importantes\n")
			} else {
				md.WriteString("Ideal para: llamadas cortas, breaks\n")
			}
		}
		md.WriteString("\n")
	}

	// Vista de mañana
	if len(report.Tomorrow) > 0 {
		md.WriteString(fmt.Sprintf("## Vista de Mañana (%s %d)\n\n",
			weekdays[(report.Date.Weekday()+1)%7], report.Date.Day()+1))

		for _, event := range report.Tomorrow {
			md.WriteString(fmt.Sprintf("- [%s] %s (%d min)\n",
				event.DateTime.Format("15:04"),
				event.Title,
				event.Duration))
		}

		// Advertencia si mañana está pesado
		totalMinutes := 0
		for _, e := range report.Tomorrow {
			totalMinutes += e.Duration
		}
		hours := float64(totalMinutes) / 60.0
		if hours > 4 {
			md.WriteString(fmt.Sprintf("\n**⚠️ Día pesado mañana: %.1f horas de eventos**\n", hours))
		}

		md.WriteString("\n")
	}

	// Sugerencias
	if len(report.Suggestions) > 0 {
		md.WriteString("## Sugerencias de la IA\n\n")
		for _, suggestion := range report.Suggestions {
			md.WriteString(fmt.Sprintf("- %s\n", suggestion))
		}
		md.WriteString("\n")
	}

	// Footer
	md.WriteString("---\n\n")
	md.WriteString(fmt.Sprintf("*Generado: %s*\n", time.Now().Format("2006-01-02 15:04")))

	return md.String()
}

// FindNextEvents encuentra los próximos N eventos a partir de ahora
func FindNextEvents(store storage.Storage, userID string, count int) ([]*calendar.Entry, error) {
	now := time.Now()
	future := now.Add(30 * 24 * time.Hour) // Próximos 30 días

	filter := calendar.NewFilter()
	filter.WithDateRange(now, future)
	filter.OnlyFuture = true

	events, err := store.ListEntries(userID, filter)
	if err != nil {
		return nil, err
	}

	// Ordenar por fecha
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateTime.Before(events[j].DateTime)
	})

	// Limitar a count eventos
	if len(events) > count {
		events = events[:count]
	}

	return events, nil
}
