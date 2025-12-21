package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/spf13/cobra"
)

var (
	listFrom  string
	listTo    string
	listRange string
	listTags  []string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List calendar events",
	Long: `List calendar events with optional filters.

Available ranges:
  today          - Today
  yesterday      - Yesterday
  48             - Yesterday + today (48 hours)
  week           - Next 7 days
  past-week      - Last 7 days
  month          - Until end of month
  past-month     - Last 30 days
  month-to-date  - From start of month until today
  year-to-date   - From start of year until today

Examples:
  clical list --user=12345
  clical list --user=12345 --from="2025-11-20"
  clical list --user=12345 --range=today
  clical list --user=12345 --range=yesterday
  clical list --user=12345 --range=48
  clical list --user=12345 --range=month-to-date
  clical list --user=12345 --tags=trabajo,reunion`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate user ID
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// Build filter
		filter := calendar.NewFilter()

		// Apply predefined range
		if listRange != "" {
			from, to, err := parseRange(listRange)
			if err != nil {
				return err
			}
			filter.From = &from
			filter.To = &to
		}

		// Apply manual from/to (overrides range)
		if listFrom != "" {
			from, err := time.Parse("2006-01-02", listFrom)
			if err != nil {
				return fmt.Errorf("error parsing --from: %w", err)
			}
			filter.From = &from
		}

		if listTo != "" {
			to, err := time.Parse("2006-01-02", listTo)
			if err != nil {
				return fmt.Errorf("error parsing --to: %w", err)
			}
			// Include entire final day
			to = to.Add(24 * time.Hour)
			filter.To = &to
		}

		// Apply tags
		if len(listTags) > 0 {
			filter.Tags = listTags
		}

		// Get events
		entries, err := store.ListEntries(userID, filter)
		if err != nil {
			return fmt.Errorf("error listing eventos: %w", err)
		}

		// Show results
		if len(entries) == 0 {
			fmt.Println("No events found")
			return nil
		}

		fmt.Printf("Found %d event(s)\n\n", len(entries))

		for _, entry := range entries {
			printEntryRow(entry)
			fmt.Println()
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listFrom, "from", "", "Start date (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listTo, "to", "", "End date (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listRange, "range", "", "Predefined range: today, yesterday, 48, week, past-week, month, past-month, month-to-date, year-to-date")
	listCmd.Flags().StringSliceVar(&listTags, "tags", []string{}, "Filter by tags")
}

// parseRange parsea rangos predefinidos como "today", "week", "month"
func parseRange(r string) (time.Time, time.Time, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.Add(-24 * time.Hour)

	switch r {
	case "today":
		return today, today.Add(24 * time.Hour), nil

	case "yesterday":
		return yesterday, today, nil

	case "48":
		// Yesterday + today (48 hours)
		return yesterday, today.Add(24 * time.Hour), nil

	case "week":
		// From today for next 7 days
		return today, today.Add(7 * 24 * time.Hour), nil

	case "past-week":
		// Last 7 days (including today)
		return today.Add(-7 * 24 * time.Hour), today.Add(24 * time.Hour), nil

	case "month":
		// From today until end of month
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
		return today, endOfMonth, nil

	case "past-month":
		// Last 30 days (including today)
		return today.Add(-30 * 24 * time.Hour), today.Add(24 * time.Hour), nil

	case "month-to-date":
		// From start of current month until today
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return startOfMonth, today.Add(24 * time.Hour), nil

	case "year-to-date":
		// From start of year until today
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		return startOfYear, today.Add(24 * time.Hour), nil

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid range: %s (use: today, yesterday, 48, week, past-week, month, past-month, month-to-date, year-to-date)", r)
	}
}

// printEntryRow imprime una entrada en formato compacto
func printEntryRow(entry *calendar.Entry) {
	fmt.Printf("[%s] %s",
		entry.DateTime.Format("2006-01-02 15:04"),
		entry.Title,
	)

	if entry.Duration > 0 {
		fmt.Printf(" (%d min)", entry.Duration)
	}

	fmt.Printf(" [ID: %s]", entry.ID)

	if entry.Location != "" {
		fmt.Printf(" - %s", entry.Location)
	}

	if len(entry.Tags) > 0 {
		fmt.Printf(" #%s", strings.Join(entry.Tags, " #"))
	}
}
