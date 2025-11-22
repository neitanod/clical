# Documentación de Modo Remoto

**Fecha:** 2025-11-21
**Tipo:** Documentation
**Estado:** Completado

## Objetivo

Documentar el workflow de "modo remoto" que permite a Claude Code recibir instrucciones vía Telegram mientras el chat de consola está abierto.

## Contexto

Durante el desarrollo de clical, se implementó un sistema que permite al usuario dar instrucciones a Claude Code desde cualquier lugar usando Telegram. Este modo es fundamental para el desarrollo asíncrono del proyecto.

## Cambios Implementados

### 1. Archivo de Especificación Creado

**Archivo:** `ai/specs/02_Remote_Mode.md`

**Contenido documentado:**

1. **Descripción del sistema**
   - Explicación del concepto de modo remoto
   - Flujo de trabajo con diagrama
   - Scripts de comunicación (`s wait-for-telegram-response` y `s gobot-send-message`)

2. **Reglas importantes**

   **Regla #1: Ejecutar con Timeout**
   - DEBE usar `timeout: 600000` (10 minutos)
   - Sin timeout, el comando se ejecuta en background automáticamente
   - En background, Claude Code nunca recibe la respuesta de Telegram
   - El modo remoto NO funciona sin timeout

   ```bash
   # Correcto
   s wait-for-telegram-response
   # Con parámetro timeout: 600000
   ```

   **Regla #2: Siempre Notificar Antes de Escuchar**
   - SIEMPRE enviar mensaje vía `s gobot-send-message` antes de volver a `wait-for-telegram-response`
   - El usuario necesita saber que una tarea terminó o que Claude está listo
   - Sin notificación, el usuario no sabe que Claude está esperando

   ```bash
   # 1. PRIMERO notificar
   s gobot-send-message "Tarea completada. Listo para siguiente instrucción."

   # 2. LUEGO escuchar
   s wait-for-telegram-response  # Con timeout: 600000
   ```

   **Regla #3: Mantener Modo Remoto Cuando Hay Dudas**
   - Si Claude Code no está seguro del siguiente paso, mantener modo remoto
   - Enviar pregunta por Telegram y esperar respuesta
   - No preguntar en chat de consola si el usuario no está presente

3. **Casos de uso**
   - Desarrollo remoto completo
   - Desarrollo híbrido (presencial + remoto)
   - Tareas largas con supervisión remota

4. **Ejemplos prácticos**
   - Sesión completa en modo remoto
   - Flujo de trabajo con notificaciones

5. **Limitaciones y beneficios**
   - Documentadas las restricciones del sistema
   - Explicados los beneficios para desarrollo asíncrono

6. **Mejores prácticas**
   - Para Claude Code
   - Para el usuario

## Comandos del Sistema

### s wait-for-telegram-response

**Propósito:** Esperar bloqueantemente un mensaje de Telegram

**Uso correcto:**
```xml
<invoke name="Bash">
<parameter name="command">s wait-for-telegram-response</parameter>
<parameter name="timeout">600000</parameter>
</invoke>
```

**Comportamiento:**
- Bloquea la sesión de Claude Code
- Espera hasta recibir mensaje de Telegram
- Imprime el mensaje en stdout y termina
- Claude procesa el mensaje como si viniera del chat

### s gobot-send-message

**Propósito:** Enviar mensaje a Telegram desde Claude Code

**Uso:**
```bash
s gobot-send-message "mensaje aquí"
```

**Comportamiento:**
- Envía mensaje inmediatamente
- Retorna confirmación de envío
- Permite mantener al usuario informado del progreso

## Flujo de Trabajo Típico

```
1. Usuario dice: "Entra en modo remoto"
2. Claude ejecuta: s wait-for-telegram-response (timeout: 600000)
3. [Bloqueado esperando...]
4. Usuario envía por Telegram: "Crea tests de storage"
5. Claude recibe mensaje
6. Claude implementa los tests
7. Claude notifica: s gobot-send-message "Tests creados. 95% cobertura."
8. Claude vuelve a paso 2
```

## Aprendizajes Clave

1. **Timeout es crítico:** Sin timeout, el modo remoto no funciona porque el comando se ejecuta en background y Claude no espera la respuesta.

2. **Notificación es esencial:** El usuario debe saber cuándo Claude está esperando. Notificar antes de cada `wait-for-telegram-response` es obligatorio.

3. **Modo remoto para dudas:** Si Claude no sabe qué hacer, debe preguntar por Telegram y esperar, no preguntar en consola donde el usuario puede no estar.

## Beneficios Implementados

1. **Desarrollo asíncrono:** Usuario puede dar instrucciones desde cualquier lugar
2. **Continuidad:** Claude puede trabajar en tareas largas sin supervisión presencial
3. **Notificaciones:** Usuario recibe actualizaciones de progreso
4. **Flexibilidad:** Alternancia entre consola y Telegram según necesidad

## Archivos Modificados

- `ai/specs/02_Remote_Mode.md` - Creado (230+ líneas)

## Testing

- ✅ Comando `s wait-for-telegram-response` con timeout funciona correctamente
- ✅ Comando `s gobot-send-message` envía mensajes exitosamente
- ✅ Flujo completo probado: notificar → esperar → recibir → procesar
- ✅ Documentación validada con ejemplos reales

## Próximos Pasos

- Ninguno - Documentación completa y operativa
- Sistema de modo remoto listo para uso continuo en desarrollo de clical

---

**Estado:** ✅ Completado
**Documentación:** ai/specs/02_Remote_Mode.md
**Fecha:** 2025-11-21
