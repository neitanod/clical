package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sebasvalencia/clical/pkg/alarm"
	"github.com/spf13/cobra"
)

var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Gestionar alarmas",
	Long:  `Gestiona alarmas para recordatorios y seguimientos programados.`,
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
	Use:   "add",
	Short: "Agregar una nueva alarma",
	Long: `Agrega una nueva alarma (one-time o recurrente).

Ejemplos:
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
			return fmt.Errorf("se requiere --user")
		}

		if alarmContext == "" {
			return fmt.Errorf("se requiere --context")
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

		return fmt.Errorf("debe especificar --at, --daily, --weekly, --monthly o --yearly")
	},
}

func addOneTimeAlarm(userID, atStr, context string) error {
	// Parsear fecha/hora
	alarmTime, err := parseDateTime(atStr)
	if err != nil {
		return fmt.Errorf("error parseando --at: %w", err)
	}

	// Verificar que sea futura
	if alarmTime.Before(time.Now()) {
		return fmt.Errorf("la fecha/hora debe ser futura")
	}

	// Redondear a minuto
	alarmTime = alarm.RoundToMinute(alarmTime)

	// Crear alarma
	alm := alarm.NewAlarm(context, alarm.RecurrenceOnce)

	// Guardar
	filename := alarm.OneTimeFilename(alarmTime)
	if err := store.SaveAlarm(userID, alarmTime, alarm.RecurrenceOnce, filename, alm); err != nil {
		return fmt.Errorf("error guardando alarma: %w", err)
	}

	fmt.Printf("✓ Alarma creada exitosamente\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Tipo:       one-time\n")
	fmt.Printf("Programada: %s\n", alarmTime.Format("2006-01-02 15:04"))
	fmt.Printf("Contexto:   %s\n", context)

	return nil
}

func addDailyAlarm(userID, timeStr, context, expiresStr string) error {
	// Parsear hora (HH:MM)
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido para --daily (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("hora inválida: %s", parts[0])
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("minuto inválido: %s", parts[1])
	}

	// Crear alarma
	alm := alarm.NewAlarm(context, alarm.RecurrenceDaily)

	// Agregar expiración si se especificó
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parseando --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Guardar
	schedule := alarm.DailySchedule{Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceDaily, filename, alm); err != nil {
		return fmt.Errorf("error guardando alarma: %w", err)
	}

	fmt.Printf("✓ Alarma creada exitosamente\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Tipo:       daily\n")
	fmt.Printf("Hora:       %02d:%02d\n", hour, minute)
	fmt.Printf("Contexto:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expira:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

func addWeeklyAlarm(userID, scheduleStr, context, expiresStr string) error {
	// Parsear "monday 14:30"
	parts := strings.Fields(scheduleStr)
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido para --weekly (debe ser DAYNAME HH:MM)")
	}

	weekday, err := alarm.ParseWeekday(parts[0])
	if err != nil {
		return err
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) != 2 {
		return fmt.Errorf("formato inválido para hora (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("hora inválida: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("minuto inválido: %s", timeParts[1])
	}

	// Crear alarma
	alm := alarm.NewAlarm(context, alarm.RecurrenceWeekly)

	// Agregar expiración si se especificó
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parseando --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Guardar
	schedule := alarm.WeeklySchedule{Weekday: weekday, Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceWeekly, filename, alm); err != nil {
		return fmt.Errorf("error guardando alarma: %w", err)
	}

	fmt.Printf("✓ Alarma creada exitosamente\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Tipo:       weekly\n")
	fmt.Printf("Día:        %s\n", weekday.String())
	fmt.Printf("Hora:       %02d:%02d\n", hour, minute)
	fmt.Printf("Contexto:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expira:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

func addMonthlyAlarm(userID, scheduleStr, context, expiresStr string) error {
	// Parsear "15 14:30"
	parts := strings.Fields(scheduleStr)
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido para --monthly (debe ser DAY HH:MM)")
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil || day < 1 || day > 31 {
		return fmt.Errorf("día inválido: %s (debe ser 1-31)", parts[0])
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) != 2 {
		return fmt.Errorf("formato inválido para hora (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("hora inválida: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("minuto inválido: %s", timeParts[1])
	}

	// Crear alarma
	alm := alarm.NewAlarm(context, alarm.RecurrenceMonthly)

	// Agregar expiración si se especificó
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parseando --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Guardar
	schedule := alarm.MonthlySchedule{Day: day, Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceMonthly, filename, alm); err != nil {
		return fmt.Errorf("error guardando alarma: %w", err)
	}

	fmt.Printf("✓ Alarma creada exitosamente\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Tipo:       monthly\n")
	fmt.Printf("Día:        %d\n", day)
	fmt.Printf("Hora:       %02d:%02d\n", hour, minute)
	fmt.Printf("Contexto:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expira:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

func addYearlyAlarm(userID, scheduleStr, context, expiresStr string) error {
	// Parsear "11-21 14:30"
	parts := strings.Fields(scheduleStr)
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido para --yearly (debe ser MM-DD HH:MM)")
	}

	dateParts := strings.Split(parts[0], "-")
	if len(dateParts) != 2 {
		return fmt.Errorf("formato inválido para fecha (debe ser MM-DD)")
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
		return fmt.Errorf("formato inválido para hora (debe ser HH:MM)")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("hora inválida: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("minuto inválido: %s", timeParts[1])
	}

	// Crear alarma
	alm := alarm.NewAlarm(context, alarm.RecurrenceYearly)

	// Agregar expiración si se especificó
	if expiresStr != "" {
		expiresAt, err := parseDateTime(expiresStr)
		if err != nil {
			return fmt.Errorf("error parseando --expires: %w", err)
		}
		alm.ExpiresAt = &expiresAt
	}

	// Guardar
	schedule := alarm.YearlySchedule{Month: time.Month(month), Day: day, Hour: hour, Minute: minute}
	filename := schedule.Filename()
	if err := store.SaveAlarm(userID, time.Now(), alarm.RecurrenceYearly, filename, alm); err != nil {
		return fmt.Errorf("error guardando alarma: %w", err)
	}

	fmt.Printf("✓ Alarma creada exitosamente\n\n")
	fmt.Printf("ID:         %s\n", alm.ID)
	fmt.Printf("Tipo:       yearly\n")
	fmt.Printf("Fecha:      %02d-%02d\n", month, day)
	fmt.Printf("Hora:       %02d:%02d\n", hour, minute)
	fmt.Printf("Contexto:   %s\n", context)
	if alm.ExpiresAt != nil {
		fmt.Printf("Expira:     %s\n", alm.ExpiresAt.Format("2006-01-02"))
	}

	return nil
}

// alarm-check
var (
	alarmCheckVerbose bool
)

var alarmCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Verificar alarmas pendientes",
	Long: `Verifica y ejecuta alarmas que corresponden al momento actual.
