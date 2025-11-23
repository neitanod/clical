# Testing Guide - clical

Esta guía explica cómo ejecutar y escribir tests para clical.

## Suite de Tests

### Tests Implementados

#### 1. `pkg/calendar/entry_test.go`

**Tests de Entry:**
- `TestNewEntry` - Creación de nuevas entradas
- `TestValidate` - Validación de campos requeridos
- `TestEndTime` - Cálculo de hora de finalización
- `TestIsPastFutureCurrent` - Estados temporales de eventos
- `TestGenerateFilename` - Generación de nombres de archivo
- `TestTags` - Gestión de tags (add, remove, has)

**Cobertura:** Todas las funciones públicas de Entry

#### 2. `pkg/calendar/filter_test.go`

**Tests de Filter:**
- `TestFilterMatches` - Coincidencia de filtros
  - Filtros de fecha
  - Filtros de query (título, notas)
  - Filtros de ubicación
  - Filtros de tags
  - Filtros de duración (min/max)
- `TestFilterWithDateRange` - Constructor con rango de fechas
- `TestFilterWithQuery` - Constructor con query
- `TestFilterWithTags` - Constructor con tags

**Cobertura:** Toda la lógica de filtrado

#### 3. `pkg/user/user_test.go`

**Tests de User:**
- `TestNewUser` - Creación de usuarios
- `TestValidate` - Validación de usuarios
  - ID requerido
  - Nombre requerido
  - Timezone válido
  - Duración por defecto válida
- `TestLocation` - Carga de timezone
- `TestFormatFunctions` - Formateo de fecha/hora
- `TestDefaultConfig` - Configuración por defecto

**Cobertura:** Todas las funciones de User

## Ejecutar Tests

### Todos los Tests

```bash
# Ejecutar todos los tests
go test ./...

# Con output verbose
go test -v ./...

# Con cobertura
go test -cover ./...

# Cobertura detallada
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Tests de un Paquete Específico

```bash
# Solo tests de calendar
go test ./pkg/calendar -v

# Solo tests de user
go test ./pkg/user -v

# Solo tests de storage
go test ./pkg/storage -v
```

### Tests Específicos

```bash
# Ejecutar un test específico
go test ./pkg/calendar -run TestNewEntry -v

# Ejecutar tests que coincidan con patrón
go test ./pkg/calendar -run TestValidate -v
```

## Resultados Esperados

### Suite Completa

```bash
$ go test ./...
ok      github.com/sebasvalencia/clical/pkg/calendar    0.003s
ok      github.com/sebasvalencia/clical/pkg/user       0.002s
?       github.com/sebasvalencia/clical/cmd/clical      [no test files]
?       github.com/sebasvalencia/clical/internal/cli    [no test files]
```

### Con Cobertura

```bash
$ go test -cover ./...
ok      github.com/sebasvalencia/clical/pkg/calendar    0.003s  coverage: 85.7% of statements
ok      github.com/sebasvalencia/clical/pkg/user       0.002s  coverage: 90.2% of statements
```

## Escribir Nuevos Tests

### Estructura de un Test

```go
package mypackage

import "testing"

