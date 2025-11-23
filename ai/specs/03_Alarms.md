# Sistema de Alarmas para IA

**Fecha:** 2025-11-23
**Tipo:** Feature Specification
**Estado:** En Diseño

## Descripción General

Sistema de alarmas que permite a la IA programar recordatorios para sí misma. Las alarmas pueden ser de una sola vez (one-time) o recurrentes (daily, weekly, monthly, yearly). Cuando una alarma se dispara, clical emite información en JSON que el caller puede redirigir a los canales apropiados (Telegram, agente IA, etc.).

## Filosofía de Diseño

**Separación de Responsabilidades (Unix Philosophy):**
- clical se encarga ÚNICAMENTE de gestionar calendario y alarmas
- clical NO conoce Telegram, agentes IA, ni canales de notificación
- clical emite JSON a stdout cuando hay alarmas
- El caller (script de cron) decide qué hacer con ese output

**Principio de Silencio:**
- Si no hay alarmas, `alarm-check` no emite output (exit 0)
- Solo con `--verbose` se muestra información de debugging

## Estructura de Directorios

```
data/alarms/<user_id>/
├── pending/              # Alarmas de una sola vez
│   ├── 2025-11-21_22-45-00.json
│   ├── 2025-11-21_23-00-00.json
│   └── 2025-11-22_10-30-00.json
│
├── recurring/
│   ├── daily/           # Formato: HH-MM-00.json
│   │   ├── 14-30-00.json
│   │   └── 09-00-00.json
│   │
│   ├── weekly/          # Formato: DAYNAME_HH-MM-00.json
│   │   ├── monday_14-30-00.json
│   │   ├── friday_17-00-00.json
│   │   └── sunday_20-00-00.json
│   │
│   ├── monthly/         # Formato: DD_HH-MM-00.json
│   │   ├── 01_09-00-00.json  # Día 1 de cada mes
│   │   ├── 15_14-30-00.json  # Día 15 de cada mes
│   │   └── 28_10-00-00.json
│   │
│   └── yearly/          # Formato: MM-DD_HH-MM-00.json
│       ├── 01-01_00-00-00.json  # 1 de enero
│       ├── 11-21_10-00-00.json  # 21 de noviembre
│       └── 12-25_08-00-00.json  # 25 de diciembre
│
└── past/
    ├── one-time/        # Alarmas ejecutadas (one-time)
    │   └── 2025-11-20_15-00-00.json
    │
    └── recurring/       # Alarmas recurrentes expiradas
        ├── daily/
        ├── weekly/
        ├── monthly/
        └── yearly/
```

## Formato de Archivos

### Archivo de Alarma

Cada archivo contiene un **array** de alarmas (puede haber múltiples alarmas para el mismo momento):

```json
[
  {
    "id": "alarm_001",
    "context": "Recordar revisar el PR de autenticación en repo clical",
    "created_at": "2025-11-21T22:30:00Z",
    "recurrence": "once"
  },
  {
    "id": "alarm_002",
    "context": "Seguimiento del deploy en producción",
    "created_at": "2025-11-21T22:35:00Z",
    "recurrence": "once"
  }
]
```

### Alarma Recurrente con Expiración

```json
[
  {
    "id": "alarm_weekly_001",
    "context": "Revisar métricas semanales del proyecto",
    "created_at": "2025-11-21T10:00:00Z",
    "recurrence": "weekly",
    "expires_at": "2025-12-31T23:59:59Z"
  }
]
```

### Campos del JSON

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `id` | string | Sí | Identificador único de la alarma |
| `context` | string | Sí | Información de contexto para la IA |
| `created_at` | ISO8601 | Sí | Cuándo se creó la alarma |
| `recurrence` | enum | Sí | Tipo: `once`, `daily`, `weekly`, `monthly`, `yearly` |
| `expires_at` | ISO8601 | No | Fecha de expiración (solo para recurrentes) |

## Comandos

### alarm-add (One-time)

Crear alarma de una sola vez:

```bash
clical alarm-add --at "2025-11-21 22:45" --context "Recordar revisar PR"

# Con user específico
clical alarm-add --user alice --at "2025-11-22 10:00" --context "Llamar al cliente"

# Formato flexible de fecha/hora
clical alarm-add --at "tomorrow 14:30" --context "Reunión de equipo"
clical alarm-add --at "+30m" --context "Revisar logs del deploy"
clical alarm-add --at "+2h" --context "Verificar métricas"
```

