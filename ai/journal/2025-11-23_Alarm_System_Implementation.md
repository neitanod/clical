# Implementación del Sistema de Alarmas

**Fecha:** 2025-11-23
**Tipo:** Feature Implementation
**Estado:** Completado

## Objetivo

Implementar un sistema completo de alarmas que permita a la IA programar recordatorios y seguimientos automáticos, tanto one-time como recurrentes (daily, weekly, monthly, yearly).

## Contexto

El usuario solicitó un sistema donde la IA pueda "ponerse recordatorios para sí misma" y hacer seguimiento de tareas programadas. El diseño siguió la filosofía Unix de separación de responsabilidades: clical gestiona alarmas, el caller decide qué hacer con el output.

## Diseño

### Filosofía

1. **Separación de responsabilidades:** clical NO conoce Telegram, IA, ni canales de notificación
2. **Principio de silencio:** Si no hay alarmas, `alarm check` no emite output
3. **Unix philosophy:** clical hace una cosa bien = gestionar calendario y alarmas
4. **JSON a stdout:** El caller procesa el output según necesite

### Arquitectura de Archivos

```
data/users/<user_id>/alarms/
├── pending/                    # One-time alarms
│   └── YYYY-MM-DD_HH-MM-00.json
├── recurring/
│   ├── daily/                  # HH-MM-00.json
│   ├── weekly/                 # dayname_HH-MM-00.json
│   ├── monthly/                # DD_HH-MM-00.json
│   └── yearly/                 # MM-DD_HH-MM-00.json
└── past/
    ├── one-time/
    └── recurring/
```

**Optimización clave:** Filename = timestamp → verificación ultra-rápida con simple file existence check.

### Tipos de Recurrencia

- `once` - One-time alarm (se ejecuta una vez y se mueve a past/)
- `daily` - Cada día a la hora especificada
- `weekly` - Cada semana en el día especificado
- `monthly` - Cada mes en el día especificado (1-31)
- `yearly` - Cada año en la fecha especificada (MM-DD)

## Implementación

### 1. Paquete `pkg/alarm/`

**`recurrence.go` (156 líneas):**
- Tipos: `Recurrence`, `DailySchedule`, `WeeklySchedule`, `MonthlySchedule`, `YearlySchedule`
- Helpers: `ParseWeekday`, `OneTimeFilename`, `Current*Filename`, `RoundToMinute`
- Generación de filenames según tipo de recurrencia

**`alarm.go` (128 líneas):**
- Model `Alarm` con campos: ID, Context, CreatedAt, Recurrence, ExpiresAt, ScheduledFor, ExecutedAt
- `NewAlarm` con generación automática de ID único
- Validación completa (context no vacío, recurrence válido, expires_at solo para recurrentes)
- `IsExpired`, `ShouldExecute`, `WithScheduledFor`, `WithExecutedAt`, `Clone`

### 2. Storage Layer

**`pkg/storage/alarm_paths.go` (87 líneas):**
- Helper `AlarmPaths` para gestionar rutas
- `PendingDir`, `RecurringDir`, `PastDir`, `PendingFile`, `RecurringFile`, `PastFile`
- `EnsureAlarmDirs` crea toda la estructura necesaria

**`pkg/storage/filesystem_alarms.go` (422 líneas):**
- `SaveAlarm`: Guarda alarma agregándola al array del archivo
- `GetAlarms`: Lee todas las alarmas de un archivo
- `DeleteAlarms`: Elimina un archivo de alarmas
- `CheckAlarms`: **Función principal** - verifica alarmas pendientes
  - Recovery automático: últimos 60 minutos de one-time alarms
  - Verifica todas las recurrencias que coincidan con el momento actual
  - Maneja expiración de alarmas recurrentes
  - Mueve alarmas ejecutadas a `past/`
- `ListActiveAlarms`: Lista todas las alarmas activas
- `ListPastAlarms`: Lista histórico de alarmas
- `CancelAlarm`: Cancela una alarma por ID

**Lógica de CheckAlarms:**
1. Verifica one-time alarms (pending/) para los últimos 60 minutos (recovery)
2. Verifica cada tipo de recurring alarm que coincida con el minuto actual
3. Para alarmas con `expires_at` que expiraron: ejecuta y mueve a past/
4. Para alarmas no expiradas: ejecuta pero NO mueve (siguen activas)
5. Retorna array de alarmas para ejecutar

### 3. CLI Commands

**`internal/cli/alarm.go` (677 líneas):**

**Comando `alarm add`:**
- Flags: `--at`, `--daily`, `--weekly`, `--monthly`, `--yearly`, `--context`, `--expires`
- Funciones: `addOneTimeAlarm`, `addDailyAlarm`, `addWeeklyAlarm`, `addMonthlyAlarm`, `addYearlyAlarm`
- Validación completa de inputs
- Redondeo automático a minuto para one-time alarms
- Parsing de formatos específicos por tipo (ej: "monday 14:30", "15 14:30", "11-21 10:00")

