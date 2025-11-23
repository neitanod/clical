# clical - Guía de Uso para Agentes IA

Este documento está diseñado para que agentes de IA aprendan a usar **clical** de manera efectiva para asistir a usuarios en la gestión de su calendario.

## Conceptos Fundamentales

### Filosofía de clical

clical es un calendario diseñado para ser consultado por una IA, no por el usuario directamente. La IA usa clical para:

1. **Conocer la agenda del usuario** - Qué eventos tiene hoy, mañana, esta semana
2. **Recordar eventos próximos** - Alertar sobre eventos que están por comenzar
3. **Sugerir organización** - Identificar bloques libres, preparación necesaria
4. **Asistir proactivamente** - Ayudar al usuario a organizarse sin que lo pida

### Modelo de Datos

Un **evento** tiene:
- `id` - Identificador único (generado automáticamente)
- `user_id` - ID del usuario dueño del evento
- `datetime` - Fecha y hora de inicio (formato: "YYYY-MM-DD HH:MM")
- `title` - Título descriptivo
- `duration` - Duración en minutos
- `location` - Ubicación (opcional)
- `notes` - Notas adicionales (opcional)
- `tags` - Etiquetas para categorizar (opcional)

Un **usuario** tiene:
- `id` - Identificador único
- `name` - Nombre del usuario
- `timezone` - Zona horaria (ej: "America/Argentina/Buenos_Aires")

## Comandos Principales

### 1. Gestión de Usuarios

#### user add - Crear usuario

```bash
clical user add --id=USER_ID --name="NOMBRE" --timezone="TIMEZONE"
```

**Argumentos:**
- `--id` (requerido) - ID único del usuario (ej: número de Telegram)
- `--name` (requerido) - Nombre del usuario
- `--timezone` (requerido) - Zona horaria válida

**Ejemplo:**
```bash
clical user add --id=123456789 --name="Juan Pérez" --timezone="America/Argentina/Buenos_Aires"
```

**Cuándo usar:**
- Al configurar clical por primera vez para un usuario
- Cuando un nuevo usuario quiere empezar a usar el calendario

#### user list - Listar usuarios

```bash
clical user list
```

**Argumentos:** Ninguno

**Salida:** Lista de todos los usuarios registrados con su ID, nombre y timezone

**Cuándo usar:**
- Para ver qué usuarios están registrados
- Para obtener IDs de usuarios disponibles

#### user show - Ver detalles de usuario

```bash
clical user show --id=USER_ID
```

**Argumentos:**
- `--id` (requerido) - ID del usuario

**Ejemplo:**
```bash
clical user show --id=123456789
```

**Cuándo usar:**
- Para ver configuración de un usuario específico
- Para verificar timezone y preferencias

---

### 2. Gestión de Eventos

#### add - Agregar evento

```bash
clical add --user=USER_ID --datetime="YYYY-MM-DD HH:MM" --title="TÍTULO" [opciones]
```

**Argumentos requeridos:**
- `--user` - ID del usuario
- `--datetime` - Fecha y hora (formato: "YYYY-MM-DD HH:MM")
- `--title` - Título del evento

**Argumentos opcionales:**
- `--duration` - Duración en minutos (default: 60)
- `--location` - Ubicación del evento
- `--notes` - Notas adicionales
- `--tags` - Tags separados por coma

**Ejemplos:**

```bash
# Evento simple
clical add --user=123456789 \
  --datetime="2025-11-21 09:00" \
  --title="Stand-up Meeting" \
  --duration=15

# Evento completo
clical add --user=123456789 \
  --datetime="2025-11-21 14:00" \
  --title="Reunión con cliente" \
  --duration=60 \
  --location="Oficina Central, Sala 3" \
  --notes="Revisar propuesta Q4. Llevar laptop y documentos impresos." \
  --tags=trabajo,cliente,importante

# Evento con múltiples tags
clical add --user=123456789 \
  --datetime="2025-11-22 11:00" \
  --title="Code Review Feature X" \
  --duration=45 \
  --tags=desarrollo,revision,feature-x
```

**Cuándo usar:**
- Cuando el usuario menciona un evento futuro
- Cuando pide agendar algo
- Al planificar la semana/mes

