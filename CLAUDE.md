# CLAUDE.md

Este archivo proporciona orientación a Claude Code (claude.ai/code) cuando trabaja con el código de este repositorio.

## Descripción del Proyecto

**clical** es un sistema de calendario multiusuario desarrollado en Go, actualmente en fase inicial de desarrollo.

## Estado del Proyecto

**Fase Actual**: Desarrollo inicial - Definiendo arquitectura base

El proyecto está comenzando su desarrollo:
- ⏳ Estructura de proyecto por establecer
- ⏳ Stack tecnológico por definir
- ⏳ Arquitectura en diseño

## Lenguaje y Comunicación

**Todas las interacciones deben ser en español (castellano).** La documentación del proyecto, mensajes de commit y journal de desarrollo se mantienen en español.

## Flujo de Trabajo de Desarrollo

### Desarrollo Basado en Sesiones

Este proyecto sigue un enfoque de desarrollo estructurado con los siguientes principios:

1. **Implementación Iterativa**: Las características se desarrollan incrementalmente a través de interacción conversacional
2. **Control de Versiones**: Cuando una característica alcanza un estado satisfactorio (indicado por el usuario), crear un commit de git
3. **Journal de Desarrollo**: Después de cada commit, documentar el trabajo en una entrada del journal en `ai/journal/`
4. **Actualización de Documentación**: Después de cada commit, actualizar la documentación en `docs/` según corresponda

### Estructura del Journal

Las entradas del journal se almacenan en `ai/journal/` con la siguiente convención de nombres:
- Formato: `YYYY-MM-DD_Session_Summary.md`
- Para sesiones largas: `YYYY-MM-DD_Session_Summary_Part_2.md`, etc.

Ejemplo de directorio journal:
```
ai/journal/
├── 2025-11-19_Session_Summary.md
├── 2025-11-19_Session_Summary_Part_2.md
└── 2025-11-20_Session_Summary.md
```

### Especificaciones del Proyecto

Las especificaciones y decisiones arquitectónicas están documentadas en `ai/specs/`:
- `00_Overview.md` - Descripción general del proyecto
- `01_Journaling.md` - Guías de journaling
- (Se añadirán más según se desarrolle el proyecto)

**Al iniciar una sesión:**
1. Leer `ai/README.md` para instrucciones de workflow
2. Revisar todos los archivos en `ai/specs/` para especificaciones del proyecto
3. Revisar `ai/journal/` para entender qué se ha implementado
4. Revisar el último journal para conocer el estado más reciente

### Estructura de Documentación

La documentación de usuario se mantiene en `docs/`:
- `index.html` - Página principal con tabla de contenidos y búsqueda
- Contenidos cargados vía AJAX desde archivos separados
- Actualizar después de cada commit con nuevas características

La estructura consiste en:
- **Header**: Encabezado superior
- **Columna izquierda**: Campo de búsqueda y tabla de contenidos
- **Sección central/derecha**: Contenido de documentación cargado dinámicamente

## Arquitectura del Proyecto

### Estructura de Directorios (Propuesta)

```
clical/
├── cmd/
│   └── clical/             # Punto de entrada principal
├── pkg/                    # Paquetes públicos exportables
│   └── (por definir)
├── internal/               # Paquetes internos
│   └── (por definir)
├── ai/                     # Documentación de desarrollo
│   ├── specs/              # Especificaciones técnicas
│   └── journal/            # Journal de desarrollo
├── docs/                   # Documentación web
├── .env                    # Configuración (no incluir en git)
├── .env.example            # Ejemplo de configuración
└── README.md               # Documentación principal
```

### Stack Tecnológico (Por Definir)

- **Lenguaje**: Go 1.24+
- **Librerías**: Por determinar según necesidades

## Prácticas de Git

### Commits

- **Nunca commitear automáticamente** - Solo crear commits cuando el usuario indica explícitamente que una característica está completa
- Los mensajes de commit deben estar en español
- **NO incluir** atribuciones automáticas de IA (Co-Authored-By) a menos que el usuario lo solicite explícitamente

### Formato de Commits

```
Título descriptivo del cambio

Descripción detallada del cambio si es necesario.
Puede incluir múltiples párrafos.
```

**Nota**: No se incluye automáticamente la atribución "Generated with Claude Code" ni "Co-Authored-By", a menos que el usuario lo pida.

### Seguridad en Git

1. **Nunca commitear tokens o secretos** - Usar .env.example para ejemplos
2. **Verificar .gitignore** antes de commits
3. **Revisar cambios** antes de commitear

## Testing y Building

### Compilación

```bash
# Compilar el proyecto
go build -o clical ./cmd/clical

# Compilar para producción con optimizaciones
go build -ldflags="-s -w" -o clical ./cmd/clical
```

### Testing

```bash
# Ejecutar todos los tests
go test ./...

# Test con cobertura
go test -cover ./...

# Test de módulo específico
go test ./pkg/<modulo> -v
```

### Ejecución en Desarrollo

```bash
# Ejecución directa
go run ./cmd/clical

# Con hot-reload (requiere air)
air
```

## Directrices de Desarrollo

### Al Agregar Nuevas Características

1. **Revisar especificaciones existentes** en `ai/specs/` para entender el diseño
2. **Mantener consistencia** con la arquitectura actual
3. **Documentar cambios** en specs si se toman decisiones arquitectónicas
4. **Actualizar README.md** con nuevas características
5. **Crear entrada en journal** después del commit
6. **Actualizar documentación web** en `docs/` si aplica

### Al Modificar Código Existente

1. **Leer el código existente** primero para entender el patrón
2. **Mantener el estilo** de código Go (gofmt, convenciones)
3. **No romper compatibilidad** sin discutir con el usuario
4. **Actualizar documentación** relacionada

### Seguridad

1. **Validar entrada de usuarios** en todas las interfaces
2. **Usar prácticas seguras** de manejo de datos
3. **Sanitizar salidas** para prevenir inyecciones
4. **Gestión segura de credenciales** (nunca en código)

### Testing

1. **Compilar antes de commitear** - `go build ./cmd/clical`
2. **Test manual** de funcionalidades nuevas o modificadas
3. **Verificar logs** no contengan errores
4. **Test de casos edge** cuando sea relevante

## Configuración

### Archivo .env (Por Definir)

```bash
# Configuración de ejemplo
# (Se definirá según las necesidades del proyecto)
```

## Características Planificadas

(Por definir según avance el proyecto)

## Próximos Pasos

Áreas de trabajo inmediatas:
- [ ] Definir requerimientos completos del sistema
- [ ] Establecer arquitectura base
- [ ] Elegir stack tecnológico
- [ ] Crear estructura inicial del proyecto
- [ ] Configurar sistema de build

## Recursos Adicionales

- **Go Docs**: https://go.dev/doc/
- **Effective Go**: https://go.dev/doc/effective_go
- (Más recursos según se definan dependencias)

## Notas Importantes

1. Este proyecto usa **Go Modules** - mantener `go.mod` actualizado
2. Seguir convenciones estándar de Go para estructura de proyecto
3. Priorizar simplicidad y mantenibilidad
4. Documentar decisiones arquitectónicas en `ai/specs/`

---

**Última actualización**: 2025-11-19
**Versión del proyecto**: v0.0 (desarrollo inicial)
**Estado**: Desarrollo - Fase de diseño e implementación inicial
