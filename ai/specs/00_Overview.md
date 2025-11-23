# Especificación General del Proyecto clical

**Fecha de creación:** 2025-11-20
**Versión:** 1.0

## Visión del Proyecto

**clical** es un sistema de calendario multiusuario con interfaz de línea de comandos (CLI) diseñado para ser asistido por Inteligencia Artificial. El propósito principal es que una IA pueda conocer la agenda del usuario y asistirlo proactivamente en la organización de su día, preparación para eventos y gestión de tareas pendientes.

## Características Principales

### 1. Almacenamiento en Texto Plano

- **Formato primario:** Markdown para legibilidad humana
- **Metadata:** JSON para datos estructurados
- **Estructura de carpetas:** Organización jerárquica por año/mes/día
- **Navegabilidad:** Archivos fáciles de explorar sin la aplicación
- **Versionable:** Compatible con Git para control de versiones

### 2. Multiusuario

- Separación por ID de usuario
- Cada usuario tiene su propio espacio de datos
- Configuración individual por usuario
- Soporte para múltiples usuarios en el mismo sistema

### 3. Interfaz CLI Completa

Subcomandos organizados en categorías:

#### a) Gestión de Entradas
- `add` - Crear eventos
- `list` - Listar eventos
- `show` - Ver detalle de evento
- `edit` - Modificar eventos
- `delete` - Eliminar eventos

#### b) Búsqueda y Filtrado
- `search` - Buscar por texto
- `filter` - Filtrar por criterios

#### c) Vistas de Calendario
- `month` - Vista mensual
- `week` - Vista semanal
- `day` - Vista diaria
- `agenda` - Próximos eventos
- `upcoming` - Eventos próximos

#### d) Gestión de Usuarios
- `user add` - Crear usuario
- `user list` - Listar usuarios
- `user show` - Ver usuario
- `user edit` - Editar usuario
- `user delete` - Eliminar usuario

#### e) Importar/Exportar
- `export` - Exportar a JSON/CSV/iCal
- `import` - Importar desde JSON/CSV/iCal

#### f) Estadísticas y Reportes
- `stats` - Estadísticas generales
- `summary` - Resumen de tiempo
- `freetime` - Bloques libres

#### g) Reportes para IA
- `daily-report` - Reporte diario completo
- `tomorrow-report` - Vista previa del día siguiente
- `upcoming-report` - Eventos próximos (horas)
- `weekly-report` - Resumen semanal

#### h) Recordatorios
- `reminder add` - Agregar recordatorio
- `reminder list` - Listar recordatorios
- `next` - Próximo evento

#### i) Configuración
- `config show` - Ver configuración
- `config set` - Configurar opciones

#### j) Utilidades
- `validate` - Validar integridad
- `cleanup` - Limpiar datos antiguos
- `backup` - Crear backup
- `restore` - Restaurar backup
- `info` - Información del sistema
- `version` - Versión

### 4. Eventos Recurrentes

Soporte completo para eventos que se repiten:
- Diarios
- Semanales (por día de semana)
- Mensuales (por día del mes o día de semana)
- Anuales
- Excepciones configurables

### 5. Reportes Optimizados para IA

Los reportes están diseñados para que una IA asista al usuario:

#### Reporte Diario (07:00 AM)
- Resumen del día completo
- Agenda cronológica
- Tareas pendientes extraídas
- Bloques de tiempo libre
- Vista previa del día siguiente
- Sugerencias de organización

#### Reporte de Mañana (20:00 PM)
- Vista previa del día siguiente
- Preparación necesaria
- Alertas de días pesados

#### Reporte de Próximos Eventos (cada hora)
- Eventos en las próximas 2 horas
- Solo eventos nuevos (no repetir)
- Recordatorios activos

#### Reporte Semanal (lunes 07:00 AM)
- Resumen de la semana
- Estadísticas
- Eventos importantes

## Arquitectura Técnica

### Stack Tecnológico

