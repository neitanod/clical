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
	Short: "Agregar un nuevo evento al calendario",
	Long: `Agrega un nuevo evento al calendario del usuario especificado.

Ejemplos:
  clical add --user=12345 --datetime="2025-11-20 14:00" --title="Reunión" --duration=60
  clical add --user=12345 --datetime="2025-11-20 14:00" --title="Llamada" --duration=30 --location="Zoom"
  clical add --user=12345 --datetime="2025-11-21 09:00" --title="Stand-up" --duration=15 --tags=trabajo,equipo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validar user ID
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Validar título
		if addTitle == "" {
			return fmt.Errorf("se requiere --title")
		}

		// Parsear datetime
		datetime, err := parseDateTime(addDatetime)
		if err != nil {
			return fmt.Errorf("error parseando --datetime: %w", err)
		}

		// Crear entrada
		entry := calendar.NewEntry(userID, addTitle, datetime, addDuration)
		entry.Location = addLocation
		entry.Notes = addNotes
		entry.Tags = addTags

		// Guardar
		if err := store.SaveEntry(userID, entry); err != nil {
			return fmt.Errorf("error guardando entrada: %w", err)
		}

		// Mostrar confirmación
		fmt.Printf("✓ Evento creado exitosamente\n\n")
		fmt.Printf("ID:       %s\n", entry.ID)
		fmt.Printf("Título:   %s\n", entry.Title)
		fmt.Printf("Fecha:    %s\n", entry.DateTime.Format("2006-01-02 15:04"))
		fmt.Printf("Duración: %d minutos\n", entry.Duration)
		if entry.Location != "" {
			fmt.Printf("Ubicación: %s\n", entry.Location)
		}
		if len(entry.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(entry.Tags, ", "))
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addDatetime, "datetime", "", "Fecha y hora del evento (YYYY-MM-DD HH:MM)")
	addCmd.Flags().StringVar(&addTitle, "title", "", "Título del evento")
	addCmd.Flags().IntVar(&addDuration, "duration", 60, "Duración en minutos")
	addCmd.Flags().StringVar(&addLocation, "location", "", "Ubicación del evento")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "Notas adicionales")
	addCmd.Flags().StringSliceVar(&addTags, "tags", []string{}, "Tags (separados por coma)")

	addCmd.MarkFlagRequired("datetime")
	addCmd.MarkFlagRequired("title")
}

// parseDateTime parsea una fecha/hora en formato "YYYY-MM-DD HH:MM"
func parseDateTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato inválido (use YYYY-MM-DD HH:MM)")
}