**Comando `alarm check`:**
- Flag: `--verbose`
- Si NO hay alarmas: exit 0, sin output
- Si hay alarmas: JSON array a stdout
- Ideal para cron: `* * * * * clical alarm check --user X | process.sh`

**Comando `alarm list`:**
- Flags: `--past`, `--json`
- Por defecto: tabla formateada de alarmas activas
- Con `--json`: output estructurado para scripting
- Con `--past`: incluye histórico

**Comando `alarm cancel`:**
- Args: `ALARM_ID`
- Busca la alarma en pending/ y recurring/
- Si hay múltiples alarmas en un archivo, solo remueve la especificada
- Si es la única, elimina el archivo completo

### 4. Interface Storage

Extendimos `pkg/storage/storage.go` con 8 nuevos métodos:
```go
SaveAlarm(userID, alarmTime, recurrence, filename, alarm)
GetAlarms(userID, recurrence, filename)
DeleteAlarms(userID, recurrence, filename)
CheckAlarms(userID, at time.Time) ([]*Alarm, error)
ListActiveAlarms(userID)
ListPastAlarms(userID)
CancelAlarm(userID, alarmID)
MoveAlarmsToPast(userID, recurrence, filename)
```

### 5. Tests

**`pkg/alarm/alarm_test.go` (327 líneas):**
- `TestNewAlarm`: Creación y fields correctos
- `TestValidate`: 9 casos de validación
- `TestIsExpired`: Manejo de expiración
- `TestShouldExecute`: Lógica de ejecución
- `TestWithScheduledFor`, `TestWithExecutedAt`: Mutadores
- `TestClone`: Deep copy correcto

**`pkg/alarm/recurrence_test.go` (321 líneas):**
- `TestRecurrenceValid`: Validación de tipos
- `TestWeeklyScheduleFilename`, `TestMonthlyScheduleFilename`, `TestYearlyScheduleFilename`, `TestDailyScheduleFilename`
- `TestParseWeekday`: 12 casos (incluyendo uppercase, mixed case, espacios, inválidos)
- `TestOneTimeFilename`: Generación de filenames
- `TestCurrentFilenames`: Helpers de filename actual
- `TestRoundToMinute`: Redondeo a minuto

**Resultado:** 100% cobertura en `pkg/alarm/`

## Cambios en Archivos

### Archivos Nuevos

1. `pkg/alarm/alarm.go` - Model de alarma
2. `pkg/alarm/recurrence.go` - Tipos y helpers de recurrencia
3. `pkg/alarm/alarm_test.go` - Tests del model
4. `pkg/alarm/recurrence_test.go` - Tests de recurrencia
5. `pkg/storage/alarm_paths.go` - Path helpers
6. `pkg/storage/filesystem_alarms.go` - Implementación de storage
7. `internal/cli/alarm.go` - Comandos CLI
8. `ai/specs/03_Alarms.md` - Especificación completa (470+ líneas)

### Archivos Modificados

1. `pkg/storage/storage.go` - Agregados 8 métodos a interface Storage
2. `internal/cli/root.go` - Registrado comando `alarmCmd`
3. `docs/USAGE.md` - Sección 9: Sistema de Alarmas (250+ líneas)
4. `docs/index.html` - Sección de alarmas en navegación y contenido

## Comandos Implementados

### alarm add

```bash
# One-time
clical alarm add --user ai-agent --at "2025-11-24 10:00" --context "..."

# Daily
clical alarm add --user ai-agent --daily "14:30" --context "..."
clical alarm add --user ai-agent --daily "09:00" --expires "2025-12-31" --context "..."

# Weekly
clical alarm add --user ai-agent --weekly "monday 14:30" --context "..."

# Monthly
clical alarm add --user ai-agent --monthly "1 09:00" --context "..."

# Yearly
clical alarm add --user ai-agent --yearly "01-01 00:00" --context "..."
```

### alarm check

```bash
# Para cron (cada minuto)
clical alarm check --user ai-agent

# Con verbose
clical alarm check --user ai-agent --verbose
```

**Output JSON:**
```json
[
  {
    "id": "alarm_once_1234567890_abcd1234",
    "scheduled_for": "2025-11-24T10:00:00Z",
    "context": "Revisar PR de autenticación",
    "created_at": "2025-11-23T14:00:00Z",
    "recurrence": "once"
  }
]
```

### alarm list

```bash
clical alarm list --user ai-agent
clical alarm list --user ai-agent --past
clical alarm list --user ai-agent --json
```

### alarm cancel

```bash
clical alarm cancel --user ai-agent alarm_once_1234567890_abcd1234
```

## Integración con Cron

**Configuración recomendada:**

```bash
# Ejecutar cada minuto
* * * * * /usr/local/bin/clical-alarm-processor
```