- **Lenguaje:** Go 1.24+
- **CLI Framework:** Cobra
- **Formato de datos:** Markdown + JSON
- **Almacenamiento:** Sistema de archivos

### Estructura del Proyecto

```
clical/
├── cmd/clical/              # Entry point
├── pkg/                     # Paquetes públicos
│   ├── calendar/            # Modelos Entry, validaciones
│   ├── storage/             # Interface + impl filesystem
│   ├── user/                # Gestión usuarios
│   ├── formatter/           # Salidas: text, json, csv, ical
│   ├── parser/              # Parsers datetime, duration
│   ├── view/                # Vistas: month, week, day
│   ├── reporter/            # Reportes para IA
│   └── importer/            # Import json, csv, ical
├── internal/                # Paquetes privados
│   ├── cli/                 # Comandos Cobra
│   ├── config/              # Configuración
│   └── util/                # Helpers
├── ai/specs/                # Especificaciones
├── ai/journal/              # Journal de desarrollo
└── docs/                    # Documentación web
```

### Formato de Almacenamiento

```
data/
└── users/
    └── {user_id}/
        ├── user.md                           # Info del usuario
        ├── user.json                         # Metadata
        ├── events/
        │   └── {year}/
        │       └── {month}/
        │           └── {day}/
        │               ├── {HH-MM-titulo}.md
        │               └── {HH-MM-titulo}.json
        ├── recurring/
        │   ├── daily/
        │   ├── weekly/{day}/
        │   ├── monthly/
        │   └── yearly/
        ├── tags/
        │   └── {tag}/
        │       └── index.md
        └── .state/
            ├── cron-state.json
            └── cron-reported.json
```

## Modelos de Datos

### Entry (Evento)

```go
type Entry struct {
    ID        string
    UserID    string
    DateTime  time.Time
    Title     string
    Duration  int  // minutos
    Location  string
    Notes     string
    Tags      []string
    Metadata  map[string]string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### RecurringEntry (Evento Recurrente)

```go
type RecurringEntry struct {
    ID      string
    UserID  string
    Pattern RecurrencePattern
    Title   string
    // ... resto de campos como Entry
}

type RecurrencePattern struct {
    Frequency   string  // daily, weekly, monthly, yearly
    Interval    int     // cada N días/semanas/meses
    DayOfWeek   int     // para weekly
    DayOfMonth  int     // para monthly
    StartDate   time.Time
    EndDate     *time.Time  // nil = indefinido
    Exceptions  []time.Time
}
```

### User

```go
type User struct {
    ID       string
    Name     string
    Timezone string
    Config   UserConfig
    Created  time.Time
}

