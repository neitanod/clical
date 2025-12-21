# Modo Remoto - Desarrollo con Telegram

**Fecha:** 2025-11-21
**Tipo:** Workflow de Desarrollo
**Estado:** Activo

## Descripción

El "modo remoto" es un workflow de desarrollo que permite a Claude Code recibir instrucciones vía Telegram mientras el chat de consola está abierto. Esto permite al usuario dar instrucciones desde cualquier lugar sin necesidad de estar frente a la computadora.

Los mensajes de Telegram son **inyectados directamente** al chat local de Claude Code, por lo que Claude los procesa como si el usuario estuviera escribiendo en la consola.

## Funcionamiento

### Script de Comunicación

El sistema utiliza un script llamado `s` (ubicado en el directorio del proyecto) que proporciona los siguientes comandos:

#### Envío de Mensajes

**`s gobot-send-message <mensaje>`**
- Envía un mensaje de texto desde Claude Code hacia Telegram
- Permite notificar al usuario sobre el progreso y estado de las tareas

#### Envío de Archivos

**`s gobot-send-file <path-to-file>`**
- Envía un archivo cualquiera al usuario vía Telegram
- Útil para compartir logs, archivos de configuración, binarios, etc.
- El archivo se envía como documento

**`s gobot-send-image <path-to-image>`**
- Envía una imagen al usuario vía Telegram
- Formatos compatibles: JPG, PNG, GIF, WebP
- La imagen se muestra directamente en el chat (no como documento)
- Útil para compartir gráficos, diagramas, capturas de pantalla

**`s gobot-send-video <path-to-video>`**
- Envía un video al usuario vía Telegram
- Formatos compatibles: MP4, AVI, MOV, etc.
- El video se puede reproducir directamente en Telegram
- Útil para compartir demos, grabaciones de pantalla

**`s gobot-send-voice <path-to-audio>`**
- Envía un archivo de audio como mensaje de voz en Telegram
- Se muestra como mensaje de audio con forma de onda (diferente a un archivo adjunto)
- Formatos compatibles: MP3, OGG, WAV, etc.
- Útil para compartir grabaciones de audio, mensajes de voz generados

### Flujo de Trabajo

```
┌─────────────┐
│ Claude Code │ 1. Claude Code trabajando en una tarea
└──────┬──────┘
       │
       │ 2. Ejecuta tareas solicitadas
       │
       ▼
┌─────────────────────┐
│ s gobot-send-message│─────► 3. Notifica progreso al usuario vía Telegram
└─────────────────────┘


┌──────────────────────┐
│ Usuario envía mensaje│◄────── 4. Usuario responde desde Telegram
│ desde Telegram       │
└──────┬───────────────┘
       │
       │ 5. Mensaje inyectado directamente al chat local
       │
       ▼
┌─────────────┐
│ Claude Code │◄────── 6. Procesa mensaje como si viniera de la consola
└──────┬──────┘
       │
       └──────► (repite el ciclo)
```

## Reglas Importantes

### 1. Reportar Progreso Proactivamente

**Regla:** Cuando trabajas en tareas complejas o largas, reporta progreso regularmente vía `s gobot-send-message`.

**Razón:** El usuario puede no estar frente a la consola. Enviar actualizaciones vía Telegram le permite monitorear el progreso desde cualquier lugar.

**Formato de Reportes:**

Los reportes deben indicar claramente si el agente continuará trabajando o si ha terminado:

**Reportes Intermedios** (agente continúa trabajando):
```bash
s gobot-send-message "Tests creados. 15 tests, todos pasan. Continúo trabajando."
s gobot-send-message "Compilación exitosa. Binary: 2.0MB. Continúo trabajando."
s gobot-send-message "Módulo de autenticación implementado. Continúo trabajando."
```

**Reportes Finales** (agente ha terminado y espera instrucciones):
```bash
s gobot-send-message "Implementación completada. Tests pasando. Tarea terminada, esperando instrucciones."
s gobot-send-message "Refactorización finalizada. Código optimizado. Tarea terminada, esperando instrucciones."
s gobot-send-message "Error encontrado: [descripción]. Tarea terminada, esperando instrucciones."
```

**Regla de terminación:**
- Si vas a continuar con más subtareas → termina con "Continúo trabajando"
- Si has completado TODO y no ejecutarás más tareas → termina con "Tarea terminada, esperando instrucciones"

### 2. Usar Telegram para Comunicación Asíncrona

**Regla:** Cuando necesites información del usuario y no haya respuesta inmediata en consola, envía la pregunta vía Telegram.