**Tips para IA:**
- Extraer la fecha/hora de lenguaje natural del usuario
- Inferir duración típica si no se especifica (meetings: 30-60 min, llamadas: 15-30 min)
- Agregar tags relevantes para facilitar búsqueda posterior
- Incluir en notes cualquier preparación necesaria

#### list - Listar eventos

```bash
clical list --user=USER_ID [filtros]
```

**Argumentos:**
- `--user` (requerido) - ID del usuario

**Filtros opcionales:**
- `--from="YYYY-MM-DD"` - Desde esta fecha
- `--to="YYYY-MM-DD"` - Hasta esta fecha
- `--range=RANGO` - Rango predefinido: "today", "week", "month"
- `--tags=TAG1,TAG2` - Filtrar por tags

**Ejemplos:**

```bash
# Todos los eventos del usuario
clical list --user=123456789

# Eventos de hoy
clical list --user=123456789 --range=today

# Eventos de esta semana
clical list --user=123456789 --range=week

# Eventos de este mes
clical list --user=123456789 --range=month

# Rango personalizado
clical list --user=123456789 --from="2025-11-20" --to="2025-11-30"

# Filtrar por tags
clical list --user=123456789 --tags=trabajo
clical list --user=123456789 --tags=cliente,importante

# Eventos de trabajo esta semana
clical list --user=123456789 --range=week --tags=trabajo
```

**Cuándo usar:**
- Para conocer la agenda completa del usuario
- Antes de sugerir agregar un evento (verificar conflictos)
- Cuando el usuario pregunta "qué tengo hoy/mañana/esta semana"

**Tips para IA:**
- Usar `--range=today` frecuentemente para conocer agenda del día
- Combinar filtros para búsquedas específicas
- Ordenar resultados por fecha al presentarlos al usuario

#### show - Ver evento específico

```bash
clical show --user=USER_ID --id=EVENT_ID
```

**Argumentos:**
- `--user` (requerido) - ID del usuario
- `--id` (requerido) - ID del evento

**Ejemplo:**
```bash
clical show --user=123456789 --id=abc123def456
```

**Cuándo usar:**
- Para ver todos los detalles de un evento (notas, metadata)
- Cuando el usuario pide información sobre un evento específico
- Para verificar datos antes de editar

#### edit - Editar evento

```bash
clical edit --user=USER_ID --id=EVENT_ID [cambios]
```

**Argumentos:**
- `--user` (requerido) - ID del usuario
- `--id` (requerido) - ID del evento a editar

**Campos editables:**
- `--title="NUEVO_TÍTULO"`
- `--datetime="YYYY-MM-DD HH:MM"`
- `--duration=MINUTOS`
- `--location="NUEVA_UBICACIÓN"`
- `--notes="NUEVAS_NOTAS"`

**Ejemplos:**

```bash
# Cambiar título
clical edit --user=123456789 --id=abc123 --title="Reunión Reprogramada"

# Cambiar fecha y hora
clical edit --user=123456789 --id=abc123 --datetime="2025-11-22 15:00"

# Cambiar duración
clical edit --user=123456789 --id=abc123 --duration=90

# Múltiples cambios
clical edit --user=123456789 --id=abc123 \
  --title="Reunión Virtual" \
  --location="Zoom" \
  --duration=45
```

**Cuándo usar:**
- Cuando el usuario pide reprogramar un evento
- Para actualizar detalles de un evento existente
- Cuando cambian circunstancias (ubicación, duración)

**Tips para IA:**
- Primero usar `list` o `show` para obtener el ID del evento
- Confirmar cambios con el usuario antes de ejecutar
- Solo editar los campos que cambian (no es necesario especificar todos)

#### delete - Eliminar evento

```bash
clical delete --user=USER_ID --id=EVENT_ID [--force]
```

**Argumentos:**
- `--user` (requerido) - ID del usuario
- `--id` (requerido) - ID del evento a eliminar
- `--force` (opcional) - Eliminar sin confirmación

**Ejemplos:**

```bash
# Con confirmación interactiva
clical delete --user=123456789 --id=abc123

# Sin confirmación (usar con precaución)
clical delete --user=123456789 --id=abc123 --force
```

**Cuándo usar:**
- Cuando el usuario cancela un evento
- Para limpiar eventos obsoletos