**Comportamiento:**
- Parsea la fecha/hora
- Genera ID único: `alarm_<timestamp>_<random>`
- Crea archivo en `pending/YYYY-MM-DD_HH-MM-00.json`
- Si el archivo ya existe, agrega al array

### alarm-add (Recurrente)

Crear alarmas recurrentes:

```bash
# Diaria
clical alarm-add --daily "14:30" --context "Revisar métricas diarias"
clical alarm-add --daily "09:00" --expires "2025-12-31" --context "Stand-up temporal"

# Semanal (días: monday, tuesday, wednesday, thursday, friday, saturday, sunday)
clical alarm-add --weekly monday "14:30" --context "Reunión semanal"
clical alarm-add --weekly friday "17:00" --context "Reporte semanal"

# Mensual (día del mes: 1-31)
clical alarm-add --monthly 1 "09:00" --context "Reporte mensual"
clical alarm-add --monthly 15 "14:30" --context "Revisión quincenal"

# Anual (formato: MM-DD)
clical alarm-add --yearly "01-01" "00:00" --context "Feliz año nuevo"
clical alarm-add --yearly "11-21" "10:00" --context "Aniversario del proyecto"
```

**Comportamiento:**
- Genera ID único: `alarm_<recurrence>_<timestamp>_<random>`
- Crea archivo en `recurring/<type>/` con formato apropiado
- Si tiene `--expires`, incluye `expires_at` en JSON

### alarm-check

Verifica alarmas pendientes y las ejecuta:

```bash
# Ejecutar desde cron cada minuto
* * * * * clical alarm-check --user alice | process-alarms.sh

# Con verbose para debugging
clical alarm-check --verbose

# Manual con user específico
clical alarm-check --user bob
```

**Comportamiento:**

1. **Chequeo de alarmas one-time (pending/):**
   - Construye timestamp actual: `YYYY-MM-DD_HH-MM-00`
   - Construye timestamps pasados (recovery): desde hace 60 minutos
   - Para cada timestamp (pasados + actual):
     - Si existe archivo `pending/<timestamp>.json`:
       - Lee el archivo
       - Agrega todas las alarmas al output
       - Mueve archivo a `past/one-time/`

2. **Chequeo de alarmas recurrentes:**

   **Daily:**
   - Construye filename: `HH-MM-00.json`
   - Si existe `recurring/daily/<filename>`:
     - Lee el archivo
     - Para cada alarma:
       - Si tiene `expires_at` y ya expiró:
         - Agrega al output
         - Mueve archivo a `past/recurring/daily/`
       - Si no expiró o no tiene expiración:
         - Agrega al output
         - NO mueve (sigue activa)

   **Weekly:**
   - Obtiene día de la semana actual (ej: "monday")
   - Construye filename: `<dayname>_HH-MM-00.json`
   - Mismo proceso que daily

   **Monthly:**
   - Obtiene día del mes actual (01-31)
   - Construye filename: `DD_HH-MM-00.json`
   - Mismo proceso que daily

   **Yearly:**
   - Obtiene mes-día actual (MM-DD)
   - Construye filename: `MM-DD_HH-MM-00.json`
   - Mismo proceso que daily

3. **Output:**
   - Si NO hay alarmas: exit 0, sin output
   - Si hay alarmas: imprime JSON array a stdout, exit 0
   - Con `--verbose`: imprime logs a stderr

**Output JSON:**

```json
[
  {
    "id": "alarm_001",
    "scheduled_for": "2025-11-21T22:45:00Z",
    "context": "Recordar revisar el PR",
    "created_at": "2025-11-21T22:30:00Z",
    "recurrence": "once"
  },
  {
    "id": "alarm_weekly_001",
    "scheduled_for": "2025-11-21T14:30:00Z",
    "context": "Revisar métricas semanales",
    "created_at": "2025-11-21T10:00:00Z",
    "recurrence": "weekly",
    "expires_at": "2025-12-31T23:59:59Z"
  }
]
```

### alarm-list

Lista alarmas activas:

