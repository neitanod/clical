package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/calendar"
	"github.com/spf13/cobra"
)

var (
	addDatetime string
	addTitle    string
	addDuration int
	addLocation string
	addNotes    string
	addTags     []string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new event to the calendar",
	Long: `Add a new event to the specified user's calendar.

Examples:
  clical add --user=12345 --datetime="2025-11-20 14:00" --title="Meeting" --duration=60
  clical add --user=12345 --datetime="2025-11-20 14:00" --title="Call" --duration=30 --location="Zoom"
  clical add --user=12345 --datetime="2025-11-21 09:00" --title="Stand-up" --duration=15 --tags=work,team`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate user ID
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// Validate title
		if addTitle == "" {
			return fmt.Errorf("--title is required")
		}

		// Parse datetime
		datetime, err := parseDateTime(addDatetime)
		if err != nil {
			return fmt.Errorf("error parsing --datetime: %w", err)
		}

		// Create entry
		entry := calendar.NewEntry(userID, addTitle, datetime, addDuration)
		entry.Location = addLocation
		entry.Notes = addNotes
		entry.Tags = addTags

		// Save
		if err := store.SaveEntry(userID, entry); err != nil {
			return fmt.Errorf("error saving entry: %w", err)
		}

		// Show confirmation
		fmt.Printf("✓ Event created successfully\n\n")
		fmt.Printf("ID:       %s\n", entry.ID)
		fmt.Printf("Title:    %s\n", entry.Title)
		fmt.Printf("Date:     %s\n", entry.DateTime.Format("2006-01-02 15:04"))
		fmt.Printf("Duration: %d minutes\n", entry.Duration)
		if entry.Location != "" {
			fmt.Printf("Location: %s\n", entry.Location)
		}
		if len(entry.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(entry.Tags, ", "))
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addDatetime, "datetime", "", "Event date and time (YYYY-MM-DD HH:MM)")
	addCmd.Flags().StringVar(&addTitle, "title", "", "Event title")
	addCmd.Flags().IntVar(&addDuration, "duration", 60, "Duration in minutes")
	addCmd.Flags().StringVar(&addLocation, "location", "", "Event location")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "Additional notes")
	addCmd.Flags().StringSliceVar(&addTags, "tags", []string{}, "Tags (comma-separated)")

	addCmd.MarkFlagRequired("datetime")
	addCmd.MarkFlagRequired("title")
}

// parseDateTime parses a date/time in "YYYY-MM-DD HH:MM" format
func parseDateTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	// Use local timezone for parsing
	loc := time.Now().Location()

	for _, format := range formats {
		t, err := time.ParseInLocation(format, s, loc)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid format (use YYYY-MM-DD HH:MM)")
}