**Tips para IA:**
- Usar `--force` en automatizaciones, no en interacciones directas
- Confirmar con el usuario antes de eliminar
- Mostrar detalles del evento que se va a eliminar

---

### 3. Reportes para IA

Estos comandos están diseñados específicamente para que una IA obtenga información estructurada del calendario.

#### daily-report - Reporte diario completo

```bash
clical daily-report --user=USER_ID [--date="YYYY-MM-DD"]
```

**Argumentos:**
- `--user` (requerido) - ID del usuario
- `--date` (opcional) - Fecha específica (default: hoy)

**Ejemplos:**

```bash
# Reporte de hoy
clical daily-report --user=123456789

# Reporte de una fecha específica
clical daily-report --user=123456789 --date="2025-11-21"
```

**Contenido del reporte:**
- Resumen del día (total eventos, horas ocupadas, tiempo libre)
- Próximo evento inmediato (con minutos restantes)
- Agenda completa cronológica
- Bloques de tiempo libre con sugerencias de uso
- Vista previa del día siguiente
- Sugerencias de organización

**Cuándo usar:**
- **07:00 AM** - Saludo matutino con agenda del día
- Cuando el usuario pregunta "qué tengo hoy"
- Al planificar el día
- Antes de sugerir agregar eventos (ver disponibilidad)

**Tips para IA:**
- Ejecutar automáticamente cada mañana
- Usar para conocer el contexto del día
- Presentar al usuario de forma conversacional
- Identificar eventos que requieren preparación

**Ejemplo de uso por IA:**

```
Usuario: "Buenos días"
IA: [Ejecuta: clical daily-report --user=123456789]
IA: "¡Buenos días! Hoy es Viernes 21 de Noviembre. Tienes 3 eventos:
     - 09:00 Stand-up Meeting (15 min)
     - 11:00 Desarrollo Feature X (2 horas)
     - 15:00 Code Review (45 min)

     Tu primer bloque libre largo es de 09:15 a 11:00 (1h 45min),
     ideal para trabajo enfocado.

     ¿Necesitas que te prepare algo para tus eventos?"
```

#### tomorrow-report - Reporte del día siguiente

```bash
clical tomorrow-report --user=USER_ID
```

**Argumentos:**
- `--user` (requerido) - ID del usuario

**Ejemplo:**
```bash
clical tomorrow-report --user=123456789
```

**Contenido:**
- Igual que daily-report pero para el día siguiente

**Cuándo usar:**
- **20:00 PM** - Al final del día
- Cuando el usuario pregunta "qué tengo mañana"
- Para planificación nocturna

**Tips para IA:**
- Ejecutar automáticamente cada noche
- Alertar si mañana hay día pesado
- Sugerir preparación nocturna si es necesario

**Ejemplo de uso por IA:**

```
[20:00 PM - Ejecuta automáticamente]
IA: "Buenas noches! Vista previa de mañana Sábado 22:

     Tienes 2 eventos:
     - 10:00 Entrevista técnica (90 min) ⭐
     - 15:00 Sprint Planning (2 horas)

     ⚠️ Día moderado: 3.5 horas de eventos.

     Considera revisar el backlog esta noche para el Sprint Planning."
```

#### upcoming-report - Próximos eventos

```bash
clical upcoming-report --user=USER_ID [--hours=N] [--count=N]
```

**Argumentos:**
- `--user` (requerido) - ID del usuario
- `--hours=N` (opcional) - Próximas N horas (default: 2)
- `--count=N` (opcional) - Próximos N eventos (sobrescribe --hours)

**Ejemplos:**

```bash
# Próximos eventos en las siguientes 2 horas
clical upcoming-report --user=123456789

# Próximas 4 horas
clical upcoming-report --user=123456789 --hours=4

# Próximos 5 eventos (sin límite de tiempo)
clical upcoming-report --user=123456789 --count=5
```

**Salida:** Lista de eventos próximos con tiempo restante

**Cuándo usar:**
- **Cada hora** durante horario laboral (09:00-18:00)
- Antes de eventos importantes (recordatorio 15-30 min antes)
- Cuando el usuario pregunta "qué tengo próximamente"

**Tips para IA:**
- Ejecutar periódicamente en background
- Alertar solo si hay eventos en las próximas 2 horas
- Incluir recordatorios de preparación necesaria

**Ejemplo de uso por IA:**

