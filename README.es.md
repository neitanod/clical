[English](README.md) · **Español**

# clical

**clical** es un sistema de calendario multiusuario con interfaz de línea de comandos (CLI), diseñado para ser asistido por Inteligencia Artificial.

## Características Principales

- 🗓️ **Calendario CLI multiusuario** - Gestión completa de eventos por línea de comandos
- 📝 **Almacenamiento en Markdown + JSON** - Archivos legibles y editables manualmente
- 🤖 **Optimizado para asistencia de IA** - Reportes diseñados para que una IA ayude proactivamente
- 📁 **Organización jerárquica** - Eventos organizados por año/mes/día
- 🔍 **Búsqueda y filtrado** - Potentes capacidades de búsqueda
- 📊 **Reportes inteligentes** - Daily, weekly, upcoming reports para IA
- ⏰ **Compatible con cron** - Ejecución programada de reportes

## Instalación

### Desde código fuente (Linux / macOS)

```bash
# Clonar repositorio
git clone https://github.com/neitanod/clical.git
cd clical

# Compilar
make build

# Instalar en /usr/local/bin (opcional)
make install-system
```

### Desde código fuente (Windows)

En Windows usá los scripts PowerShell equivalentes (no requieren `make`):

```powershell
git clone https://github.com/neitanod/clical.git
cd clical

# Compilar (genera .\clical.exe)
.\build.ps1

# Instalar en %USERPROFILE%\go\bin (opcional)
.\install.ps1
```

Para que el comando `clical` esté disponible en cualquier terminal, asegurate
de que `%USERPROFILE%\go\bin` esté en tu `PATH` de usuario.

### Instalación asistida por un agente de IA

Si usás un agente con acceso a tu terminal (Claude Code, Cursor, etc.), podés
instalar clical pegándole el siguiente prompt:

<https://github.com/neitanod/clical/blob/main/install_prompt.md>

El prompt guía al agente para detectar el sistema operativo, verificar
requisitos, clonar, compilar e instalar tanto en Linux/macOS como en Windows.

### Requisitos

- Go 1.23 o superior

## Uso Rápido

### 1. Crear un usuario

```bash
clical user add --id=12345 --name="Tu Nombre" --timezone="America/Argentina/Buenos_Aires"
```

### 2. Agregar un evento

```bash
clical add --user=12345 \
  --datetime="2025-11-21 14:00" \
  --title="Reunión con cliente" \
  --duration=60 \
  --location="Oficina Central" \
  --notes="Revisar propuesta Q4"
```

### 3. Listar eventos

```bash
# Todos los eventos
clical list --user=12345

# Eventos de hoy
clical list --user=12345 --range=today

# Eventos de esta semana
clical list --user=12345 --range=week
```

### 4. Ver reporte diario

```bash
clical daily-report --user=12345
```

## Comandos Disponibles

### Gestión de Usuarios

```bash
# Crear usuario
clical user add --id=ID --name="Nombre" --timezone="Timezone"

# Listar usuarios
clical user list

# Ver detalles de usuario
clical user show --id=ID
```

### Gestión de Eventos

```bash
# Agregar evento
clical add --user=ID --datetime="YYYY-MM-DD HH:MM" --title="Título" [opciones]

# Listar eventos
clical list --user=ID [--from=FECHA] [--to=FECHA] [--range=RANGO] [--tags=TAG1,TAG2]

# Ver evento
clical show --user=ID --id=EVENT_ID

# Editar evento
clical edit --user=ID --id=EVENT_ID [--title="Nuevo"] [--datetime="YYYY-MM-DD HH:MM"]

# Eliminar evento
clical delete --user=ID --id=EVENT_ID [--force]
```

### Reportes para IA

```bash
# Reporte diario completo
clical daily-report --user=ID [--date=YYYY-MM-DD]

# Reporte de mañana
clical tomorrow-report --user=ID

# Próximos eventos
clical upcoming-report --user=ID --hours=2
clical upcoming-report --user=ID --count=5

# Reporte semanal
clical weekly-report --user=ID
```

### Otros

```bash
# Versión
clical version

# Ayuda
clical --help
clical COMANDO --help
```

## Uso con Cron

### Reportes automáticos

Editar crontab:

```bash
crontab -e
```

Agregar líneas:

