package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/alarm"
	"github.com/spf13/cobra"
)

var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Manage alarms",
	Long:  `Manage alarms for reminders and scheduled follow-ups.`,
}

// alarm-add variables
var (
	alarmContext string
	alarmAt      string
	alarmDaily   string
	alarmWeekly  string
	alarmMonthly string
	alarmYearly  string
	alarmExpires string
)

var alarmAddCmd = &cobra.Command{
	Use:          "add",
	Short:        "Add a new alarm",
	SilenceUsage: true,
	Long: `Add a new alarm (one-time or recurring).

Examples:
  # One-time
  clical alarm add --user alice --at "2025-11-23 14:30" --context "Revisar PR"
  clical alarm add --user alice --at "tomorrow 10:00" --context "Llamar cliente"
  clical alarm add --user alice --at "+30m" --context "Verificar deploy"

  # Recurrente daily
  clical alarm add --user alice --daily "14:30" --context "Revisar métricas"
  clical alarm add --user alice --daily "09:00" --expires "2025-12-31" --context "Stand-up temporal"

  # Recurrente weekly
  clical alarm add --user alice --weekly "monday 14:30" --context "Reunión semanal"
  clical alarm add --user alice --weekly "friday 17:00" --context "Reporte semanal"

  # Recurrente monthly
  clical alarm add --user alice --monthly "1 09:00" --context "Reporte mensual"
  clical alarm add --user alice --monthly "15 14:30" --context "Revisión quincenal"

  # Recurrente yearly
  clical alarm add --user alice --yearly "01-01 00:00" --context "Feliz año nuevo"
  clical alarm add --user alice --yearly "11-21 10:00" --context "Aniversario del proyecto"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		if alarmContext == "" {
			return fmt.Errorf("--context is required")
		}

		// Determinar tipo de alarma
		if alarmAt != "" {
			return addOneTimeAlarm(userID, alarmAt, alarmContext)
		} else if alarmDaily != "" {
			return addDailyAlarm(userID, alarmDaily, alarmContext, alarmExpires)
		} else if alarmWeekly != "" {
			return addWeeklyAlarm(userID, alarmWeekly, alarmContext, alarmExpires)
		} else if alarmMonthly != "" {
			return addMonthlyAlarm(userID, alarmMonthly, alarmContext, alarmExpires)
		} else if alarmYearly != "" {
			return addYearlyAlarm(userID, alarmYearly, alarmContext, alarmExpires)
		}

		return fmt.Errorf("must specify --at, --daily, --weekly, --monthly or --yearly")
	},
}

func addOneTimeAlarm(userID, atStr, context string) error {
	// Parse date/time
	alarmTime, err := parseDateTime(atStr)
	if err != nil {
		return fmt.Errorf("error parsing --at: %w", err)
	}

	// Verify it's in the future
	if alarmTime.Before(time.Now()) {
		printRedError("date/time must be in the future")
		os.Exit(1)
	}

	// Round to minute
	alarmTime = alarm.RoundToMinute(alarmTime)

	// Create alarm
	alm := alarm.NewAlarm(context, alarm.RecurrenceOnce)

	// Save
	filename := alarm.OneTimeFilename(alarmTime)
	if err := store.SaveAlarm(userID, alarmTime, alarm.RecurrenceOnce, filename, alm); err != nil {
		return fmt.Errorf("error saving alarm: %w", err)
	}

	fmt.Printf("✓ Alarm created successfully\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Type:       one-time\n")
	fmt.Printf("Scheduled: %s\n", alarmTime.Format("2006-01-02 15:04"))
	fmt.Printf("Context:   %s\n", context)

	return nil
}

func addDailyAlarm(userID, timeStr, context, expiresStr string) error {
	// Parse time (HH:MM)
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format for --daily (must be HH:MM)")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("invalid hour: %s", parts[0])
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid minute: %s", parts[1])
	}

	// Create alarm
	alm := alarm.NewAlarm(context, alarm.RecurrenceDaily)

	// Add expiration if specified
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parsing --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Save
	schedule := alarm.DailySchedule{Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceDaily, filename, alm); err != nil {
		return fmt.Errorf("error saving alarm: %w", err)
	}

	fmt.Printf("✓ Alarm created successfully\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Type:       daily\n")
	fmt.Printf("Time:       %02d:%02d\n", hour, minute)
	fmt.Printf("Context:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expires:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

func addWeeklyAlarm(userID, scheduleStr, context, expiresStr string) error {
	// Parse "monday 14:30"
	parts := strings.Fields(scheduleStr)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format para --weekly (debe ser DAYNAME HH:MM)")
	}

	weekday, err := alarm.ParseWeekday(parts[0])
	if err != nil {
		return err
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) != 2 {
		return fmt.Errorf("invalid format para hora (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("invalid hour: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid minute: %s", timeParts[1])
	}

	// Create alarm
	alm := alarm.NewAlarm(context, alarm.RecurrenceWeekly)

	// Add expiration if specified
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parsing --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Save
	schedule := alarm.WeeklySchedule{Weekday: weekday, Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceWeekly, filename, alm); err != nil {
		return fmt.Errorf("error saving alarm: %w", err)
	}

	fmt.Printf("✓ Alarm created successfully\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Type:       weekly\n")
	fmt.Printf("Day:        %s\n", weekday.String())
	fmt.Printf("Time:       %02d:%02d\n", hour, minute)
	fmt.Printf("Context:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expires:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

func addMonthlyAlarm(userID, scheduleStr, context, expiresStr string) error {
	// Parsear "15 14:30"
	parts := strings.Fields(scheduleStr)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format para --monthly (debe ser DAY HH:MM)")
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil || day < 1 || day > 31 {
		return fmt.Errorf("día inválido: %s (debe ser 1-31)", parts[0])
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) != 2 {
		return fmt.Errorf("invalid format para hora (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("invalid hour: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid minute: %s", timeParts[1])
	}

	// Create alarm
	alm := alarm.NewAlarm(context, alarm.RecurrenceMonthly)

	// Add expiration if specified
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parsing --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Save
	schedule := alarm.MonthlySchedule{Day: day, Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceMonthly, filename, alm); err != nil {
		return fmt.Errorf("error saving alarm: %w", err)
	}

	fmt.Printf("✓ Alarm created successfully\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Type:       monthly\n")
	fmt.Printf("Day:        %d\n", day)
	fmt.Printf("Time:       %02d:%02d\n", hour, minute)
	fmt.Printf("Context:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expires:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

func addYearlyAlarm(userID, scheduleStr, context, expiresStr string) error {
	// Parsear "11-21 14:30"
	parts := strings.Fields(scheduleStr)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format para --yearly (debe ser MM-DD HH:MM)")
	}

	dateParts := strings.Split(parts[0], "-")
	if len(dateParts) != 2 {
		return fmt.Errorf("invalid format para fecha (debe ser MM-DD)")
	}

	month, err := strconv.Atoi(dateParts[0])
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("mes inválido: %s (debe ser 1-12)", dateParts[0])
	}

	day, err := strconv.Atoi(dateParts[1])
	if err != nil || day < 1 || day > 31 {
		return fmt.Errorf("día inválido: %s (debe ser 1-31)", dateParts[1])
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) != 2 {
		return fmt.Errorf("invalid format para hora (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("invalid hour: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid minute: %s", timeParts[1])
	}

	// Create alarm
	alm := alarm.NewAlarm(context, alarm.RecurrenceYearly)

	// Add expiration if specified
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parsing --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Save
	schedule := alarm.YearlySchedule{Month: time.Month(month), Day: day, Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceYearly, filename, alm); err != nil {
		return fmt.Errorf("error saving alarm: %w", err)
	}

	fmt.Printf("✓ Alarm created successfully\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Type:       yearly\n")
	fmt.Printf("Date:      %02d-%02d\n", month, day)
	fmt.Printf("Time:       %02d:%02d\n", hour, minute)
	fmt.Printf("Context:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expires:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

// alarm-check
var (
	alarmCheckVerbose bool
	alarmCheckJSON    bool
)

var alarmCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check pending alarms",
	Long: `Check and execute alarms for current time.
Includes automatic recovery of missed alarms (últimos 60 minutos).

This command is designed to run from cron every minute.
If no alarms, produces no output (silent).

Examples:
  clical alarm check --user alice
  clical alarm check --user alice --verbose
  clical alarm check --user alice --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// Verificar alarmas en el momento actual
		now := time.Now()
		alarms, err := store.CheckAlarms(userID, now)
		if err != nil {
			return fmt.Errorf("error verifying alarmas: %w", err)
		}

		// Si no hay alarmas, salir silenciosamente
		if len(alarms) == 0 {
			if alarmCheckVerbose {
				fmt.Fprintf(cmd.ErrOrStderr(), "No alarms to execute at this time\n")
			}
			return nil
		}

		// Output en JSON o texto
		if alarmCheckJSON {
			// Emitir JSON a stdout
			jsonData, err := json.MarshalIndent(alarms, "", "  ")
			if err != nil {
				return fmt.Errorf("error serializing alarms: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			// Emitir reporte en texto
			for _, alm := range alarms {
				fmt.Printf("=== %s\n", alm.Context)
				fmt.Printf("    ID: %s\n", alm.ID)
				fmt.Printf("    Recurrence: %s\n", capitalizeRecurrence(alm.Recurrence))
				if !alm.ScheduledFor.IsZero() {
					fmt.Printf("    Scheduled for: %s\n", alm.ScheduledFor.Format("2006-01-02T15:04:05-07:00"))
				}
				fmt.Println()
			}
		}

		return nil
	},
}

// capitalizeRecurrence capitaliza el tipo de recurrencia para display
func capitalizeRecurrence(r alarm.Recurrence) string {
	switch r {
	case alarm.RecurrenceOnce:
		return "One-time"
	case alarm.RecurrenceDaily:
		return "Daily"
	case alarm.RecurrenceWeekly:
		return "Weekly"
	case alarm.RecurrenceMonthly:
		return "Monthly"
	case alarm.RecurrenceYearly:
		return "Yearly"
	default:
		return string(r)
	}
}

// alarm-list
var (
	alarmListPast bool
	alarmListJSON bool
)

var alarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List alarms",
	Long: `List all active alarms for user.

Examples:
  clical alarm list --user alice
  clical alarm list --user alice --past
  clical alarm list --user alice --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		// List active alarms
		activeAlarms, err := store.ListActiveAlarms(userID)
		if err != nil {
			return fmt.Errorf("error listing alarmas activas: %w", err)
		}

		var pastAlarms []*alarm.Alarm
		if alarmListPast {
			pastAlarms, err = store.ListPastAlarms(userID)
			if err != nil {
				return fmt.Errorf("error listing alarmas pasadas: %w", err)
			}
		}

		// Output JSON
		if alarmListJSON {
			output := map[string][]*alarm.Alarm{
				"active": activeAlarms,
			}
			if alarmListPast {
				output["past"] = pastAlarms
			}

			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("error serializing alarmas: %w", err)
			}

			fmt.Println(string(jsonData))
			return nil
		}

		// Output tabla
		if len(activeAlarms) == 0 && len(pastAlarms) == 0 {
			fmt.Println("No alarms")
			return nil
		}

		if len(activeAlarms) > 0 {
			fmt.Println("ACTIVE ALARMS:")
			fmt.Println()
			fmt.Printf("%-25s %-10s %-20s %s\n", "ID", "TIPO", "PROGRAMADA", "CONTEXTO")
			fmt.Println(strings.Repeat("-", 100))

			for _, alm := range activeAlarms {
				scheduled := formatSchedule(alm)
				context := alm.Context
				if len(context) > 40 {
					context = context[:37] + "..."
				}
				fmt.Printf("%-25s %-10s %-20s %s\n", alm.ID, alm.Recurrence, scheduled, context)
			}
			fmt.Println()
		}

		if alarmListPast && len(pastAlarms) > 0 {
			fmt.Println("PAST ALARMS:")
			fmt.Println()
			fmt.Printf("%-25s %-10s %-20s %s\n", "ID", "TIPO", "EXECUTED", "CONTEXTO")
			fmt.Println(strings.Repeat("-", 100))

			for _, alm := range pastAlarms {
				executed := ""
				if alm.ExecutedAt != nil {
					executed = alm.ExecutedAt.Format("2006-01-02 15:04")
				}
				context := alm.Context
				if len(context) > 40 {
					context = context[:37] + "..."
				}
				fmt.Printf("%-25s %-10s %-20s %s\n", alm.ID, alm.Recurrence, executed, context)
			}
			fmt.Println()
		}

		return nil
	},
}

func formatSchedule(alm *alarm.Alarm) string {
	// Esto es una simplificación - en producción parsearíamos el filename
	switch alm.Recurrence {
	case alarm.RecurrenceOnce:
		if !alm.ScheduledFor.IsZero() {
			return alm.ScheduledFor.Format("2006-01-02 15:04")
		}
		return "pending"
	case alarm.RecurrenceDaily:
		return "daily (ver filename)"
	case alarm.RecurrenceWeekly:
		return "weekly (ver filename)"
	case alarm.RecurrenceMonthly:
		return "monthly (ver filename)"
	case alarm.RecurrenceYearly:
		return "yearly (ver filename)"
	default:
		return "unknown"
	}
}

// alarm-cancel
var alarmCancelCmd = &cobra.Command{
	Use:   "cancel ALARM_ID",
	Short: "Cancel an alarm",
	Long: `Cancel (delete) an active alarm by its ID.

Examples:
  clical alarm cancel --user alice alarm_once_1234567890_abcd1234
  clical alarm cancel --user alice alarm_weekly_1234567890_abcd1234`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		alarmID := args[0]

		if err := store.CancelAlarm(userID, alarmID); err != nil {
			return fmt.Errorf("error canceling alarma: %w", err)
		}

		fmt.Printf("✓ Alarm canceled successfully: %s\n", alarmID)

		return nil
	},
}

func init() {
	// alarm add
	alarmAddCmd.Flags().StringVar(&alarmContext, "context", "", "Alarm context (required)")
	alarmAddCmd.Flags().StringVar(&alarmAt, "at", "", "Date/time for one-time alarm (eg: '2025-11-23 14:30', 'tomorrow 10:00', '+30m')")
	alarmAddCmd.Flags().StringVar(&alarmDaily, "daily", "", "Time for daily alarm (eg: '14:30')")
	alarmAddCmd.Flags().StringVar(&alarmWeekly, "weekly", "", "Day and time for weekly alarm (eg: 'monday 14:30')")
	alarmAddCmd.Flags().StringVar(&alarmMonthly, "monthly", "", "Day of month and time (eg: '15 14:30')")
	alarmAddCmd.Flags().StringVar(&alarmYearly, "yearly", "", "Yearly date and time (eg: '11-21 14:30')")
	alarmAddCmd.Flags().StringVar(&alarmExpires, "expires", "", "Expiration date for recurring alarms (eg: '2025-12-31')")

	// alarm check
	alarmCheckCmd.Flags().BoolVarP(&alarmCheckVerbose, "verbose", "v", false, "Show debugging logs")
	alarmCheckCmd.Flags().BoolVar(&alarmCheckJSON, "json", false, "Output in JSON format")

	// alarm list
	alarmListCmd.Flags().BoolVar(&alarmListPast, "past", false, "Include past alarms")
	alarmListCmd.Flags().BoolVar(&alarmListJSON, "json", false, "Output en formato JSON")

	// Add subcommands
	alarmCmd.AddCommand(alarmAddCmd)
	alarmCmd.AddCommand(alarmCheckCmd)
	alarmCmd.AddCommand(alarmListCmd)
	alarmCmd.AddCommand(alarmCancelCmd)
}
