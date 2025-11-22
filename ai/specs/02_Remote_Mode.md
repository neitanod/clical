# Modo Remoto - Desarrollo con Telegram

**Fecha:** 2025-11-21
**Tipo:** Workflow de Desarrollo
**Estado:** Activo

## Descripción

El "modo remoto" es un workflow de desarrollo que permite a Claude Code recibir instrucciones vía Telegram mientras el chat de consola está abierto. Esto permite al usuario dar instrucciones desde cualquier lugar sin necesidad de estar frente a la computadora.

## Funcionamiento

### Script de Comunicación

El sistema utiliza un script llamado `s` (ubicado en el directorio del proyecto) que proporciona dos comandos principales:

1. **`s wait-for-telegram-response`**
   - Espera bloqueantemente un mensaje de Telegram
   - Cuando llega un mensaje, lo imprime en stdout y termina
   - Claude Code procesa ese mensaje como si viniera del chat de consola

2. **`s gobot-send-message <mensaje>`**
   - Envía un mensaje desde Claude Code hacia Telegram
   - Permite notificar al usuario sobre el progreso

### Flujo de Trabajo

```
┌─────────────┐
│ Claude Code │
└──────┬──────┘
       │
       │ 1. Ejecuta: s wait-for-telegram-response (BLOQUEANTE)
       │
       ▼
┌─────────────────┐
│ Script espera   │◄────── 2. Usuario envía mensaje desde Telegram
│ mensaje Telegram│
└──────┬──────────┘
       │
       │ 3. Script imprime mensaje y termina
       │
       ▼
┌─────────────┐
│ Claude Code │◄────── 4. Procesa mensaje como si viniera del chat
└──────┬──────┘
       │
       │ 5. Ejecuta tareas solicitadas
       │
       ▼
┌─────────────┐
│ s send-to-  │────── 6. Notifica progreso al usuario vía Telegram
│ telegram    │
└─────────────┘
       │
       │ 7. Vuelve a ejecutar: s wait-for-telegram-response
       │
       └──────► (repite el ciclo)
```

## Reglas Importantes

### 1. Ejecutar con Timeout para Bloqueo Correcto

**❌ INCORRECTO:**
```bash
# Sin timeout (se ejecuta en background automáticamente)
s wait-for-telegram-response

# O ejecutando manualmente en background
s wait-for-telegram-response &
# o usando run_in_background: true
```

**✅ CORRECTO:**
```bash
# Con timeout de 10 minutos (600000 ms)
s wait-for-telegram-response
# Usando parámetro timeout: 600000
```

**Implementación en Claude Code:**
```xml
<invoke name="Bash">
<parameter name="command">s wait-for-telegram-response</parameter>
<parameter name="timeout">600000</parameter>
</invoke>
```

**Razón:** Sin el parámetro `timeout`, Claude Code ejecuta automáticamente el comando en background. Cuando un comando corre en background, Claude Code no espera su resultado antes de continuar, por lo que **nunca recibirá la respuesta de Telegram** y el modo remoto no funcionará. El timeout de 600000ms (10 minutos) hace que el comando sea bloqueante, forzando a Claude Code a esperar la respuesta antes de poder hacer cualquier otra cosa.

### 2. Siempre Notificar Antes de Escuchar

**Regla:** Cuando estás en modo remoto, **SIEMPRE** debes enviar un mensaje a Telegram ANTES de volver a ejecutar `wait-for-telegram-response`.

**Razón:** El usuario necesita saber que completaste una tarea o que estás listo para la siguiente instrucción. Si simplemente vuelves a escuchar sin notificar, el usuario no sabrá que estás esperando.

**Ejemplo correcto:**

```bash
# Claude Code termina una tarea
# 1. PRIMERO notifica
s gobot-send-message "Tests de storage completados. 95% cobertura. Todos pasan."

# 2. LUEGO vuelve a escuchar
s wait-for-telegram-response  # Con timeout: 600000
```

**Ejemplos de mensajes apropiados:**
- "Tarea completada. Listo para la siguiente instrucción."
- "Compilación exitosa. ¿Qué sigue?"
- "Tests creados. 15 tests, todos pasan. Esperando instrucciones."
- "Error encontrado: [descripción]. ¿Cómo procedo?"

### 3. Mantener Modo Remoto Cuando Hay Dudas

**Regla:** Si Claude Code tiene dudas sobre qué hacer o necesita instrucciones adicionales del usuario, **DEBE mantener el modo remoto activo**.