```bash
# Reporte diario a las 7:00 AM
0 7 * * * /usr/local/bin/clical daily-report --user=12345 | mail -s "Agenda de Hoy" tu@email.com

# Reporte de mañana a las 8:00 PM
0 20 * * * /usr/local/bin/clical tomorrow-report --user=12345 | mail -s "Agenda de Mañana" tu@email.com

# Alertas cada hora durante horario laboral
0 9-18 * * * /usr/local/bin/clical upcoming-report --user=12345 --hours=2 | mail -s "Próximos Eventos" tu@email.com

# Reporte semanal los lunes a las 7:00 AM
0 7 * * 1 /usr/local/bin/clical weekly-report --user=12345 | mail -s "Agenda Semanal" tu@email.com
```

### Integración con Telegram (si tienes un bot)

```bash
# Reporte diario por Telegram
0 7 * * * OUTPUT=$(/usr/local/bin/clical daily-report --user=12345); [ -n "$OUTPUT" ] && tu-comando-telegram "$OUTPUT"
```

## Formato de Almacenamiento

Los datos se guardan en `~/.clical/data/` (configurable con `--data-dir`):

```
~/.clical/data/
└── users/
    └── 12345/
        ├── user.md              # Info del usuario (Markdown)
        ├── user.json            # Metadata del usuario
        ├── events/
        │   └── 2025/
        │       └── 11/
        │           └── 21/
        │               ├── 09-00-stand-up-meeting.md
        │               ├── 09-00-stand-up-meeting.json
        │               ├── 14-00-reunion-con-cliente.md
        │               └── 14-00-reunion-con-cliente.json
        └── .state/
            └── report-state.json
```

### Ejemplo de archivo Markdown

```markdown
# Reunión con cliente

**Fecha:** 2025-11-21
**Hora:** 14:00
**Duración:** 60 minutos
**Ubicación:** Oficina Central
**Tags:** #trabajo #cliente

## Notas

Revisar propuesta Q4 y discutir timeline.

---

*Creado: 2025-11-20 16:18*
*Actualizado: 2025-11-20 16:18*
*ID: e36e10014ea57372*
```

## Configuración

### Variables de Entorno

```bash
# Directorio de datos (default: ~/.clical/data)
export CLICAL_DATA_DIR="/ruta/personalizada/data"

# Usuario por defecto
export CLICAL_USER_ID="12345"
```

### Timezones Comunes

- `America/Argentina/Buenos_Aires`
- `America/Mexico_City`
- `America/New_York`
- `Europe/Madrid`
- `UTC`

Lista completa: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones

## Uso con IA

clical está diseñado para que una IA te asista proactivamente. Ejemplo de flujo:

1. **07:00 AM** - Cron ejecuta `daily-report`
   - IA recibe agenda del día
   - IA saluda y presenta eventos
   - IA identifica tareas pendientes
   - IA sugiere organización

2. **Durante el día** - Cron ejecuta `upcoming-report` cada hora
   - IA alerta eventos próximos
   - IA recuerda preparación necesaria

3. **20:00 PM** - Cron ejecuta `tomorrow-report`
   - IA presenta vista de mañana
   - IA sugiere preparación nocturna

## Desarrollo

### Estructura del Proyecto

```
clical/
├── cmd/clical/       # Entry point
├── pkg/              # Paquetes públicos
│   ├── calendar/     # Modelos Entry, Filter
│   ├── storage/      # Storage filesystem
│   ├── user/         # User management
│   └── reporter/     # Report generation
├── internal/         # Paquetes privados
│   ├── cli/          # Comandos Cobra
│   └── config/       # Configuración
├── ai/               # Documentación de desarrollo
│   ├── specs/        # Especificaciones
│   └── journal/      # Journal de desarrollo
└── docs/             # Documentación web
```

### Compilar

```bash
make build
```

### Tests

```bash
make test
```

### Formatear código

```bash
make fmt
```

## Contribuir

Ver [ai/specs/00_Overview.md](ai/specs/00_Overview.md) para especificaciones completas.

## Licencia

MIT

## Autor

Desarrollado por Sebastián Valencia con asistencia de Claude (Anthropic).

---

**Versión:** 0.1.0
**Estado:** MVP funcional - En desarrollo activo