**Script procesador:**
```bash
#!/bin/bash
USER="ai-agent"
OUTPUT=$(clical alarm check --user "$USER" 2>/dev/null)

if [ -n "$OUTPUT" ]; then
  # Enviar a Telegram
  echo "$OUTPUT" | jq -r '.[] | .context' | while read -r ctx; do
    s gobot-send-message "[ALARMA] $ctx"
  done

  # Notificar a agente IA
  echo "$OUTPUT" | ai-agent-notify --stdin
fi
```

## Testing Manual

```bash
# Crear alarma
$ clical alarm add --user test --at "2025-11-24 10:00" --context "Test"
✓ Alarma creada exitosamente
ID: alarm_once_1763918422_b17f3770

# Listar
$ clical alarm list --user test
ALARMAS ACTIVAS:
ID                              TIPO    PROGRAMADA          CONTEXTO
alarm_once_1763918422_b17f3770  once    pending             Test

# Check (no hay alarmas para este minuto)
$ clical alarm check --user test
$ echo $?
0

# Cancelar
$ clical alarm cancel --user test alarm_once_1763918422_b17f3770
✓ Alarma cancelada exitosamente
```

## Características Clave

1. **Recovery Automático:** `alarm check` recupera alarmas perdidas de los últimos 60 minutos
2. **Múltiples alarmas por minuto:** Los archivos contienen arrays, soportan múltiples alarmas
3. **Expiración para recurrentes:** Las alarmas recurrentes pueden tener `expires_at`
4. **Silencio por defecto:** Sin alarmas = sin output (perfecto para cron)
5. **JSON estructurado:** Output parseable para integración con scripts
6. **IDs únicos:** Cada alarma tiene ID único con timestamp + random
7. **Persistencia simple:** Archivos JSON legibles por humanos
8. **Verificación eficiente:** Filename = timestamp → solo check file existence

## Casos de Uso Implementados

### Caso 1: Seguimiento de Tareas
```
Usuario: "Recordame en 30 minutos revisar el deploy"
IA: clical alarm add --at "2025-11-24 14:45" --context "Revisar deploy"
```

### Caso 2: Reportes Automáticos
```
Usuario: "Cada viernes a las 5 PM dame un resumen"
IA: clical alarm add --weekly "friday 17:00" --context "Generar resumen semanal"
```

### Caso 3: Alarmas Temporales
```
Usuario: "Durante diciembre, recordame cada día revisar métricas"
IA: clical alarm add --daily "09:00" --expires "2025-12-31" --context "Revisar métricas"
```

## Decisiones de Diseño

### ¿Por qué archivos separados por minuto?

Verificación ultrarrápida: solo verificar si existe archivo para el minuto actual. No parsear JSON hasta estar seguro que hay alarmas.

### ¿Por qué JSON arrays?

Soporta múltiples alarmas al mismo tiempo naturalmente. La IA puede programar varios recordatorios para el mismo minuto.

### ¿Por qué no background daemon?

- Más robusto: se autorestaura después de reinicios
- Más simple: no necesita gestión de procesos
- Cron cada minuto es suficiente para este caso de uso

### ¿Por qué no borrar alarmas pasadas?

Mantener historial es útil para debugging y auditoría. Los archivos en `past/` se pueden limpiar manualmente si es necesario.

## Documentación

1. **ai/specs/03_Alarms.md** - Especificación técnica completa (470 líneas)
2. **docs/USAGE.md** - Guía de uso para agentes IA (250+ líneas nuevas)
3. **docs/index.html** - Documentación web con ejemplos interactivos

## Métricas

- **Archivos nuevos:** 8
- **Archivos modificados:** 4
- **Líneas de código:** ~2,100
- **Líneas de tests:** ~650
- **Cobertura:** 100% en pkg/alarm/
- **Comandos:** 4 (add, check, list, cancel)
- **Tipos de recurrencia:** 5 (once, daily, weekly, monthly, yearly)

## Próximos Pasos Sugeridos

- [ ] Tests de integración end-to-end
- [ ] Tests de storage (filesystem_alarms.go)
- [ ] Parser de lenguaje natural para fechas ("en 30 minutos", "mañana a las 10")
- [ ] Snooze de alarmas
- [ ] Prioridades (urgente, normal, low)
- [ ] Alarmas condicionales ("si el PR no está mergeado...")

## Aprendizajes

1. **Filename como clave:** Usar el filename como timestamp fue clave para optimización
2. **Separación de responsabilidades:** NO integrar Telegram/IA directamente hizo el sistema más flexible
3. **Recovery automático:** Incluir recovery desde el principio evitó problemas de alarmas perdidas
4. **JSON arrays:** Decisión simple que resolvió el caso de múltiples alarmas elegantemente

---

**Estado:** ✅ Completado e instalado
**Compilación:** Exitosa
**Tests:** 100% passing (pkg/alarm/)
**Documentación:** Completa
**Fecha:** 2025-11-23