type UserConfig struct {
    DefaultDuration int
    DateFormat      string
    TimeFormat      string
}
```

## Interfaces Principales

### Storage

```go
type Storage interface {
    // Entries
    SaveEntry(userID string, entry *Entry) error
    GetEntry(userID, entryID string) (*Entry, error)
    ListEntries(userID string, filter Filter) ([]*Entry, error)
    DeleteEntry(userID, entryID string) error
    UpdateEntry(userID string, entry *Entry) error

    // Recurring
    SaveRecurringEntry(userID string, entry *RecurringEntry) error
    ListRecurringEntries(userID string) ([]*RecurringEntry, error)

    // Users
    SaveUser(user *User) error
    GetUser(userID string) (*User, error)
    ListUsers() ([]*User, error)
    DeleteUser(userID string) error

    // State
    GetReportState(userID string) (*ReportState, error)
    SaveReportState(userID string, state *ReportState) error
}
```

### Formatter

```go
type Formatter interface {
    FormatEntry(entry *Entry) (string, error)
    FormatEntries(entries []*Entry) (string, error)
    FormatDailyReport(report *DailyReport) (string, error)
}
```

### Importer

```go
type Importer interface {
    Import(reader io.Reader) ([]*Entry, error)
}
```

## Casos de Uso Principales

### Uso Diario con Asistencia de IA

1. **07:00 AM** - Cron ejecuta `daily-report`
   - IA recibe agenda completa del día
   - IA saluda al usuario y presenta el día
   - IA identifica tareas pendientes
   - IA sugiere organización del día

2. **Durante el día** - Cron ejecuta `upcoming-report` cada hora
   - IA alerta eventos próximos
   - IA recuerda preparación necesaria
   - IA ofrece ayuda contextual

3. **20:00 PM** - Cron ejecuta `tomorrow-report`
   - IA presenta vista de mañana
   - IA sugiere preparación nocturna
   - IA alerta días pesados

4. **Usuario interactúa** - Comandos manuales
   - Usuario agrega eventos: `clical add`
   - Usuario consulta: `clical day`
   - Usuario busca: `clical search "reunión"`

### Gestión Manual

Usuario puede:
- Navegar archivos Markdown directamente
- Editar eventos con cualquier editor de texto
- Versionar con Git
- Hacer grep/search en archivos
- Backup simple copiando carpetas

## Características Avanzadas

### 1. Parser de Lenguaje Natural (Futuro)

```bash
clical add "Reunión con cliente mañana a las 2pm por 1 hora"
clical add "Llamada importante el viernes que viene a las 10"
```

### 2. Integración con Herramientas Externas

- Notificaciones del sistema
- Envío por email
- Integración con bots de Telegram
- Webhooks

### 3. Exportación Avanzada

- iCalendar (.ics) para importar en Google Calendar, Outlook, etc.
- CSV para Excel/Google Sheets
- JSON para procesamiento programático

### 4. Estadísticas

- Tiempo ocupado vs. libre
- Eventos por categoría/tag
- Patrones de productividad
- Heat maps de ocupación

## Requisitos No Funcionales

### Rendimiento

- Búsquedas rápidas incluso con miles de eventos
- Inicio instantáneo de comandos
- Índices para optimización

### Usabilidad

- Mensajes de error claros
- Ayuda contextual
- Ejemplos en --help
- Confirmaciones para operaciones destructivas

### Mantenibilidad

- Código modular
- Interfaces claras
- Tests unitarios
- Documentación completa

### Portabilidad

- Linux, macOS, Windows
- Sin dependencias externas complejas
- Binario único compilado

## Casos Edge

### 1. Zonas Horarias

- Soporte completo para timezones
- Conversión automática
- Storage en UTC + timezone del usuario

### 2. Eventos que Cruzan Medianoche

- Duración puede exceder 24 horas
- Vista correcta en múltiples días

### 3. Eventos Pasados

- Mantener histórico
- Cleanup opcional de eventos antiguos
- Búsqueda en histórico

### 4. Conflictos de Eventos

- Detección opcional de overlapping
- Advertencias al usuario
- No bloqueante (usuario decide)

## Fases de Implementación

### Fase 1: MVP (Mínimo Viable)
- Modelos de datos básicos
- Storage filesystem
- Comandos: add, list, show, delete
- Formato Markdown básico

### Fase 2: Vistas y Búsqueda
- Vistas: day, week, month
- Comandos: search, filter
- Formatters: text, json

### Fase 3: Reportes para IA
- daily-report
- tomorrow-report
- upcoming-report
- weekly-report

### Fase 4: Eventos Recurrentes
- Modelo RecurringEntry
- Storage de recurrentes
- Expansión a eventos individuales
- Excepciones

### Fase 5: Importar/Exportar
- Export: JSON, CSV, iCal
- Import: JSON, CSV
- Integración con calendarios externos

### Fase 6: Características Avanzadas
- Parser de lenguaje natural
- Estadísticas avanzadas
- TUI interactivo
- Recordatorios con notificaciones

## Referencias

- [Especificación de Comandos](./04_Commands.md)
- [Formato de Almacenamiento](./03_Storage_Format.md)
- [Arquitectura de Módulos](./02_Architecture.md)
- [Guías de Journaling](./01_Journaling.md)

---

**Estado:** Diseño completo - Pendiente de implementación
**Próximo paso:** Inicializar proyecto Go y estructura de directorios