```bash
# Listar todas las alarmas activas (formato tabla)
clical alarm-list

# Incluir alarmas pasadas
clical alarm-list --past

# Output en JSON (para scripting)
clical alarm-list --json

# Con user específico
clical alarm-list --user alice --json
```

**Output (tabla por defecto):**

```
ID              TYPE      SCHEDULED               CONTEXT
alarm_001       once      2025-11-21 22:45:00     Recordar revisar el PR
alarm_weekly_01 weekly    monday 14:30:00         Reunión semanal
alarm_daily_02  daily     14:30:00                Revisar métricas
```

**Output (con --json):**

```json
{
  "active": [
    {
      "id": "alarm_001",
      "scheduled_for": "2025-11-21T22:45:00Z",
      "context": "Recordar revisar el PR",
      "created_at": "2025-11-21T22:30:00Z",
      "recurrence": "once"
    },
    {
      "id": "alarm_weekly_001",
      "scheduled_for": "weekly monday 14:30:00",
      "context": "Revisar métricas semanales",
      "created_at": "2025-11-21T10:00:00Z",
      "recurrence": "weekly"
    }
  ],
  "past": []
}
```

**Con --past:**

```json
{
  "active": [...],
  "past": [
    {
      "id": "alarm_002",
      "scheduled_for": "2025-11-20T15:00:00Z",
      "context": "Deploy completado",
      "created_at": "2025-11-20T14:00:00Z",
      "recurrence": "once",
      "executed_at": "2025-11-20T15:00:00Z"
    }
  ]
}
```

### alarm-cancel

Cancelar una alarma activa:

```bash
# Cancelar por ID
clical alarm-cancel alarm_001

# Con user específico
clical alarm-cancel --user alice alarm_weekly_001
```

**Comportamiento:**
- Busca el archivo que contiene esa alarma en `pending/` o `recurring/`
- Si el archivo tiene solo esa alarma: mueve a `past/` (con subfolder apropiado)
- Si el archivo tiene múltiples alarmas: remueve del array y reescribe archivo
- Si no encuentra la alarma: error "Alarm not found"

## Integración con User

Las alarmas están asociadas al mismo `--user` que los eventos del calendario:

```bash
# Crear alarma para usuario alice
clical alarm-add --user alice --at "tomorrow 10:00" --context "Revisar proyecto"

# Chequear alarmas de alice
clical alarm-check --user alice

# Listar alarmas de alice
clical alarm-list --user alice
```

**Estructura:**
```
data/
├── events/
│   └── alice/
│       └── 2025/11/21/...
└── alarms/
    └── alice/
        ├── pending/
        ├── recurring/
        └── past/
```

## Caso de Uso: Integración con Cron

**Script de cron (`/usr/local/bin/clical-alarm-processor`):**

```bash
#!/bin/bash
OUTPUT=$(clical alarm-check --user ai-assistant 2>/dev/null)

if [ -n "$OUTPUT" ]; then
  # Enviar a Telegram
  echo "$OUTPUT" | jq -r '.[] | .context' | while read -r ctx; do
    s gobot-send-message "[ALARMA] $ctx"
  done

  # Notificar a agente IA
  echo "$OUTPUT" | ai-agent-notify --stdin
fi
```

**Crontab:**

```cron
# Ejecutar cada minuto
* * * * * /usr/local/bin/clical-alarm-processor
```

## Arquitectura de Código

### Paquetes

```
pkg/alarm/
├── alarm.go           # Model: Alarm struct
├── storage.go         # Interface + filesystem implementation
└── recurrence.go      # Tipos y helpers de recurrencia

internal/cli/
└── alarm.go           # Comandos: alarm-add, alarm-check, alarm-list, alarm-cancel
```

### Model: Alarm

```go
package alarm

import "time"

type Recurrence string

const (
    RecurrenceOnce    Recurrence = "once"
    RecurrenceDaily   Recurrence = "daily"
    RecurrenceWeekly  Recurrence = "weekly"
    RecurrenceMonthly Recurrence = "monthly"
    RecurrenceYearly  Recurrence = "yearly"
)

type Alarm struct {
    ID          string      `json:"id"`
    Context     string      `json:"context"`
    CreatedAt   time.Time   `json:"created_at"`
    Recurrence  Recurrence  `json:"recurrence"`
    ExpiresAt   *time.Time  `json:"expires_at,omitempty"`

    // Campos adicionales para output de alarm-check
    ScheduledFor time.Time  `json:"scheduled_for,omitempty"`
    ExecutedAt   *time.Time `json:"executed_at,omitempty"`
}

func NewAlarm(context string, recurrence Recurrence) *Alarm
func (a *Alarm) Validate() error
func (a *Alarm) IsExpired() bool
```