**Razón:** El usuario puede estar ocupado o lejos de la consola. Telegram permite comunicación asíncrona - el usuario responderá cuando esté disponible, y su respuesta será inyectada automáticamente al chat local.

**Ejemplo:**

```bash
# Claude Code necesita clarificación
s gobot-send-message "Implementación completada. ¿Qué características debo agregar ahora?"
# Usuario responderá vía Telegram cuando esté disponible
# Respuesta será inyectada automáticamente al chat local
```

### 3. Compartir Archivos Cuando Sea Solicitado

**Regla:** Si el usuario solicita un archivo específico, enviarlo usando el comando apropiado según el tipo de archivo.

**Razón:** El usuario puede necesitar revisar logs, binarios, gráficos o cualquier archivo generado sin tener acceso directo a la máquina.

**Ejemplos:**

```bash
# Usuario pide: "Envíame el log de la compilación"
s gobot-send-file "/var/log/myapp.log"
s gobot-send-message "Log enviado. Tarea terminada, esperando instrucciones."

# Usuario pide: "Muéstrame un diagrama de la arquitectura"
s gobot-send-image "/tmp/architecture-diagram.png"
s gobot-send-message "Diagrama enviado. Tarea terminada, esperando instrucciones."

# Usuario pide: "Envíame el binario compilado"
s gobot-send-file "bin/myapp"
s gobot-send-message "Binary enviado (2.0MB). Tarea terminada, esperando instrucciones."

# Usuario pide: "Mándame ese video de la demo"
s gobot-send-video "/tmp/demo-recording.mp4"
s gobot-send-message "Video enviado. Tarea terminada, esperando instrucciones."
```

## Casos de Uso

### Caso 1: Desarrollo Remoto Completo

Usuario está lejos de la computadora y quiere que Claude Code trabaje:

1. Usuario: Deja la consola abierta con Claude Code activo
2. Usuario: Envía instrucción vía Telegram (ej: "implementa autenticación")
3. Sistema: Inyecta mensaje al chat local automáticamente
4. Claude Code: Recibe mensaje en consola, implementa la tarea
5. Claude Code: Reporta progreso vía `s gobot-send-message "Autenticación implementada. Tests pasan."`
6. Usuario: Continúa enviando instrucciones vía Telegram según necesite

### Caso 2: Desarrollo Híbrido

Usuario alterna entre estar presente y ausente:

- Cuando está presente: Usa el chat de consola normalmente
- Cuando se va: Envía instrucciones desde Telegram
- Claude Code: Procesa ambas fuentes de la misma manera (todo llega al chat local)
- Claude Code: Reporta progreso vía Telegram para mantener al usuario informado

### Caso 3: Tareas Largas con Supervisión Remota

Claude Code está ejecutando una tarea larga con múltiples subtareas:

1. Claude Code: Reporta inicio vía `s gobot-send-message "Iniciando refactorización del módulo storage. Continúo trabajando."`
2. Claude Code: Ejecuta primera subtarea (refactorizar filesystem.go)
3. Claude Code: Reporta vía `s gobot-send-message "filesystem.go refactorizado. Continúo trabajando."`
4. Claude Code: Ejecuta segunda subtarea (actualizar tests)
5. Claude Code: Reporta vía `s gobot-send-message "Tests actualizados. Continúo trabajando."`
6. Claude Code: Ejecuta tercera subtarea (ejecutar test suite completo)
7. Claude Code: Reporta completitud final vía `s gobot-send-message "Refactorización completa. 25 tests, todos pasan. Tarea terminada, esperando instrucciones."`
8. Usuario: Puede responder con siguiente instrucción cuando esté disponible

### Caso 4: Compartir Archivos Generados

Usuario solicita archivos específicos vía Telegram:

1. Usuario: Envía vía Telegram "Compila el proyecto y envíame el binario"
2. Sistema: Inyecta mensaje al chat local
3. Claude Code: Compila el proyecto
4. Claude Code: Reporta vía `s gobot-send-message "Compilación exitosa. Binary: 2.0MB. Continúo trabajando."`
5. Claude Code: Envía binario vía `s gobot-send-file "bin/myapp"`
6. Claude Code: Reporta vía `s gobot-send-message "Binary enviado. Tarea terminada, esperando instrucciones."`
7. Usuario: Recibe el archivo en Telegram y puede descargarlo

## Ejemplo Práctico

### Sesión Completa en Modo Remoto

