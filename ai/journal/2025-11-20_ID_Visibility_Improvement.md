# Mejora de Visibilidad de IDs en Listados

**Fecha:** 2025-11-20
**Tipo:** Feature Enhancement
**Estado:** Completado

## Objetivo

Mejorar la usabilidad para agentes IA mostrando el ID de los eventos directamente en todos los listados, eliminando la necesidad de ejecutar comandos adicionales para obtener el ID antes de operar sobre un evento.

## Motivación

**Problema anterior:**
- Los agentes IA necesitaban ejecutar 2 comandos para operar sobre un evento:
  1. `list` para encontrar el evento
  2. Buscar manualmente el ID en la salida
  3. Luego ejecutar `edit`/`delete` con ese ID

**Solución:**
- Incluir el ID directamente en todas las salidas de listados
- Permitir que la IA opere en un solo paso

## Cambios Implementados

### 1. `pkg/reporter/daily.go`

**Sección "Próximo Evento":**
- Agregada línea: `- ID: {id}`
- Posición: Después del título, antes de duración

**Sección "Agenda de Hoy":**
- Agregada línea: `- ID: {id}` en cada evento
- Posición: Consistente con "Próximo Evento"

**Formato:**
```markdown
**[09:00 - 09:15] Stand-up Meeting**
- ID: 161cbab48ae81b66
- Duración: 15 min
- Tags: #trabajo
```

### 2. `internal/cli/reports.go` - upcoming-report

**Cambio:**
- Agregado emoji: `🆔 {id}`
- Posición: Después del título y duración

**Formato:**
```
⏰ **En 15 minutos** (09:00)
   Stand-up Meeting (15 min)
   🆔 161cbab48ae81b66
   📍 Sala 2
```

### 3. `internal/cli/reports.go` - weekly-report

**Cambio:**
- Agregado `[ID: {id}]` en línea compacta

**Formato:**
```
- [09:00] Stand-up Meeting (15 min) [ID: 161cbab48ae81b66]
```

### 4. `internal/cli/list.go`

**Cambio:**
- Movido ID a la misma línea con formato `[ID: {id}]`
- Removida línea separada de ID

**Antes:**
```
[2025-11-21 09:00] Stand-up Meeting (15 min) #trabajo
  ID: 161cbab48ae81b66
```

**Después:**
```
[2025-11-21 09:00] Stand-up Meeting (15 min) [ID: 161cbab48ae81b66] #trabajo
```

## Archivos Modificados

1. `pkg/reporter/daily.go` - 2 ubicaciones
2. `internal/cli/reports.go` - 2 comandos (upcoming, weekly)
3. `internal/cli/list.go` - 1 función

**Total:** 3 archivos, 5 cambios

## Testing

### Comandos Probados

```bash
# list - ID visible en línea
clical list --user=12345 --range=today
✓ Salida: [2025-11-21 09:00] Daily Stand-up (15 min) [ID: 161cbab48ae81b66] #trabajo

# daily-report - ID en cada evento de agenda
clical daily-report --user=12345 --date="2025-11-21"
✓ Salida: Cada evento muestra "- ID: {id}"

# upcoming-report - ID con emoji
clical upcoming-report --user=12345 --hours=24
✓ Salida: Cada evento muestra "🆔 {id}"

# weekly-report - ID compacto
clical weekly-report --user=12345
✓ Salida: Cada evento muestra "[ID: {id}]"
```

### Verificación Manual

- ✅ IDs son copiables fácilmente
- ✅ Formato consistente entre comandos
- ✅ No rompe parsing existente
- ✅ Mejora legibilidad para humanos también

## Beneficios

### Para Agentes IA

1. **Reducción de comandos**: 1 en vez de 2
2. **Menos errores**: No hay ambigüedad al identificar eventos
3. **Mejor eficiencia**: Parsing directo del output
4. **Operaciones más rápidas**: Copiar ID y ejecutar inmediatamente

### Para Usuarios Humanos

1. **Más información visible**: Todo en un comando
2. **Copiar/pegar ID fácil**: Para operaciones manuales
3. **Debugging más simple**: Ver IDs directamente

## Ejemplo de Uso Mejorado

**Antes:**
```bash
# Usuario (o IA) necesita 2 pasos:
$ clical list --user=12345 --range=today
[2025-11-21 09:00] Meeting (15 min) #trabajo
  ID: 161cbab48ae81b66

# Copiar ID manualmente...
$ clical edit --user=12345 --id=161cbab48ae81b66 --duration=30
```

**Después:**
```bash
# Ahora en 1 paso (IA puede extraer ID directamente):
$ clical list --user=12345 --range=today
[2025-11-21 09:00] Meeting (15 min) [ID: 161cbab48ae81b66] #trabajo

# IA parsea: id="161cbab48ae81b66" y ejecuta:
$ clical edit --user=12345 --id=161cbab48ae81b66 --duration=30
```

## Compatibilidad

- ✅ Backwards compatible
- ✅ No rompe scripts existentes
- ✅ Salida de JSON no afectada (solo texto)
- ✅ Parsing más fácil (IDs siempre en misma posición)

## Próximos Pasos

- [ ] Actualizar USAGE.md con ejemplos de parsing de IDs
- [ ] Agregar ejemplos de uso de IA con IDs visibles
- [ ] Documentar formato de salida en specs

## Notas

- Los IDs son hexadecimales de 16 caracteres
- Formato consistente en todos los comandos de reporte
- Emoji 🆔 usado solo en upcoming-report para diferenciación visual

---

**Compilación:** Exitosa
**Testing:** Completo
**Instalación:** /usr/local/bin/clical actualizado
**Estado:** ✅ Implementado y probado