```
[Ejecuta cada hora]
[13:45 - Detecta evento a las 14:00]
IA: "⏰ Recordatorio: En 15 minutos

     Reunión con cliente (14:00-15:00)
     📍 Oficina Central, Sala 3

     ✅ Checklist:
     - Laptop con presentación
     - Documentos impresos
     - Tarjetas de presentación"
```

#### weekly-report - Reporte semanal

```bash
clical weekly-report --user=USER_ID
```

**Argumentos:**
- `--user` (requerido) - ID del usuario

**Ejemplo:**
```bash
clical weekly-report --user=123456789
```

**Contenido:**
- Resumen de la semana (lunes a domingo)
- Eventos agrupados por día
- Estadísticas semanales

**Cuándo usar:**
- **Lunes 07:00 AM** - Inicio de semana
- Cuando el usuario pregunta "qué tengo esta semana"
- Para planificación semanal

**Tips para IA:**
- Ejecutar automáticamente cada lunes
- Identificar días pesados
- Sugerir reorganización si es necesario

---

## Patrones de Uso para IA

### Patrón 1: Saludo Matutino (07:00 AM)

```bash
# Ejecutar automáticamente
REPORT=$(clical daily-report --user=123456789)

# Procesar y presentar al usuario conversacionalmente
```

**Script de ejemplo:**

```bash
#!/bin/bash
USER_ID="123456789"

# Obtener reporte
REPORT=$(clical daily-report --user=$USER_ID)

# Enviar al usuario (ejemplo con Telegram)
send-telegram-message "$REPORT"
```

### Patrón 2: Monitoreo Horario (cada hora 09:00-18:00)

```bash
# Ejecutar cada hora
UPCOMING=$(clical upcoming-report --user=123456789 --hours=2)

# Si hay eventos, alertar
if [ -n "$UPCOMING" ]; then
    send-alert "$UPCOMING"
fi
```

### Patrón 3: Resumen Nocturno (20:00 PM)

```bash
# Vista de mañana
TOMORROW=$(clical tomorrow-report --user=123456789)

send-telegram-message "$TOMORROW"
```

### Patrón 4: Agregar Evento desde Conversación

Cuando el usuario dice: *"Tengo reunión con el cliente mañana a las 2 de la tarde"*

**Proceso de IA:**

1. **Extraer información:**
   - Título: "Reunión con cliente"
   - Fecha: mañana → calcular fecha
   - Hora: 2 de la tarde → 14:00
   - Duración: inferir (reuniones ~60 min)
   - Tags: inferir (cliente, trabajo)

2. **Verificar conflictos:**
```bash
# Ver qué tiene ese día
clical list --user=123456789 --date="2025-11-22"
```

3. **Agregar evento:**
```bash
clical add --user=123456789 \
  --datetime="2025-11-22 14:00" \
  --title="Reunión con cliente" \
  --duration=60 \
  --tags=cliente,trabajo
```

4. **Confirmar:**
   - "✓ Agendé tu reunión con cliente para mañana 22 de noviembre a las 14:00 (1 hora)"

### Patrón 5: Responder "Qué tengo hoy/mañana/esta semana"

**Usuario:** "Qué tengo hoy?"

```bash
# Opción 1: Reporte completo
clical daily-report --user=123456789

# Opción 2: Lista simple
clical list --user=123456789 --range=today
```

**Presentar al usuario:** Procesar y formatear conversacionalmente

**Usuario:** "Qué tengo mañana?"

```bash
clical tomorrow-report --user=123456789
```

**Usuario:** "Qué tengo esta semana?"

```bash
clical weekly-report --user=123456789
# o
clical list --user=123456789 --range=week
```

### Patrón 6: Buscar Eventos Específicos

**Usuario:** "Cuándo es mi próxima reunión con el cliente?"

```bash
# Buscar eventos con "cliente"
clical list --user=123456789 --tags=cliente

# O buscar en upcoming
clical upcoming-report --user=123456789 --count=20 | grep -i cliente
```

### Patrón 7: Sugerir Bloques Libres

**Usuario:** "Cuándo tengo tiempo libre para trabajar en el proyecto?"

```bash
# Obtener daily report (incluye bloques libres)
clical daily-report --user=123456789

# Identificar bloques de 2+ horas
# Sugerir al usuario
```

