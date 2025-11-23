.PHONY: build test install clean run help

# Variables
BINARY_NAME=clical
BUILD_DIR=.
CMD_DIR=./cmd/clical

# Build the application
build:
	@echo "Compilando $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✓ Compilación exitosa: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Ejecutando tests..."
	go test -v ./...

# Test with coverage
test-coverage:
	@echo "Ejecutando tests con cobertura..."
	go test -cover ./...

# Install to $GOPATH/bin
install:
	@echo "Instalando $(BINARY_NAME) en $$GOPATH/bin..."
	go install $(CMD_DIR)
	@echo "✓ Instalación exitosa"

# Install to /usr/local/bin (requires sudo)
install-system:
	@echo "Instalando $(BINARY_NAME) en /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Instalación exitosa en /usr/local/bin"

# Clean build artifacts
clean:
	@echo "Limpiando archivos generados..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	go clean
	@echo "✓ Limpieza completada"

# Run the application (for development)
run:
	go run $(CMD_DIR)

# Format code
fmt:
	@echo "Formateando código..."
	go fmt ./...
	@echo "✓ Formato aplicado"

# Vet code
vet:
	@echo "Analizando código..."
	go vet ./...
	@echo "✓ Análisis completado"

# Show help
help:
	@echo "Makefile para $(BINARY_NAME)"
	@echo ""
	@echo "Uso:"
	@echo "  make build           - Compilar el binario"
	@echo "  make test            - Ejecutar tests"
	@echo "  make test-coverage   - Ejecutar tests con cobertura"
	@echo "  make install         - Instalar en \$$GOPATH/bin"
	@echo "  make install-system  - Instalar en /usr/local/bin (requiere sudo)"
	@echo "  make clean           - Limpiar archivos generados"
	@echo "  make run             - Ejecutar en modo desarrollo"
	@echo "  make fmt             - Formatear código"
	@echo "  make vet             - Analizar código"
	@echo "  make help            - Mostrar esta ayuda"
