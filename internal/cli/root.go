package cli

import (
	"fmt"
	"os"

	"github.com/sebasvalencia/clical/internal/config"
	"github.com/sebasvalencia/clical/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	store      storage.Storage
	configPath string
	dataDir    string
	userID     string
)

// rootCmd es el comando raíz
var rootCmd = &cobra.Command{
	Use:   "clical",
	Short: "Multi-user calendar CLI system",
	Long: `clical is a multi-user calendar system with command line interface,
designed for AI assistance.

Data is stored in Markdown + JSON files, organized by date,
which allows manual navigation and editing of events.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Cargar configuración desde archivo
		var err error
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Si se especificó --data-dir, sobrescribir config
		if cmd.Flags().Changed("data-dir") {
			cfg.DataDir = dataDir
		} else {
			dataDir = cfg.DataDir
		}

		// Si se especificó --user, sobrescribir config
		if cmd.Flags().Changed("user") {
			cfg.UserID = userID
		} else if cfg.UserID != "" {
			userID = cfg.UserID
		}

		// Inicializar storage
		store, err = storage.NewFilesystemStorage(cfg.DataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute ejecuta el comando raíz
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Flags globales
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "/etc/clical/config.env", "Configuration file")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Data directory (overrides config)")
	rootCmd.PersistentFlags().StringVar(&userID, "user", "", "User ID (overrides config)")

	// Agregar subcomandos
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(alarmCmd)
	rootCmd.AddCommand(versionCmd)
}
