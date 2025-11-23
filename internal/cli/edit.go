package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	editID       string
	editTitle    string
	editDatetime string
	editDuration int
	editLocation string
	editNotes    string
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Editar un evento existente",
	Long: `Modifica los campos de un evento existente.

Ejemplos:
  clical edit --user=12345 --id=abc123 --title="Nuevo título"
  clical edit --user=12345 --id=abc123 --datetime="2025-11-21 15:00"
  clical edit --user=12345 --id=abc123 --duration=90 --location="Sala 2"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validar argumentos
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}
		if editID == "" {
			return fmt.Errorf("se requiere --id")
		}

		// Obtener evento existente
		entry, err := store.GetEntry(userID, editID)
		if err != nil {
			return fmt.Errorf("error obteniendo evento: %w", err)
		}

		// Aplicar cambios
		modified := false

		if cmd.Flags().Changed("title") {
			entry.Title = editTitle
			modified = true
		}

		if cmd.Flags().Changed("datetime") {
			datetime, err := parseDateTime(editDatetime)
			if err != nil {
				return fmt.Errorf("error parseando --datetime: %w", err)
			}
			entry.DateTime = datetime
			modified = true
		}

		if cmd.Flags().Changed("duration") {
			entry.Duration = editDuration
			modified = true
		}

		if cmd.Flags().Changed("location") {
			entry.Location = editLocation
			modified = true
		}

		if cmd.Flags().Changed("notes") {
			entry.Notes = editNotes
			modified = true
		}

		if !modified {
			return fmt.Errorf("no se especificaron cambios")
		}

		// Actualizar timestamp
		entry.UpdatedAt = time.Now()

		// Guardar
		if err := store.UpdateEntry(userID, entry); err != nil {
			return fmt.Errorf("error actualizando evento: %w", err)
		}

		fmt.Printf("✓ Evento actualizado exitosamente\n\n")
		fmt.Printf("ID:       %s\n", entry.ID)
		fmt.Printf("Título:   %s\n", entry.Title)
		fmt.Printf("Fecha:    %s\n", entry.DateTime.Format("2006-01-02 15:04"))
		fmt.Printf("Duración: %d minutos\n", entry.Duration)
		if entry.Location != "" {
			fmt.Printf("Ubicación: %s\n", entry.Location)
		}

		return nil
	},
}

func init() {
	editCmd.Flags().StringVar(&editID, "id", "", "ID del evento a editar")
	editCmd.Flags().StringVar(&editTitle, "title", "", "Nuevo título")
	editCmd.Flags().StringVar(&editDatetime, "datetime", "", "Nueva fecha/hora")
	editCmd.Flags().IntVar(&editDuration, "duration", 0, "Nueva duración")
	editCmd.Flags().StringVar(&editLocation, "location", "", "Nueva ubicación")
	editCmd.Flags().StringVar(&editNotes, "notes", "", "Nuevas notas")

	editCmd.MarkFlagRequired("id")
}
