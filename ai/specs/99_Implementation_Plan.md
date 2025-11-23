# Plan de Implementación de clical

**Fecha:** 2025-11-20
**Versión:** 1.0

## Resumen

Este documento describe el plan paso a paso para implementar clical desde cero hasta tener un MVP funcional.

## Fase 1: Inicialización del Proyecto (30 min)

### Paso 1.1: Inicializar Go Module
```bash
cd /home/sebas/doc/prj/clical
go mod init github.com/tu-usuario/clical
```

### Paso 1.2: Instalar Dependencias
```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
```

### Paso 1.3: Crear Estructura de Directorios
```bash
mkdir -p cmd/clical
mkdir -p pkg/{calendar,storage,user,formatter,parser,view,reporter,importer}
mkdir -p internal/{cli,config,util}
mkdir -p data/users
mkdir -p ai/specs
```

### Paso 1.4: Crear .gitignore
```
# Binarios
clical
*.exe

# Datos locales
/data/

# Configuración local
.env

# Go
*.o
*.a
*.so
```

## Fase 2: Modelos de Datos (45 min)

### Paso 2.1: Implementar pkg/calendar/entry.go
- Struct Entry con todos los campos
- Métodos básicos: Validate(), GenerateID()
- Constantes para estados

### Paso 2.2: Implementar pkg/user/user.go
- Struct User
- Struct UserConfig
- Métodos de validación

### Paso 2.3: Implementar pkg/calendar/filter.go
- Struct Filter para búsquedas
- Métodos de construcción de filtros

### Paso 2.4: Tests Básicos
- entry_test.go
- user_test.go

## Fase 3: Storage Layer (1.5 horas)

### Paso 3.1: Definir Interface pkg/storage/storage.go
```go
type Storage interface {
    SaveEntry(userID string, entry *calendar.Entry) error
    GetEntry(userID, entryID string) (*calendar.Entry, error)
    ListEntries(userID string, filter calendar.Filter) ([]*calendar.Entry, error)
    DeleteEntry(userID, entryID string) error
    UpdateEntry(userID string, entry *calendar.Entry) error

    SaveUser(user *user.User) error
    GetUser(userID string) (*user.User, error)
    ListUsers() ([]*user.User, error)
}
```

### Paso 3.2: Implementar pkg/storage/filesystem.go
- Struct FilesystemStorage
- Constructor NewFilesystemStorage(dataDir string)
- Implementar SaveEntry con generación de Markdown + JSON
- Implementar GetEntry
- Implementar ListEntries con filtrado
- Implementar DeleteEntry
- Implementar UpdateEntry

### Paso 3.3: Implementar pkg/storage/markdown.go
- Función entryToMarkdown(*Entry) string
- Función markdownToEntry(string) (*Entry, error)
- Función entryToJSON(*Entry) ([]byte, error)
- Función jsonToEntry([]byte) (*Entry, error)

### Paso 3.4: Implementar pkg/storage/paths.go
- Función getEntryPath(dataDir, userID, entryID, datetime string) string
- Función getUserPath(dataDir, userID string) string
- Función parseEntryFilename(filename string) (datetime, title string, error)

### Paso 3.5: Tests
- filesystem_test.go con casos completos
- markdown_test.go

## Fase 4: CLI Foundation (1 hora)

### Paso 4.1: Setup Cobra en cmd/clical/main.go
```go
package main

import (
    "github.com/tu-usuario/clical/internal/cli"
)

func main() {
    cli.Execute()
}
```

### Paso 4.2: Implementar internal/cli/root.go
- Comando raíz con descripción
- Flags globales: --data-dir, --user
- Configuración de Viper

### Paso 4.3: Implementar internal/config/config.go
- Struct Config
- LoadConfig() función
- DefaultConfig() función
- Lectura de variables de entorno

## Fase 5: Comandos Básicos (2 horas)

### Paso 5.1: Implementar internal/cli/add.go
```bash
clical add --user=12345 --datetime="2025-11-20 14:00" --title="Reunión" --duration=60
```
- Parse de argumentos
- Validación
- Crear Entry
- Guardar con Storage
- Formatear output

