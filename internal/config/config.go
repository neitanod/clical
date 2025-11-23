package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config contiene la configuración global de la aplicación
type Config struct {
	DataDir  string
	UserID   string // Usuario por defecto si no se especifica
	LogLevel string
}

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		DataDir:  "/var/clical/data",
		UserID:   "",
		LogLevel: "info",
	}
}

// LoadConfig carga la configuración desde archivo .env
func LoadConfig(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	// Cargar desde archivo si existe
	if configPath != "" {
		if err := loadFromFile(cfg, configPath); err != nil {
			// Si no existe el archivo, continuar con defaults
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("error cargando config: %w", err)
			}
		}
	}

	// Permitir override desde variables de entorno
	if dataDir := os.Getenv("CLICAL_DATA_DIR"); dataDir != "" {
		cfg.DataDir = dataDir
	}

	if userID := os.Getenv("CLICAL_USER_ID"); userID != "" {
		cfg.UserID = userID
	}

	if logLevel := os.Getenv("CLICAL_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return cfg, nil
}

// loadFromFile carga configuración desde archivo .env
func loadFromFile(cfg *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignorar líneas vacías y comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parsear KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Aplicar configuración
		switch key {
		case "CLICAL_DATA_DIR":
			if value != "" {
				cfg.DataDir = value
			}
		case "CLICAL_USER_ID":
			if value != "" {
				cfg.UserID = value
			}
		case "CLICAL_LOG_LEVEL":
			if value != "" {
				cfg.LogLevel = value
			}
		}
	}

	return scanner.Err()
}
