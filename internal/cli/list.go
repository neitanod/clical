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
	Short: "Listar eventos del calendario",
	Long: `Lista eventos del calendario con filtros opcionales.

Ejemplos:
  clical list --user=12345
  clical list --user=12345 --from="2025-11-20"
  clical list --user=12345 --range=today
  clical list --user=12345 --range=week
  clical list --user=12345 --tags=trabajo,reunion`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validar user ID
		if userID == "" {
			return fmt.Errorf("se requiere --user")
		}

		// Construir filtro
		filter := calendar.NewFilter()

		// Aplicar rango predefinido
		if listRange != "" {
			from, to, err := parseRange(listRange)
			if err != nil {
				return err
			}
			filter.From = &from
			filter.To = &to
		}

		// Aplicar from/to manual (sobrescribe range)
		if listFrom != "" {
			from, err := time.Parse("2006-01-02", listFrom)
			if err != nil {
				return fmt.Errorf("error parseando --from: %w", err)
			}
			filter.From = &from
		}

		if listTo != "" {
			to, err := time.Parse("2006-01-02", listTo)
			if err != nil {
				return fmt.Errorf("error parseando --to: %w", err)
			}
			// Incluir todo el día final
			to = to.Add(24 * time.Hour)
			filter.To = &to
		}

		// Aplicar tags
		if len(listTags) > 0 {
			filter.Tags = listTags
		}

		// Obtener eventos
		entries, err := store.ListEntries(userID, filter)
		if err != nil {
			return fmt.Errorf("error listando eventos: %w", err)
		}

		// Mostrar resultados
		if len(entries) == 0 {
			fmt.Println("No se encontraron eventos")
			return nil
		}

		fmt.Printf("Se encontraron %d evento(s)\n\n", len(entries))

		for _, entry := range entries {
			printEntryRow(entry)
			fmt.Println()
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listFrom, "from", "", "Fecha inicial (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listTo, "to", "", "Fecha final (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&listRange, "range", "", "Rango predefinido: today, week, month")
	listCmd.Flags().StringSliceVar(&listTags, "tags", []string{}, "Filtrar por tags")
}

// parseRange parsea rangos predefinidos como "today", "week", "month"
func parseRange(r string) (time.Time, time.Time, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch r {
	case "today":
		return today, today.Add(24 * time.Hour), nil

	case "week":
		// Desde hoy hasta dentro de 7 días
		return today, today.Add(7 * 24 * time.Hour), nil

	case "month":
		// Desde hoy hasta fin de mes
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
		return today, endOfMonth, nil

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("rango inválido: %s (use: today, week, month)", r)
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