### Paso 5.2: Implementar internal/cli/list.go
```bash
clical list --user=12345 --from="2025-11-20" --to="2025-11-30"
```
- Parse de argumentos
- Construir Filter
- Obtener entries de Storage
- Formatear como tabla

### Paso 5.3: Implementar internal/cli/show.go
```bash
clical show --user=12345 --id=abc123
```
- Obtener entry
- Formatear detallado

### Paso 5.4: Implementar internal/cli/delete.go
```bash
clical delete --user=12345 --id=abc123
```
- Confirmación
- Eliminar con Storage

### Paso 5.5: Implementar internal/cli/edit.go
```bash
clical edit --user=12345 --id=abc123 --title="Nuevo título"
```
- Obtener entry existente
- Actualizar campos
- Guardar

## Fase 6: Formatters (45 min)

### Paso 6.1: Definir Interface pkg/formatter/formatter.go
```go
type Formatter interface {
    FormatEntry(entry *calendar.Entry) (string, error)
    FormatEntries(entries []*calendar.Entry) (string, error)
}
```

### Paso 6.2: Implementar pkg/formatter/text.go
- Formato legible para humanos
- Uso de colores ANSI (opcional)

### Paso 6.3: Implementar pkg/formatter/table.go
- Formato tabla para listas
- Columnas: DateTime, Title, Duration, Location

### Paso 6.4: Implementar pkg/formatter/json.go
- JSON pretty-printed

## Fase 7: Views (1 hora)

### Paso 7.1: Implementar pkg/view/day.go
```go
func RenderDay(entries []*calendar.Entry, date time.Time) string
```
- Timeline del día
- Eventos cronológicos
- Bloques libres

### Paso 7.2: Implementar pkg/view/week.go
- Vista semanal
- Columnas por día

### Paso 7.3: Implementar pkg/view/month.go
- Calendario mensual estilo `cal`
- Marcadores en días con eventos

### Paso 7.4: Implementar comandos CLI
- internal/cli/day.go
- internal/cli/week.go
- internal/cli/month.go

## Fase 8: Reportes para IA (2 horas)

### Paso 8.1: Implementar pkg/reporter/daily.go
```go
type DailyReport struct {
    Date           time.Time
    Events         []*calendar.Entry
    Summary        Summary
    FreetimeBlocks []FreetimeBlock
    NextDay        []calendar.Entry
    Suggestions    []string
}

func GenerateDailyReport(storage Storage, userID string, date time.Time) (*DailyReport, error)
```

### Paso 8.2: Implementar pkg/reporter/markdown.go
- Función reportToMarkdown(*DailyReport) string
- Formato optimizado para IA
- Incluir todos los campos del reporte

### Paso 8.3: Implementar pkg/reporter/state.go
- Gestión de .state/daily-reports.json
- Tracking de eventos ya reportados
- Evitar duplicados

### Paso 8.4: Implementar internal/cli/daily_report.go
```bash
clical daily-report --user=12345
clical daily-report --user=12345 --date="2025-11-20"
```

### Paso 8.5: Implementar internal/cli/tomorrow_report.go
```bash
clical tomorrow-report --user=12345
```

### Paso 8.6: Implementar internal/cli/upcoming_report.go
```bash
clical upcoming-report --user=12345 --hours=2
clical upcoming-report --user=12345 --hours=2 --only-new
```

### Paso 8.7: Implementar internal/cli/weekly_report.go
```bash
clical weekly-report --user=12345
```

## Fase 9: User Management (45 min)

### Paso 9.1: Implementar internal/cli/user.go
- Subcomando raíz para user

### Paso 9.2: Implementar comandos user
```bash
clical user add --id=12345 --name="Juan" --timezone="America/Argentina/Buenos_Aires"
clical user list
clical user show --id=12345
clical user delete --id=12345
```

### Paso 9.3: Storage de usuarios
- Implementar SaveUser, GetUser en filesystem.go
- user.md y user.json