Incluye recovery automático de alarmas perdidas (últimos 60 minutos).

Este comando está diseñado para ejecutarse desde cron cada minuto.
Si no hay alarmas, no produce output (silencio).

Ejemplos:
  clical alarm check --user alice
  clical alarm check --user alice --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Verificar alarmas en el momento actual
		now := time.Now()
		alarms, err := store.CheckAlarms(userID, now)
		if err != nil {
			return fmt.Errorf("error verificando alarmas: %w", err)
		}

		// Si no hay alarmas, salir silenciosamente
		if len(alarms) == 0 {
			if alarmCheckVerbose {
				fmt.Fprintf(cmd.ErrOrStderr(), "No hay alarmas para ejecutar en este momento\n")
			}
			return nil
		}

		// Emitir JSON a stdout
		jsonData, err := json.MarshalIndent(alarms, "", "  ")
		if err != nil {
			return fmt.Errorf("error serializando alarmas: %w", err)
		}

		fmt.Println(string(jsonData))

		return nil
	},
}

// alarm-list
var (
	alarmListPast bool
	alarmListJSON bool
)

var alarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar alarmas",
	Long: `Lista todas las alarmas activas del usuario.

Ejemplos:
  clical alarm list --user alice
  clical alarm list --user alice --past
  clical alarm list --user alice --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Listar alarmas activas
		activeAlarms, err := store.ListActiveAlarms(userID)
		if err != nil {
			return fmt.Errorf("error listando alarmas activas: %w", err)
		}

		var pastAlarms []*alarm.Alarm
		if alarmListPast {
			pastAlarms, err = store.ListPastAlarms(userID)
			if err != nil {
				return fmt.Errorf("error listando alarmas pasadas: %w", err)
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
				return fmt.Errorf("error serializando alarmas: %w", err)
			}

			fmt.Println(string(jsonData))
			return nil
		}

		// Output tabla
		if len(activeAlarms) == 0 && len(pastAlarms) == 0 {
			fmt.Println("No hay alarmas")
			return nil
		}

		if len(activeAlarms) > 0 {
			fmt.Println("ALARMAS ACTIVAS:")
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
			fmt.Println("ALARMAS PASADAS:")
			fmt.Println()
			fmt.Printf("%-25s %-10s %-20s %s\n", "ID", "TIPO", "EJECUTADA", "CONTEXTO")
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
	Short: "Cancelar una alarma",
	Long: `Cancela (elimina) una alarma activa por su ID.