---

## Configuración para Automatización

### Crontab para IA Assistant

```bash
# Editar crontab
crontab -e

# Agregar:

# Reporte diario (7:00 AM)
0 7 * * * /usr/local/bin/clical daily-report --user=123456789 | tu-script-ia

# Monitoreo horario (9 AM - 6 PM, cada hora)
0 9-18 * * * /usr/local/bin/clical upcoming-report --user=123456789 --hours=2 | tu-script-ia

# Reporte nocturno (8:00 PM)
0 20 * * * /usr/local/bin/clical tomorrow-report --user=123456789 | tu-script-ia

# Reporte semanal (Lunes 7:00 AM)
0 7 * * 1 /usr/local/bin/clical weekly-report --user=123456789 | tu-script-ia
```

---

## Tips Avanzados para IA

### 1. Inferir Duración de Eventos

Si el usuario no especifica duración:

- Stand-ups: 15 min
- Llamadas cortas: 30 min
- Reuniones: 60 min
- Talleres/workshops: 120-180 min
- Entrevistas: 60-90 min

### 2. Extraer Tags Automáticamente

De las palabras del usuario:
- "reunión" → tag: reunion
- "cliente" → tag: cliente
- "desarrollo" → tag: desarrollo
- "llamada" → tag: llamada
- "importante/urgente" → tag: importante

### 3. Detectar Preparación Necesaria

Si el evento menciona:
- "presentación" → Recordar laptop y preparar slides
- "documentos" → Recordar imprimir/llevar
- "sala/ubicación física" → Recordar llegar 5 min antes
- "virtual/zoom" → Recordar link y probar audio/video

### 4. Sugerir Bloques de Tiempo

Al presentar bloques libres:
- < 30 min: "Ideal para: emails, llamadas cortas"
- 30-60 min: "Ideal para: reuniones, tareas pequeñas"
- 60-120 min: "Ideal para: trabajo enfocado, desarrollo"
- 120+ min: "Ideal para: trabajo profundo, proyectos grandes"

### 5. Alertas Inteligentes

- 15 min antes: Eventos importantes
- 30 min antes: Eventos que requieren desplazamiento
- 1 día antes: Eventos que requieren preparación
- Al inicio del día: Si hay >5 eventos (día pesado)

---

## Formatos de Fecha/Hora

### Entrada (al agregar/editar eventos)

Formato estricto: `"YYYY-MM-DD HH:MM"`

Ejemplos válidos:
- `"2025-11-21 09:00"`
- `"2025-11-21 14:30"`
- `"2025-12-01 08:00"`

**Conversión desde lenguaje natural (tarea de la IA):**

- "mañana a las 2 pm" → calcular fecha + "14:00"
- "el viernes a las 10" → calcular fecha del próximo viernes + "10:00"
- "en 2 horas" → calcular datetime actual + 2 horas

### Salida (en reportes)

Los reportes usan formato legible:
- Fechas: "Viernes 21 de Noviembre, 2025"
- Horas: "14:00", "09:00"
- Rangos: "14:00 - 15:00"

---

## Errores Comunes y Soluciones

### Error: "user_id es requerido"

**Causa:** No se especificó `--user`

**Solución:** Siempre incluir `--user=USER_ID` en comandos de eventos

### Error: "datetime es requerido"

**Causa:** Falta `--datetime` al agregar evento

**Solución:** Incluir `--datetime="YYYY-MM-DD HH:MM"`

### Error: "entrada no encontrada"

**Causa:** ID de evento inválido

**Solución:** Usar `list` primero para obtener IDs correctos

### Error: "formato inválido"

**Causa:** Formato de fecha/hora incorrecto

**Solución:** Usar formato exacto `"YYYY-MM-DD HH:MM"` con comillas

---

## Resumen de Comandos por Frecuencia de Uso

### Uso Diario (IA debe ejecutar frecuentemente)

1. `daily-report` - Cada mañana
2. `upcoming-report` - Cada hora
3. `list --range=today` - Verificar agenda
4. `add` - Agregar eventos según conversación

### Uso Semanal

1. `weekly-report` - Cada lunes
2. `list --range=week` - Planificación semanal

### Uso Ocasional

