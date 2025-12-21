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
	Short: "Complete daily calendar report",
	Long: `Generate complete daily report optimized for AI assistance.

Includes:
- Day summary (events, busy hours, free time)
- Next immediate event
- Complete day agenda
- Free time blocks
- Next day preview
- Organization suggestions

Examples:
  clical daily-report --user=12345
  clical daily-report --user=12345 --date="2025-11-21"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// Parse date (default today)
		date := time.Now()
		if dailyReportDate != "" {
			parsed, err := time.Parse("2006-01-02", dailyReportDate)
			if err != nil {
				return fmt.Errorf("error parsing --date: %w", err)
			}
			date = parsed
		}

		// Generar reporte
		report, err := reporter.GenerateDailyReport(store, userID, date)
		if err != nil {
			return fmt.Errorf("error generating reporte: %w", err)
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
	Short: "Tomorrow's events report",
	Long: `Generate report with next day preview.

Useful to run at end of day (eg: 8pm) to prepare for tomorrow.

Examples:
  clical tomorrow-report --user=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// Tomorrow
		tomorrow := time.Now().Add(24 * time.Hour)

		// Generate tomorrow report
		report, err := reporter.GenerateDailyReport(store, userID, tomorrow)
		if err != nil {
			return fmt.Errorf("error generating reporte: %w", err)
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
	Short: "Upcoming events report",
	Long: `Show upcoming events in the next hours.

Useful to run periodically (eg: hourly) to remind of upcoming events.

Examples:
  clical upcoming-report --user=12345 --hours=2
  clical upcoming-report --user=12345 --count=5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
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
			return fmt.Errorf("error getting events: %w", err)
		}

		// Show events
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

			fmt.Printf("⏰ **In %d minutes** (%s)\n", minutesUntil, event.DateTime.Format("15:04"))
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
	Short: "Weekly calendar report",
	Long: `Generate report with week view.

Useful to run at start of week (Monday) for planning.

Examples:
  clical weekly-report --user=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// Start and end of week
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}

		// Monday of this week
		monday := now.AddDate(0, 0, -(weekday - 1))
		monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

		// Sunday of this week
		sunday := monday.AddDate(0, 0, 7)

		// Get week events
		filter := calendar.NewFilter()
		filter.WithDateRange(monday, sunday)
		events, err := store.ListEntries(userID, filter)
		if err != nil {
			return fmt.Errorf("error getting events: %w", err)
		}

		// Show report
		fmt.Printf("# Reporte Semanal: %s al %s\n\n",
			monday.Format("2006-01-02"),
			sunday.Format("2006-01-02"))

		fmt.Printf("**Total de eventos:** %d\n\n", len(events))

		// Group by day
		eventsByDay := make(map[string][]*calendar.Entry)
		for _, event := range events {
			day := event.DateTime.Format("2006-01-02")
			eventsByDay[day] = append(eventsByDay[day], event)
		}

		// Show each day
		weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		currentDay := monday

		for i := 0; i < 7; i++ {
			dayStr := currentDay.Format("2006-01-02")
			dayEvents := eventsByDay[dayStr]

			fmt.Printf("## %s %s\n\n", weekdays[i], currentDay.Format("02/01"))

			if len(dayEvents) == 0 {
				fmt.Printf("*No events*\n\n")
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

// yesterday-today-report command
var yesterdayTodayReportCmd = &cobra.Command{
	Use:   "yesterday-today-report",
	Short: "Yesterday and today events report",
	Long: `Generate combined report with yesterday and today events (48 hours).

Useful to get complete context of last 48 hours.

Examples:
  clical yesterday-today-report --user=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		yesterday := today.Add(-24 * time.Hour)
		tomorrow := today.Add(24 * time.Hour)

		// Create filter for yesterday + today
		filter := calendar.NewFilter()
		filter.From = &yesterday
		filter.To = &tomorrow

		// Get events
		entries, err := store.ListEntries(userID, filter)
		if err != nil {
			return fmt.Errorf("error listing eventos: %w", err)
		}

		// Generar reporte
		fmt.Println("# REPORT: YESTERDAY + TODAY")
		fmt.Println()
		fmt.Printf("Period: %s - %s\n", yesterday.Format("2006-01-02"), today.Format("2006-01-02"))
		fmt.Println()

		if len(entries) == 0 {
			fmt.Println("No events in this period")
			return nil
		}

		// Separate events by day
		var yesterdayEvents []*calendar.Entry
		var todayEvents []*calendar.Entry

		for _, entry := range entries {
			// Normalize to same location for comparison
			entryDate := time.Date(entry.DateTime.Year(), entry.DateTime.Month(), entry.DateTime.Day(), 0, 0, 0, 0, now.Location())
			if entryDate.Equal(yesterday) {
				yesterdayEvents = append(yesterdayEvents, entry)
			} else if entryDate.Equal(today) {
				todayEvents = append(todayEvents, entry)
			}
		}

		// Show events de ayer
		fmt.Printf("## YESTERDAY (%s)\n\n", yesterday.Format("Monday, 02 Jan 2006"))
		if len(yesterdayEvents) == 0 {
			fmt.Println("No events")
		} else {
			for _, event := range yesterdayEvents {
				fmt.Printf("- [%s] %s", event.DateTime.Format("15:04"), event.Title)
				if event.Duration > 0 {
					fmt.Printf(" (%d min)", event.Duration)
				}
				fmt.Printf(" [ID: %s]", event.ID)
				if event.Location != "" {
					fmt.Printf(" - %s", event.Location)
				}
				fmt.Println()
			}
		}
		fmt.Println()

		// Show events de hoy
		fmt.Printf("## TODAY (%s)\n\n", today.Format("Monday, 02 Jan 2006"))
		if len(todayEvents) == 0 {
			fmt.Println("No events")
		} else {
			for _, event := range todayEvents {
				fmt.Printf("- [%s] %s", event.DateTime.Format("15:04"), event.Title)
				if event.Duration > 0 {
					fmt.Printf(" (%d min)", event.Duration)
				}
				fmt.Printf(" [ID: %s]", event.ID)
				if event.Location != "" {
					fmt.Printf(" - %s", event.Location)
				}
				fmt.Println()
			}
		}
		fmt.Println()

		// Resumen
		fmt.Println("## SUMMARY")
		fmt.Printf("- Total events: %d\n", len(entries))
		fmt.Printf("- Yesterday events: %d\n", len(yesterdayEvents))
		fmt.Printf("- Today events: %d\n", len(todayEvents))

		return nil
	},
}

func init() {
	// daily-report
	dailyReportCmd.Flags().StringVar(&dailyReportDate, "date", "", "Report date (YYYY-MM-DD, default: today)")

	// upcoming-report
	upcomingReportCmd.Flags().IntVar(&upcomingHours, "hours", 2, "Hours ahead to search for events")
	upcomingReportCmd.Flags().IntVar(&upcomingCount, "count", 0, "Show next N events (overrides --hours)")

	// Agregar a root
	rootCmd.AddCommand(dailyReportCmd)
	rootCmd.AddCommand(tomorrowReportCmd)
	rootCmd.AddCommand(upcomingReportCmd)
	rootCmd.AddCommand(weeklyReportCmd)
	rootCmd.AddCommand(yesterdayTodayReportCmd)
}