```bash
# Usuario está lejos de la computadora
# Envía por Telegram: "Crea tests para el módulo de storage y luego optimiza el código"

# [Sistema inyecta mensaje al chat local de Claude Code]

# Claude Code ve en consola:
User: Crea tests para el módulo de storage y luego optimiza el código

# Claude Code empieza a trabajar:
# - Lee pkg/storage/filesystem.go
# - Crea pkg/storage/filesystem_test.go

# Claude Code reporta progreso intermedio:
$ s gobot-send-message "Tests de storage creados. 15 tests. Continúo trabajando."

# Claude Code continúa:
# - Ejecuta: go test ./pkg/storage
# - Analiza resultados

$ s gobot-send-message "Tests ejecutados. 95% cobertura, todos pasan. Continúo trabajando."

# Claude Code continúa con optimización:
# - Analiza filesystem.go
# - Aplica optimizaciones
# - Ejecuta tests nuevamente

$ s gobot-send-message "Código optimizado. Tests siguen pasando. Tarea terminada, esperando instrucciones."

# Usuario recibe notificación en Telegram
# Sabe que puede enviar siguiente instrucción
```

## Limitaciones

1. **Requiere Script:** Necesita que el script `s` esté funcionando correctamente para enviar mensajes
2. **Requiere Consola Abierta:** La consola de Claude Code debe permanecer abierta para recibir mensajes inyectados
3. **Sin Confirmación de Recepción:** Claude Code no sabe si el mensaje vía Telegram fue entregado exitosamente

## Beneficios

1. **Desarrollo Asíncrono:** Usuario puede dar instrucciones desde cualquier lugar
2. **Continuidad:** Claude Code trabaja continuamente sin necesidad de supervisión presencial
3. **Notificaciones:** Usuario recibe actualizaciones de progreso vía Telegram
4. **Flexibilidad:** Usuario puede enviar instrucciones vía Telegram sin interrumpir el flujo de trabajo
5. **Simplicidad:** Inyección directa de mensajes simplifica el workflow (sin comandos de espera)

## Mejores Prácticas

### Para Claude Code

1. **Reportar progreso periódicamente:** Usar `s gobot-send-message` durante tareas largas para mantener al usuario informado
2. **Indicar estado claramente:**
   - Si continúas trabajando → terminar con "Continúo trabajando"
   - Si terminaste TODO → terminar con "Tarea terminada, esperando instrucciones"
3. **Mensajes claros y concisos:** Ser específico sobre qué se hizo y qué sigue
4. **Reportar errores inmediatamente:** Si algo falla, notificar vía Telegram con descripción del problema y terminar con "Tarea terminada, esperando instrucciones"
5. **Reportar al terminar:** SIEMPRE enviar reporte final cuando completes todas las tareas solicitadas
6. **Compartir archivos cuando se solicite:** Si el usuario pide un archivo, usar el comando apropiado:
   - `s gobot-send-file` para archivos generales (logs, binarios, configs)
   - `s gobot-send-image` para imágenes y gráficos
   - `s gobot-send-video` para videos y demos
   - `s gobot-send-voice` para archivos de audio

### Para Usuario

1. **Instrucciones claras:** Ser específico en lo que se pide vía Telegram
2. **Confirmar recepción:** Si Claude Code reporta algo importante, responder para dar siguiente instrucción
3. **Mantener consola abierta:** La consola debe permanecer activa para recibir mensajes inyectados

## Resumen de Comandos

### Comandos Disponibles

```bash
# Enviar mensaje de texto
s gobot-send-message "mensaje aquí"

# Enviar archivo genérico (documento)
s gobot-send-file "/path/to/file"

# Enviar imagen (se muestra directamente en el chat)
s gobot-send-image "/path/to/image.png"

# Enviar video (reproducible en Telegram)
s gobot-send-video "/path/to/video.mp4"

# Enviar audio como mensaje de voz
s gobot-send-voice "/path/to/audio.mp3"
```

### Recepción de Mensajes

Los mensajes del usuario vía Telegram son inyectados automáticamente al chat local de Claude Code. No se requiere ningún comando especial para recibirlos.

## Cómo Funciona Internamente

1. **Usuario escribe en Telegram:** El mensaje se envía al bot de Telegram
2. **Sistema de inyección:** Un proceso detecta el mensaje y lo inyecta directamente al stdin de Claude Code
3. **Claude Code recibe:** El mensaje aparece en la consola como si el usuario lo hubiera escrito localmente
4. **Claude Code procesa:** Responde y ejecuta tareas normalmente
5. **Claude Code reporta:** Usa `s gobot-send-message` para enviar actualizaciones al usuario

---

**Documentado:** 2025-11-21
**Actualizado:** 2025-12-10
**Uso:** Activo en desarrollo actual
**Estado:** Operativo