## Fase 10: Search & Filter (45 min)

### Paso 10.1: Implementar internal/cli/search.go
```bash
clical search --user=12345 --query="reunión"
clical search --user=12345 --title="cliente"
```

### Paso 10.2: Implementar internal/cli/filter.go
```bash
clical filter --user=12345 --location="Oficina"
clical filter --user=12345 --duration-min=30
```

### Paso 10.3: Mejorar pkg/storage para búsquedas eficientes
- Indexación opcional
- Grep en archivos markdown

## Fase 11: Utils & Config (30 min)

### Paso 11.1: Implementar internal/cli/config.go
```bash
clical config show
clical config set --user=12345 --default-duration=60
```

### Paso 11.2: Implementar internal/cli/info.go
```bash
clical info
clical version
```

## Fase 12: Testing & Documentation (1 hora)

### Paso 12.1: Tests Unitarios
- Completar tests de todos los paquetes
- Coverage mínimo 70%

### Paso 12.2: Integration Tests
- Test end-to-end de comandos
- Test de storage con datos reales

### Paso 12.3: Actualizar README.md
- Instalación
- Ejemplos de uso
- Comandos principales

### Paso 12.4: Actualizar docs/
- Documentación web
- Guía de usuario
- Referencia de comandos

## Fase 13: Build & Deploy (30 min)

### Paso 13.1: Crear Makefile
```makefile
build:
    go build -o clical ./cmd/clical

test:
    go test ./...

install:
    go install ./cmd/clical

clean:
    rm -f clical
```

### Paso 13.2: Compilar
```bash
make build
./clical --help
```

### Paso 13.3: Test de uso real
- Crear usuario de prueba
- Agregar eventos
- Probar todos los comandos
- Generar reportes

### Paso 13.4: Configurar crontab de ejemplo
```bash
# Editar crontab
crontab -e

# Agregar líneas:
0 7 * * * /usr/local/bin/clical daily-report --user=12345 | slg-gobot-send-message "$(cat)"
0 20 * * * /usr/local/bin/clical tomorrow-report --user=12345 | slg-gobot-send-message "$(cat)"
```

## Cronograma Estimado

| Fase | Descripción | Tiempo | Acumulado |
|------|-------------|--------|-----------|
| 1 | Inicialización | 30 min | 30 min |
| 2 | Modelos | 45 min | 1h 15min |
| 3 | Storage | 1.5h | 2h 45min |
| 4 | CLI Foundation | 1h | 3h 45min |
| 5 | Comandos Básicos | 2h | 5h 45min |
| 6 | Formatters | 45 min | 6h 30min |
| 7 | Views | 1h | 7h 30min |
| 8 | Reportes IA | 2h | 9h 30min |
| 9 | User Management | 45 min | 10h 15min |
| 10 | Search/Filter | 45 min | 11h |
| 11 | Utils/Config | 30 min | 11h 30min |
| 12 | Testing/Docs | 1h | 12h 30min |
| 13 | Build/Deploy | 30 min | 13h |

**Total estimado:** 13 horas de desarrollo

## Prioridades

### MVP Mínimo (Fases 1-5)
- Inicialización
- Modelos
- Storage
- CLI básico
- Comandos: add, list, show, delete

### MVP Funcional (+ Fases 6-8)
- Formatters
- Views
- Reportes para IA (CARACTERÍSTICA PRINCIPAL)

### Completo (+ Fases 9-13)
- User management
- Search/Filter
- Utils
- Testing completo
- Documentación

## Orden de Ejecución

1. **Comenzar ahora:** Fases 1-2 (Inicialización + Modelos)
2. **Siguiente sesión:** Fases 3-5 (Storage + CLI básico)
3. **Siguiente sesión:** Fases 6-8 (Formatters + Views + Reportes)
4. **Pulir:** Fases 9-13 (Features adicionales + Testing)

## Próximo Paso Inmediato

**Ejecutar Fase 1: Inicializar proyecto Go**

---

**Estado:** Plan completo definido
**Listo para comenzar implementación:** ✅