1. `edit` - Cuando cambian planes
2. `delete` - Cuando se cancelan eventos
3. `show` - Para ver detalles específicos
4. `user add/list/show` - Gestión de usuarios

---

## Checklist para Implementar IA Assistant

- [ ] Configurar cron para reportes automáticos
- [ ] Implementar parser de lenguaje natural → datetime
- [ ] Implementar sistema de notificaciones (Telegram/Email)
- [ ] Crear lógica de inferencia (duración, tags)
- [ ] Implementar detección de conflictos
- [ ] Crear templates de mensajes conversacionales
- [ ] Implementar almacenamiento de preferencias de usuario
- [ ] Agregar logging de interacciones
- [ ] Configurar alertas inteligentes
- [ ] Testear flujos completos

---

## 9. Sistema de Alarmas

### 9.1 Conceptos

Las alarmas son recordatorios programados que permiten a la IA realizar seguimientos y tareas programadas. A diferencia de los eventos del calendario, las alarmas se disparan automáticamente vía cron y pueden ser:

- **One-time**: Se ejecutan una sola vez en una fecha/hora específica
- **Recurrentes**: Se ejecutan repetidamente (daily, weekly, monthly, yearly)

### 9.2 Comandos

#### `alarm add` - Agregar Alarma

**Alarmas One-Time:**
```bash
# Formato básico
clical alarm add --user <user_id> --at "YYYY-MM-DD HH:MM" --context "..."

# Ejemplos
clical alarm add --user ai-agent --at "2025-11-24 10:00" --context "Revisar PR de autenticación"
clical alarm add --user ai-agent --at "2025-11-24 14:30" --context "Seguimiento del deploy en producción"
```

**Alarmas Recurrentes - Daily:**
```bash
# Cada día a la hora especificada
clical alarm add --user ai-agent --daily "14:30" --context "Revisar métricas diarias"

# Con fecha de expiración
clical alarm add --user ai-agent --daily "09:00" --expires "2025-12-31" --context "Stand-up temporal"
```

**Alarmas Recurrentes - Weekly:**
```bash
# Cada semana en el día especificado
clical alarm add --user ai-agent --weekly "monday 14:30" --context "Reunión semanal"
clical alarm add --user ai-agent --weekly "friday 17:00" --context "Reporte semanal"
```

**Alarmas Recurrentes - Monthly:**
```bash
# Cada mes en el día especificado (1-31)
clical alarm add --user ai-agent --monthly "1 09:00" --context "Reporte mensual"
clical alarm add --user ai-agent --monthly "15 14:30" --context "Revisión quincenal"
```

**Alarmas Recurrentes - Yearly:**
```bash
# Cada año en la fecha especificada
clical alarm add --user ai-agent --yearly "01-01 00:00" --context "Feliz año nuevo"
clical alarm add --user ai-agent --yearly "11-21 10:00" --context "Aniversario del proyecto"
```

#### `alarm check` - Verificar Alarmas

```bash
# Ejecutar verificación (para cron)
clical alarm check --user ai-agent

# Con verbose para debugging
clical alarm check --user ai-agent --verbose
```

**Comportamiento:**
- Si NO hay alarmas: no produce output (exit 0)
- Si hay alarmas: emite JSON a stdout con las alarmas

**Output JSON:**
```json
[
  {
    "id": "alarm_once_1234567890_abcd1234",
    "scheduled_for": "2025-11-24T10:00:00Z",
    "context": "Revisar PR de autenticación",
    "created_at": "2025-11-23T14:00:00Z",
    "recurrence": "once"
  },
  {
    "id": "alarm_weekly_1234567890_efgh5678",
    "scheduled_for": "2025-11-24T14:30:00Z",
    "context": "Reunión semanal",
    "created_at": "2025-11-20T10:00:00Z",
    "recurrence": "weekly",
    "expires_at": "2025-12-31T23:59:59Z"
  }
]
```

#### `alarm list` - Listar Alarmas

```bash
# Listar alarmas activas (formato tabla)
clical alarm list --user ai-agent

# Incluir alarmas pasadas
clical alarm list --user ai-agent --past

# Output JSON (para scripting)
clical alarm list --user ai-agent --json
```

