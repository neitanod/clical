package cli

import (
	"fmt"
	"strconv"
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

// parseDateTime parses a date/time in multiple formats:
// - Relative: "+5m" (5 minutes), "+2h" (2 hours), "+1d" (1 day)
// - Absolute: "YYYY-MM-DD HH:MM", "YYYY-MM-DDTHH:MM"
// - Keywords: "tomorrow HH:MM"
func parseDateTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	now := time.Now()
	loc := now.Location()

	// Try relative time formats: +5m, +2h, +1d
	if strings.HasPrefix(s, "+") {
		value := s[1:len(s)-1]
		unit := s[len(s)-1:]

		duration, err := strconv.Atoi(value)
		if err == nil {
			switch unit {
			case "m":
				return now.Add(time.Duration(duration) * time.Minute), nil
			case "h":
				return now.Add(time.Duration(duration) * time.Hour), nil
			case "d":
				return now.Add(time.Duration(duration) * 24 * time.Hour), nil
			}
		}
	}

	// Try "tomorrow HH:MM" format
	if strings.HasPrefix(strings.ToLower(s), "tomorrow ") {
		timeStr := strings.TrimPrefix(strings.ToLower(s), "tomorrow ")
		timeStr = strings.TrimSpace(timeStr)

		// Parse time part (HH:MM)
		timeParts := strings.Split(timeStr, ":")
		if len(timeParts) == 2 {
			hour, err1 := strconv.Atoi(timeParts[0])
			minute, err2 := strconv.Atoi(timeParts[1])
			if err1 == nil && err2 == nil {
				tomorrow := now.Add(24 * time.Hour)
				return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), hour, minute, 0, 0, loc), nil
			}
		}
	}

	// Try absolute formats
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		t, err := time.ParseInLocation(format, s, loc)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid format (use: YYYY-MM-DD HH:MM, +5m, +2h, tomorrow 10:00)")
}