**Razón:** Si el usuario no está presente en la consola, no podrá responder preguntas. Al mantener modo remoto activo, el usuario puede responder vía Telegram cuando esté disponible.

**Ejemplo:**

```bash
# Claude Code termina una tarea pero no está seguro del siguiente paso
# ❌ INCORRECTO: Preguntar en el chat y esperar
# ✅ CORRECTO: Enviar pregunta por Telegram y volver a modo remoto

s gobot-send-message "Tarea completada. ¿Qué debo hacer ahora?"
s wait-for-telegram-response  # Con timeout: 600000
```

## Casos de Uso

### Caso 1: Desarrollo Remoto Completo

Usuario está lejos de la computadora y quiere que Claude Code trabaje:

1. Usuario: Inicia sesión y pone a Claude Code en modo remoto
2. Claude Code: Ejecuta `s wait-for-telegram-response`
3. Usuario: Envía instrucciones vía Telegram (ej: "implementa autenticación")
4. Claude Code: Recibe mensaje, implementa, reporta progreso vía `s gobot-send-message`
5. Claude Code: Vuelve a ejecutar `s wait-for-telegram-response`
6. Usuario: Continúa enviando instrucciones según necesite

### Caso 2: Desarrollo Híbrido

Usuario alterna entre estar presente y ausente:

- Cuando está presente: Usa el chat de consola normalmente
- Cuando se va: Le dice a Claude Code "entra en modo remoto"
- Claude Code: Ejecuta `s wait-for-telegram-response` y queda esperando
- Usuario: Puede enviar instrucciones desde donde esté

### Caso 3: Tareas Largas con Supervisión Remota

Claude Code está ejecutando una tarea larga:

1. Claude Code: Reporta inicio de tarea vía `s gobot-send-message "Iniciando compilación..."`
2. Claude Code: Ejecuta tarea
3. Claude Code: Reporta completitud vía `s gobot-send-message "Compilación exitosa"`
4. Claude Code: Vuelve a modo remoto con `s wait-for-telegram-response`

## Ejemplo Práctico

### Sesión Completa en Modo Remoto

```bash
# Usuario dice: "Entra en modo remoto"
$ s wait-for-telegram-response

# [Bloqueado esperando Telegram...]
# Usuario envía por Telegram: "Crea tests para el módulo de storage"

# Script imprime y termina:
"Crea tests para el módulo de storage"

# Claude Code procesa el mensaje:
# - Lee pkg/storage/filesystem.go
# - Crea pkg/storage/filesystem_test.go
# - Ejecuta: go test ./pkg/storage

# Claude Code reporta:
$ s gobot-send-message "Tests de storage creados. 15 tests, 95% cobertura. Todos pasan."

# Claude Code vuelve a modo remoto:
$ s wait-for-telegram-response

# [Bloqueado esperando siguiente instrucción...]
```

## Limitaciones

1. **No Bidireccional Simultáneo:** No se puede estar en chat de consola Y Telegram al mismo tiempo
2. **Espera Bloqueante:** Mientras espera Telegram, el chat de consola no responde
3. **Requiere Script:** Necesita que el script `s` esté funcionando correctamente

## Beneficios

1. **Desarrollo Asíncrono:** Usuario puede dar instrucciones desde cualquier lugar
2. **Continuidad:** Claude Code puede trabajar en tareas largas sin supervisión presencial
3. **Notificaciones:** Usuario recibe actualizaciones de progreso vía Telegram
4. **Flexibilidad:** Usuario puede alternar entre consola y Telegram según necesidad

## Mejores Prácticas

### Para Claude Code

1. **Siempre ejecutar bloqueantemente:** Usar timeout 600000 para `wait-for-telegram-response`
2. **Reportar progreso:** Usar `s gobot-send-message` para mantener informado al usuario
3. **Ante dudas, modo remoto:** Si no sabes qué hacer, quédate esperando instrucciones
4. **Mensajes claros:** Al reportar vía Telegram, ser conciso y claro

### Para Usuario

1. **Instrucciones claras:** Ser específico en lo que se pide vía Telegram
2. **Confirmar recepción:** Si Claude Code reporta algo importante, responder
3. **Salir de modo remoto:** Avisar explícitamente cuando quieras usar la consola

## Resumen de Comandos

```bash
# Entrar en modo remoto (espera bloqueante con timeout)
s wait-for-telegram-response
# IMPORTANTE: Usar con timeout: 600000

# Enviar mensaje a Telegram
s gobot-send-message "mensaje aquí"
```

---

**Documentado:** 2025-11-21
**Uso:** Activo en desarrollo de clical
**Estado:** Operativo
