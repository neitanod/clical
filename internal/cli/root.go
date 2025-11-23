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
	Short: "Sistema de calendario multiusuario CLI",
	Long: `clical es un sistema de calendario multiusuario con interfaz de línea de comandos,
diseñado para ser asistido por Inteligencia Artificial.

Los datos se almacenan en archivos Markdown + JSON, organizados por fecha,
lo que permite navegar y editar manualmente los eventos.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Cargar configuración desde archivo
		var err error
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error cargando configuración: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "Error inicializando storage: %v\n", err)
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
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "/etc/clical/config.env", "Archivo de configuración")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Directorio de datos (sobrescribe config)")
	rootCmd.PersistentFlags().StringVar(&userID, "user", "", "ID de usuario (sobrescribe config)")

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
