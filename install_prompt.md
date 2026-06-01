# Prompt de instalación de clical

Este archivo está pensado para ser pegado tal cual a un agente de IA con acceso
a una terminal (Claude Code, Cursor, etc.) para que instale **clical** de forma
autónoma en Linux, macOS o Windows.

---

## Prompt

Instalá **clical** desde el código fuente en esta máquina. clical es un sistema
de calendario CLI multiusuario escrito en Go. El repositorio oficial es
<https://github.com/neitanod/clical>.

Seguí estos pasos en orden y reportá al final un resumen breve de qué quedó
instalado y dónde.

### 1. Detectar el sistema operativo

Identificá si estás en Linux, macOS o Windows antes de empezar. Los pasos de
compilación e instalación cambian según la plataforma.

### 2. Verificar requisitos

clical necesita:

- **Go 1.23 o superior** — verificá con `go version`.
- **git** — verificá con `git --version`.
- En Linux/macOS además: **make** — verificá con `make --version`.

Si alguno falta, avisá al usuario y pedile que lo instale antes de continuar.
No intentes instalar Go automáticamente.

### 3. Clonar el repositorio

Elegí un directorio razonable según el sistema:

- Linux/macOS: `~/code/clical` o `~/src/clical`.
- Windows: `%USERPROFILE%\code\clical`.

Cloná con:

```
git clone https://github.com/neitanod/clical.git <ruta-elegida>
```

Si el directorio ya existe y tiene un clon válido, hacé `git pull --ff-only` en
vez de volver a clonar.

### 4. Compilar e instalar

**En Linux o macOS:**

```bash
cd <ruta-elegida>
make build           # genera ./clical en el directorio del repo
make install         # instala en $GOPATH/bin (suele ser ~/go/bin)
# Alternativa que requiere sudo:
# make install-system  # copia el binario a /usr/local/bin
```

**En Windows (PowerShell):**

```powershell
cd <ruta-elegida>
.\build.ps1          # genera .\clical.exe
.\install.ps1        # instala en %USERPROFILE%\go\bin
```

Si PowerShell bloquea el script por política de ejecución, usá:

```powershell
powershell -ExecutionPolicy Bypass -File .\build.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

### 5. Verificar que el binario esté en PATH

Después de instalar, abrí una terminal nueva y corré:

```
clical version
clical --help
```

Si el comando no se encuentra, asegurate de que el directorio de instalación
esté en `PATH`:

- Linux/macOS: agregá `export PATH="$HOME/go/bin:$PATH"` al `~/.bashrc` o
  `~/.zshrc` según corresponda, y recargá el shell.
- Windows: agregá `%USERPROFILE%\go\bin` al `PATH` de usuario con:

  ```powershell
  [Environment]::SetEnvironmentVariable('Path', "$env:Path;$env:USERPROFILE\go\bin", 'User')
  ```

  Y abrí una terminal nueva.

### 6. Crear el primer usuario (smoke test)

Pediéndole al usuario su nombre y timezone (o usando defaults razonables),
creá un usuario de prueba y listá:

```
clical user add --id=<id> --name="<nombre>" --timezone="<tz>"
clical user list
```

Confirmá que los archivos se hayan creado en `~/.clical/data/users/<id>/`
(en Windows: `%USERPROFILE%\.clical\data\users\<id>\`).

### 7. Reportar resultado

Al final, informá al usuario:

- Ruta del binario instalado.
- Ruta del directorio de datos (`~/.clical/data`).
- Si tuviste que modificar el `PATH` o pediste alguna acción manual.
- Cualquier advertencia o paso que no pudiste completar.

### Notas importantes

- **No modifiques `git config --global`.** Si necesitás identidad git para
  algún paso, usá `git config` local al repo y avisalo.
- **No instales dependencias del sistema sin confirmar** con el usuario
  (Go, make, gcc, etc.).
- Si el repositorio no tiene el archivo `cmd/clical/main.go` después del
  clone, no intentes crearlo: avisale al usuario que el repositorio puede
  estar en un estado inconsistente y detené la instalación.
- clical guarda datos por usuario en `~/.clical/data/`; no requiere base de
  datos ni servicios externos.
