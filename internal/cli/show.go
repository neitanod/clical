package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var showID string

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show event details",
	Long: `Muestra todos los detalles de un evento específico.

Examples:
  clical show --user=12345 --id=abc123def456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments
		if userID == "" {
			return fmt.Errorf("--user is required")
		}
		if showID == "" {
			return fmt.Errorf("--id is required")
		}

		// Get event
		entry, err := store.GetEntry(userID, showID)
		if err != nil {
			return fmt.Errorf("error getting event: %w", err)
		}

		// Show complete details
		fmt.Printf("═══════════════════════════════════════════\n")
		fmt.Printf(" %s\n", entry.Title)
		fmt.Printf("═══════════════════════════════════════════\n\n")

		fmt.Printf("ID:        %s\n", entry.ID)
		fmt.Printf("Fecha:     %s\n", entry.DateTime.Format("2006-01-02"))
		fmt.Printf("Hora:      %s\n", entry.DateTime.Format("15:04"))
		fmt.Printf("Duration:  %d minutes\n", entry.Duration)
		fmt.Printf("Fin:       %s\n", entry.EndTime().Format("15:04"))

		if entry.Location != "" {
			fmt.Printf("Ubicación: %s\n", entry.Location)
		}

		if len(entry.Tags) > 0 {
			fmt.Printf("Tags:      #%s\n", strings.Join(entry.Tags, " #"))
		}

		if entry.Notes != "" {
			fmt.Printf("\nNotas:\n%s\n", entry.Notes)
		}

		if len(entry.Metadata) > 0 {
			fmt.Printf("\nMetadata:\n")
			for k, v := range entry.Metadata {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		fmt.Printf("\n───────────────────────────────────────────\n")
		fmt.Printf("Creado:      %s\n", entry.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("Actualizado: %s\n", entry.UpdatedAt.Format("2006-01-02 15:04"))

		return nil
	},
}

func init() {
	showCmd.Flags().StringVar(&showID, "id", "", "ID del evento")
	showCmd.MarkFlagRequired("id")
}