### Storage Interface

```go
package alarm

type Storage interface {
    // One-time alarms
    SaveOneTime(userID string, when time.Time, alarm *Alarm) error
    GetOneTime(userID string, when time.Time) ([]*Alarm, error)
    DeleteOneTime(userID string, when time.Time) error

    // Recurring alarms
    SaveRecurring(userID string, recurrence Recurrence, filename string, alarm *Alarm) error
    GetRecurring(userID string, recurrence Recurrence, filename string) ([]*Alarm, error)
    DeleteRecurring(userID string, recurrence Recurrence, filename string) error

    // Check alarms
    CheckAlarms(userID string, at time.Time) ([]*Alarm, error)

    // List alarms
    ListActive(userID string) ([]*Alarm, error)
    ListPast(userID string) ([]*Alarm, error)

    // Cancel alarm
    CancelAlarm(userID string, alarmID string) error

    // Move to past
    MoveToOnePast(userID string, when time.Time) error
    MoveToRecurringPast(userID string, recurrence Recurrence, filename string) error
}
```

## Validaciones

### alarm-add

- `--at` debe ser fecha/hora válida y futura
- `--daily`, `--weekly`, etc.: formato de hora válido
- `--weekly`: día de semana válido (monday-sunday)
- `--monthly`: día del mes válido (1-31)
- `--yearly`: fecha válida (MM-DD)
- `--expires`: fecha futura (solo para recurrentes)
- `--context`: no vacío, máximo 500 caracteres

### alarm-cancel

- ID debe existir
- Solo se pueden cancelar alarmas activas (no pasadas)

## Colisiones

**Múltiples alarmas al mismo tiempo:**
- ✅ Permitido
- Razón: contextos diferentes, la IA puede necesitar múltiples recordatorios simultáneos
- El array en el JSON soporta esto naturalmente

**Ejemplo:**
```json
// pending/2025-11-21_14-30-00.json
[
  {"id": "alarm_001", "context": "Revisar PR"},
  {"id": "alarm_002", "context": "Llamar cliente"},
  {"id": "alarm_003", "context": "Deploy a producción"}
]
```

## Recovery

**Alarmas perdidas (sistema apagado, cron no ejecutado):**

El comando `alarm-check` incluye recovery automático:
- Revisa últimos 60 minutos de alarmas pendientes
- Ejecuta todas las alarmas perdidas
- Las mueve a `past/`

**Nota:** No hay límite en cuántas alarmas perdidas recuperar. Si el sistema estuvo apagado días, ejecutará todas las one-time pendientes.

## Consideraciones Futuras

**No implementar ahora, pero tener en cuenta:**

1. **Alarmas con offset:**
   - Ej: "30 minutos antes del evento X"
   - Requiere integración con pkg/calendar

2. **Alarmas condicionales:**
   - Ej: "Si el PR no está mergeado a las 17:00, recordar"
   - Requiere integración con APIs externas

3. **Snooze:**
   - Posponer alarma X minutos
   - `clical alarm-snooze alarm_001 30m`

4. **Prioridades:**
   - Alarmas urgentes vs normales
   - Afecta el formato del mensaje

## Testing

**Unit tests requeridos:**
- `pkg/alarm/alarm_test.go`: Validación de Alarm struct
- `pkg/alarm/storage_test.go`: Operaciones de storage
- `pkg/alarm/recurrence_test.go`: Cálculo de fechas recurrentes

**Integration tests:**
- Crear alarma → alarm-check → verificar output
- Alarmas múltiples en mismo minuto
- Recovery de alarmas perdidas
- Expiración de alarmas recurrentes
- Cancelación de alarmas

**Manual testing:**
- Configurar cron con usuario de prueba
- Verificar integración con Telegram
- Verificar que no hay output cuando no hay alarmas

---

**Estado:** ✅ Especificación completa
**Próximo paso:** Implementación
**Fecha:** 2025-11-23