func TestMyFunction(t *testing.T) {
    // Arrange - Preparar datos
    input := "test"
    expected := "expected result"

    // Act - Ejecutar función
    got := MyFunction(input)

    // Assert - Verificar resultado
    if got != expected {
        t.Errorf("MyFunction(%s) = %s, want %s", input, got, expected)
    }
}
```

### Tests con Table-Driven

```go
func TestMultipleCases(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "case 1",
            input:    "input1",
            expected: "output1",
            wantErr:  false,
        },
        {
            name:     "case 2",
            input:    "input2",
            expected: "output2",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.expected {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## Convenciones de Naming

- Archivos de test: `*_test.go`
- Funciones de test: `TestXxx(t *testing.T)`
- Benchmarks: `BenchmarkXxx(b *testing.B)`
- Examples: `ExampleXxx()`

## Tests Pendientes

### Próximos Tests a Implementar

1. **pkg/storage/filesystem_test.go**
   - SaveEntry / GetEntry
   - ListEntries con filtros
   - UpdateEntry / DeleteEntry
   - SaveUser / GetUser
   - ReportState

2. **pkg/reporter/daily_test.go**
   - GenerateDailyReport
   - FormatDailyReport
   - calculateSummary
   - calculateFreetime

3. **internal/cli Tests** (integration tests)
   - add command
   - list command
   - edit/delete commands
   - daily-report command

4. **internal/config Tests**
   - LoadConfig from file
   - LoadConfig from env
   - Priority order

## Integración Continua (CI)

### GitHub Actions (ejemplo)

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - run: go test -v -cover ./...
```

## Makefile Targets

```bash
# Ejecutar tests
make test

# Tests con cobertura
make test-coverage

# Limpiar archivos de cobertura
make clean
```

## Tips para Testing

### 1. Usar Subtests

```go
t.Run("subtest name", func(t *testing.T) {
    // test code
})
```

**Ventajas:**
- Tests más organizados
- Fácil identificar qué falló
- Ejecutar subtests específicos

### 2. Helper Functions

```go
func TestMain(m *testing.M) {
    // Setup
    os.Exit(m.Run())
    // Teardown
}

func setup() {
    // preparar ambiente
}

func teardown() {
    // limpiar ambiente
}
```

### 3. Fixtures

Crear datos de prueba en archivos separados:

```go
func createTestEntry() *Entry {
    return &Entry{
        ID:       "test123",
        UserID:   "12345",
        Title:    "Test Event",
        DateTime: time.Now(),
        Duration: 60,
    }
}
```

### 4. Mocking

Para tests que requieren filesystem, usar:
- Temporary directories: `t.TempDir()`
- In-memory storage
- Interfaces para dependency injection

## Debugging Tests

```bash
# Ejecutar con más información
go test -v ./pkg/calendar

# Ver qué tests se ejecutan sin ejecutarlos
go test -v -run=TestNewEntry ./pkg/calendar -dry-run

# Con race detector
go test -race ./...

# Con timeout personalizado
go test -timeout 30s ./...
```

## Benchmarks (Futuro)

```go
func BenchmarkGenerateID(b *testing.B) {
    for i := 0; i < b.N; i++ {
        GenerateID()
    }
}

// Ejecutar: go test -bench=. -benchmem
```

## Coverage Goals

### Objetivos de Cobertura

- **Critical packages** (calendar, storage): 90%+
- **Business logic** (reporter): 80%+
- **CLI commands**: 70%+
- **Overall**: 75%+

### Generar Reporte de Cobertura

```bash
# Generar coverage profile
go test -coverprofile=coverage.out ./...

# Ver cobertura por función
go tool cover -func=coverage.out

# HTML report
go tool cover -html=coverage.out -o coverage.html

# Abrir en navegador
open coverage.html
```

## Troubleshooting

### Tests Fallan en CI pero No Localmente

- Verificar timezones: usar UTC en tests
- Verificar permisos de archivos
- Verificar variables de entorno

### Tests Lentos

- Usar `t.Parallel()` para tests independientes
- Mockear I/O pesado
- Reducir datos de prueba

### Tests Flaky

- Evitar dependencias de tiempo real (usar mocks)
- Evitar dependencias de orden de ejecución
- Usar seeds fijos para random

## Recursos

- **Go Testing Package**: https://pkg.go.dev/testing
- **Go Testing Best Practices**: https://go.dev/doc/tutorial/add-a-test
- **Table-Driven Tests**: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests

---

**Estado Actual:** 3 paquetes con tests (calendar, filter, user)
**Cobertura Actual:** ~85%
**Próximo:** Tests de storage y reporter
