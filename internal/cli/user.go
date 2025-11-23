package cli

import (
	"fmt"

	"github.com/sebasvalencia/clical/pkg/user"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Gestión de usuarios",
	Long:  "Comandos para gestionar usuarios del sistema de calendario",
}

// user add
var (
	userAddID       string
	userAddName     string
	userAddTimezone string
)

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Crear un nuevo usuario",
	Long: `Crea un nuevo usuario en el sistema.

Ejemplos:
  clical user add --id=12345 --name="Juan Pérez" --timezone="America/Argentina/Buenos_Aires"
  clical user add --id=67890 --name="María García" --timezone="America/Mexico_City"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validar argumentos
		if userAddID == "" {
			return fmt.Errorf("se requiere --id")
		}
		if userAddName == "" {
			return fmt.Errorf("se requiere --name")
		}
		if userAddTimezone == "" {
			return fmt.Errorf("se requiere --timezone")
		}

		// Crear usuario
		u := user.NewUser(userAddID, userAddName, userAddTimezone)

		// Validar
		if err := u.Validate(); err != nil {
			return fmt.Errorf("usuario inválido: %w", err)
		}

		// Guardar
		if err := store.SaveUser(u); err != nil {
			return fmt.Errorf("error guardando usuario: %w", err)
		}

		fmt.Printf("✓ Usuario creado exitosamente\n\n")
		fmt.Printf("ID:       %s\n", u.ID)
		fmt.Printf("Nombre:   %s\n", u.Name)
		fmt.Printf("Timezone: %s\n", u.Timezone)
		fmt.Printf("Creado:   %s\n", u.Created.Format("2006-01-02 15:04"))

		return nil
	},
}

// user list
var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar todos los usuarios",
	Long:  "Lista todos los usuarios registrados en el sistema",
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := store.ListUsers()
		if err != nil {
			return fmt.Errorf("error listando usuarios: %w", err)
		}

		if len(users) == 0 {
			fmt.Println("No hay usuarios registrados")
			return nil
		}

		fmt.Printf("Se encontraron %d usuario(s)\n\n", len(users))

		for _, u := range users {
			fmt.Printf("ID: %s\n", u.ID)
			fmt.Printf("  Nombre:   %s\n", u.Name)
			fmt.Printf("  Timezone: %s\n", u.Timezone)
			fmt.Printf("  Creado:   %s\n", u.Created.Format("2006-01-02"))
			fmt.Println()
		}

		return nil
	},
}

// user show
var userShowID string

var userShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Mostrar detalles de un usuario",
	Long: `Muestra información detallada de un usuario.

Ejemplos:
  clical user show --id=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userShowID == "" {
			return fmt.Errorf("se requiere --id")
		}

		u, err := store.GetUser(userShowID)
		if err != nil {
			return fmt.Errorf("error obteniendo usuario: %w", err)
		}

		fmt.Printf("═══════════════════════════════════════════\n")
		fmt.Printf(" Usuario: %s\n", u.Name)
		fmt.Printf("═══════════════════════════════════════════\n\n")

		fmt.Printf("ID:       %s\n", u.ID)
		fmt.Printf("Timezone: %s\n", u.Timezone)
		fmt.Printf("Creado:   %s\n\n", u.Created.Format("2006-01-02 15:04"))

		fmt.Printf("Configuración:\n")
		fmt.Printf("  Duración por defecto: %d minutos\n", u.Config.DefaultDuration)
		fmt.Printf("  Formato de fecha:     %s\n", u.Config.DateFormat)
		fmt.Printf("  Formato de hora:      %s\n", u.Config.TimeFormat)

		firstDay := "Domingo"
		if u.Config.FirstDayOfWeek == 1 {
			firstDay = "Lunes"
		}
		fmt.Printf("  Primer día semana:    %s\n", firstDay)

		return nil
	},
}

func init() {
	// user add
	userAddCmd.Flags().StringVar(&userAddID, "id", "", "ID del usuario")
	userAddCmd.Flags().StringVar(&userAddName, "name", "", "Nombre del usuario")
	userAddCmd.Flags().StringVar(&userAddTimezone, "timezone", "", "Timezone (ej: America/Argentina/Buenos_Aires)")
	userAddCmd.MarkFlagRequired("id")
	userAddCmd.MarkFlagRequired("name")
	userAddCmd.MarkFlagRequired("timezone")

	// user show
	userShowCmd.Flags().StringVar(&userShowID, "id", "", "ID del usuario")
	userShowCmd.MarkFlagRequired("id")

	// Agregar subcomandos a user
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userShowCmd)
}
