package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	deleteID    string
	deleteForce bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an event",
	Long: `Delete a calendar event.

Examples:
  clical delete --user=12345 --id=abc123def456
  clical delete --user=12345 --id=abc123def456 --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments
		if userID == "" {
			return fmt.Errorf("--user is required")
		}
		if deleteID == "" {
			return fmt.Errorf("--id is required")
		}

		// Get event first to show what will be deleted
		entry, err := store.GetEntry(userID, deleteID)
		if err != nil {
			return fmt.Errorf("error getting event: %w", err)
		}

		// Show event
		fmt.Printf("Evento a eliminar:\n")
		fmt.Printf("  %s - %s (%d min)\n",
			entry.DateTime.Format("2006-01-02 15:04"),
			entry.Title,
			entry.Duration,
		)

		// Confirmar a menos que sea --force
		if !deleteForce {
			fmt.Printf("\n¿Está seguro que desea eliminar este evento? (s/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading respuesta: %w", err)
			}

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "s" && response != "si" && response != "sí" {
				fmt.Println("Operación cancelada")
				return nil
			}
		}

		// Delete
		if err := store.DeleteEntry(userID, deleteID); err != nil {
			return fmt.Errorf("error deleting evento: %w", err)
		}

		fmt.Println("✓ Event deleted successfully")

		return nil
	},
}

func init() {
	deleteCmd.Flags().StringVar(&deleteID, "id", "", "ID del evento a eliminar")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Eliminar sin confirmación")
	deleteCmd.MarkFlagRequired("id")
}