**Output tabla:**
```
ALARMAS ACTIVAS:

ID                        TIPO       PROGRAMADA           CONTEXTO
----------------------------------------------------------------------------------------------------
alarm_once_1234567890_... once       2025-11-24 10:00     Revisar PR de autenticación
alarm_weekly_1234567890_..weekly     monday 14:30:00      Reunión semanal
alarm_daily_1234567890_...daily      14:30:00             Revisar métricas
```

#### `alarm cancel` - Cancelar Alarma

```bash
# Cancelar por ID
clical alarm cancel --user ai-agent alarm_once_1234567890_abcd1234
```

### 9.3 Integración con Cron

**Configurar cron para ejecutar cada minuto:**

```bash
# Editar crontab
crontab -e

# Agregar línea
* * * * * /usr/local/bin/clical-alarm-processor
```

**Script procesador (`/usr/local/bin/clical-alarm-processor`):**

```bash
#!/bin/bash

# Usuario para las alarmas
USER="ai-agent"

# Ejecutar check
OUTPUT=$(clical alarm check --user "$USER" 2>/dev/null)

# Si hay alarmas, procesarlas
if [ -n "$OUTPUT" ]; then
  # Enviar a Telegram
  echo "$OUTPUT" | jq -r '.[] | .context' | while read -r ctx; do
    s gobot-send-message "[ALARMA] $ctx"
  done

  # Notificar a agente IA
  echo "$OUTPUT" | ai-agent-notify --stdin
fi
```

### 9.4 Casos de Uso para IA

#### Caso 1: Seguimiento de Tareas

```
Usuario: "Recordame en 30 minutos revisar si el deploy terminó"

IA interpreta:
- Momento: now + 30 minutos
- Tipo: one-time
- Contexto: "Revisar estado del deploy a producción"

Comando:
clical alarm add --user ai-agent --at "2025-11-24 14:45" --context "Revisar estado del deploy a producción"
```

#### Caso 2: Reportes Automáticos

```
Usuario: "Quiero que cada viernes me des un resumen semanal"

IA configura:
clical alarm add --user ai-agent --weekly "friday 17:00" --context "Generar resumen semanal de actividades"
```

#### Caso 3: Alarmas con Expiración

```
Usuario: "Durante diciembre, recordame todos los días a las 9 AM revisar métricas de ventas"

IA configura:
clical alarm add --user ai-agent --daily "09:00" --expires "2025-12-31" --context "Revisar métricas de ventas navideñas"
```

### 9.5 Recovery Automático

El comando `alarm check` incluye recovery automático:
- Revisa últimos 60 minutos de alarmas one-time perdidas
- Si el sistema estuvo apagado, ejecuta todas las pendientes
- Alarmas ejecutadas se mueven a `past/`

### 9.6 Almacenamiento

**Estructura de directorios:**
```
data/users/<user_id>/alarms/
├── pending/                 # One-time alarms
│   └── 2025-11-24_10-00-00.json
├── recurring/
│   ├── daily/
│   │   └── 14-30-00.json
│   ├── weekly/
│   │   └── monday_14-30-00.json
│   ├── monthly/
│   │   └── 15_14-30-00.json
│   └── yearly/
│       └── 11-21_10-00-00.json
└── past/
    ├── one-time/
    └── recurring/
```

**Formato de archivo (JSON array):**
```json
[
  {
    "id": "alarm_once_1234567890_abcd1234",
    "context": "Contexto de la alarma",
    "created_at": "2025-11-23T14:00:00Z",
    "recurrence": "once"
  }
]
```

### 9.7 Consejos para IA

1. **Parsear el contexto:** El campo `context` debe contener toda la información necesaria para que la IA sepa qué hacer
2. **IDs visibles:** El ID de la alarma se muestra en la salida de `alarm-check`, úsalo para tracking
3. **Combinar con calendario:** Las alarmas pueden referenciar eventos del calendario
4. **Batch processing:** `alarm check` puede retornar múltiples alarmas en un solo JSON

**Ejemplo de procesamiento:**
```bash
# La IA ejecuta esto cada minuto
clical alarm check --user ai-agent | jq -r '.[] | "\(.id)|\(.context)"' | while IFS='|' read -r id context; do
  echo "Procesando alarma $id: $context"
  # La IA procesa el contexto y toma acción
done
```

---

**Versión:** 1.1
**Última actualización:** 2025-11-23
**Proyecto:** clical v0.2.0
