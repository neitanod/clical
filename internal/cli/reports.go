package cli

import (
	"fmt"
	"time"

	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/sebasvalencia/clical/pkg/reporter"
	"github.com/spf13/cobra"
)

var (
	dailyReportDate string
	tomorrowReportDate string
	upcomingHours   int
	upcomingCount   int
)

// daily-report command
var dailyReportCmd = &cobra.Command{
	Use:   "daily-report",
	Short: "Reporte diario completo del calendario",
	Long: `Genera un reporte diario completo optimizado para asistencia de IA.

Incluye:
- Resumen del día (eventos, horas ocupadas, tiempo libre)
- Próximo evento inmediato
- Agenda completa del día
- Bloques de tiempo libre
- Vista previa del día siguiente
- Sugerencias de organización

Ejemplos:
  clical daily-report --user=12345
  clical daily-report --user=12345 --date="2025-11-21"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Parsear fecha (default hoy)
		date := time.Now()
		if dailyReportDate != "" {
			parsed, err := time.Parse("2006-01-02", dailyReportDate)
			if err != nil {
				return fmt.Errorf("error parseando --date: %w", err)
			}
			date = parsed
		}

		// Generar reporte
		report, err := reporter.GenerateDailyReport(store, userID, date)
		if err != nil {
			return fmt.Errorf("error generando reporte: %w", err)
		}

		// Formatear y mostrar
		output := reporter.FormatDailyReport(report)
		fmt.Print(output)

		return nil
	},
}

// tomorrow-report command
var tomorrowReportCmd = &cobra.Command{
	Use:   "tomorrow-report",
	Short: "Reporte de eventos de mañana",
	Long: `Genera un reporte con vista previa del día siguiente.

Útil para ejecutar al final del día (ej: 8pm) para prepararse para mañana.

Ejemplos:
  clical tomorrow-report --user=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Mañana
		tomorrow := time.Now().Add(24 * time.Hour)

		// Generar reporte de mañana
		report, err := reporter.GenerateDailyReport(store, userID, tomorrow)
		if err != nil {
			return fmt.Errorf("error generando reporte: %w", err)
		}

		// Formatear y mostrar
		output := reporter.FormatDailyReport(report)
		fmt.Print(output)

		return nil
	},
}

// upcoming-report command
var upcomingReportCmd = &cobra.Command{
	Use:   "upcoming-report",
	Short: "Reporte de próximos eventos",
	Long: `Muestra los próximos eventos en las siguientes horas.

Útil para ejecutar periódicamente (ej: cada hora) para recordar eventos próximos.

Ejemplos:
  clical upcoming-report --user=12345 --hours=2
  clical upcoming-report --user=12345 --count=5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		var events []*calendar.Entry
		var err error

		if upcomingCount > 0 {
			// Próximos N eventos
			events, err = reporter.FindNextEvents(store, userID, upcomingCount)
		} else {
			// Próximos eventos en las siguientes X horas
			now := time.Now()
			future := now.Add(time.Duration(upcomingHours) * time.Hour)

			filter := calendar.NewFilter()
			filter.WithDateRange(now, future)
			filter.OnlyFuture = true

			events, err = store.ListEntries(userID, filter)
		}

		if err != nil {
			return fmt.Errorf("error obteniendo eventos: %w", err)
		}

		// Mostrar eventos
		if len(events) == 0 {
			fmt.Printf("No hay eventos próximos en las siguientes %d horas\n", upcomingHours)
			return nil
		}

		fmt.Printf("📅 Próximos eventos:\n\n")

		for _, event := range events {
			minutesUntil := int(event.DateTime.Sub(time.Now()).Minutes())

			if minutesUntil < 0 {
				continue // Skip eventos que ya empezaron
			}

			fmt.Printf("⏰ **En %d minutos** (%s)\n", minutesUntil, event.DateTime.Format("15:04"))
			fmt.Printf("   %s (%d min)\n", event.Title, event.Duration)
			fmt.Printf("   🆔 %s\n", event.ID)

			if event.Location != "" {
				fmt.Printf("   📍 %s\n", event.Location)
			}

			if event.Notes != "" {
				fmt.Printf("   📝 %s\n", event.Notes)
			}

			fmt.Println()
		}

		return nil
	},
}

// weekly-report command
var weeklyReportCmd = &cobra.Command{
	Use:   "weekly-report",
	Short: "Reporte semanal del calendario",
	Long: `Genera un reporte con vista de la semana.

Útil para ejecutar al inicio de semana (lunes) para planificar.

Ejemplos:
  clical weekly-report --user=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Inicio y fin de semana
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Domingo = 7
		}

		// Lunes de esta semana
		monday := now.AddDate(0, 0, -(weekday - 1))
		monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

		// Domingo de esta semana
		sunday := monday.AddDate(0, 0, 7)

		// Obtener eventos de la semana
		filter := calendar.NewFilter()
		filter.WithDateRange(monday, sunday)
		events, err := store.ListEntries(userID, filter)
		if err != nil {
			return fmt.Errorf("error obteniendo eventos: %w", err)
		}

		// Mostrar reporte
		fmt.Printf("# Reporte Semanal: %s al %s\n\n",
			monday.Format("2006-01-02"),
			sunday.Format("2006-01-02"))

		fmt.Printf("**Total de eventos:** %d\n\n", len(events))

		// Agrupar por día
		eventsByDay := make(map[string][]*calendar.Entry)
		for _, event := range events {
			day := event.DateTime.Format("2006-01-02")
			eventsByDay[day] = append(eventsByDay[day], event)
		}

		// Mostrar cada día
		weekdays := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}
		currentDay := monday

		for i := 0; i < 7; i++ {
			dayStr := currentDay.Format("2006-01-02")
			dayEvents := eventsByDay[dayStr]

			fmt.Printf("## %s %s\n\n", weekdays[i], currentDay.Format("02/01"))

			if len(dayEvents) == 0 {
				fmt.Printf("*Sin eventos*\n\n")
			} else {
				for _, event := range dayEvents {
					fmt.Printf("- [%s] %s (%d min) [ID: %s]\n",
						event.DateTime.Format("15:04"),
						event.Title,
						event.Duration,
						event.ID)
				}
				fmt.Println()
			}

			currentDay = currentDay.AddDate(0, 0, 1)
		}

		return nil
	},
}

func init() {
	// daily-report
	dailyReportCmd.Flags().StringVar(&dailyReportDate, "date", "", "Fecha para el reporte (YYYY-MM-DD, default: hoy)")

	// upcoming-report
	upcomingReportCmd.Flags().IntVar(&upcomingHours, "hours", 2, "Horas hacia adelante para buscar eventos")
	upcomingReportCmd.Flags().IntVar(&upcomingCount, "count", 0, "Mostrar próximos N eventos (sobrescribe --hours)")

	// Agregar a root
	rootCmd.AddCommand(dailyReportCmd)
	rootCmd.AddCommand(tomorrowReportCmd)
	rootCmd.AddCommand(upcomingReportCmd)
	rootCmd.AddCommand(weeklyReportCmd)
}