Ejemplos:
  clical alarm cancel --user alice alarm_once_1234567890_abcd1234
  clical alarm cancel --user alice alarm_weekly_1234567890_abcd1234`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		alarmID := args[0]

		if err := store.CancelAlarm(userID, alarmID); err != nil {
			return fmt.Errorf("error cancelando alarma: %w", err)
		}

		fmt.Printf("✓ Alarma cancelada exitosamente: %s\n", alarmID)

		return nil
	},
}

func init() {
	// alarm add
	alarmAddCmd.Flags().StringVar(&alarmContext, "context", "", "Contexto de la alarma (requerido)")
	alarmAddCmd.Flags().StringVar(&alarmAt, "at", "", "Fecha/hora para alarma one-time (ej: '2025-11-23 14:30', 'tomorrow 10:00', '+30m')")
	alarmAddCmd.Flags().StringVar(&alarmDaily, "daily", "", "Hora para alarma diaria (ej: '14:30')")
	alarmAddCmd.Flags().StringVar(&alarmWeekly, "weekly", "", "Día y hora para alarma semanal (ej: 'monday 14:30')")
	alarmAddCmd.Flags().StringVar(&alarmMonthly, "monthly", "", "Día del mes y hora (ej: '15 14:30')")
	alarmAddCmd.Flags().StringVar(&alarmYearly, "yearly", "", "Fecha y hora anual (ej: '11-21 14:30')")
	alarmAddCmd.Flags().StringVar(&alarmExpires, "expires", "", "Fecha de expiración para recurrentes (ej: '2025-12-31')")

	// alarm check
	alarmCheckCmd.Flags().BoolVarP(&alarmCheckVerbose, "verbose", "v", false, "Mostrar logs de debugging")

	// alarm list
	alarmListCmd.Flags().BoolVar(&alarmListPast, "past", false, "Incluir alarmas pasadas")
	alarmListCmd.Flags().BoolVar(&alarmListJSON, "json", false, "Output en formato JSON")

	// Agregar subcomandos
	alarmCmd.AddCommand(alarmAddCmd)
	alarmCmd.AddCommand(alarmCheckCmd)
	alarmCmd.AddCommand(alarmListCmd)
	alarmCmd.AddCommand(alarmCancelCmd)
}
