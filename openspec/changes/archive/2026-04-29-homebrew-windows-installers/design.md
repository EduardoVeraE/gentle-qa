# Design: Homebrew + Windows Installers

## Technical Approach

Modificar el script de instalación existente (`scripts/install.sh`) para detectar y manejar plataformas Windows usando la API nativa de GoReleaser para archivos ZIP y extensión `.exe`.

## Architecture Decisions

### Decision: Windows Detection via OSTYPE

**Choice**: Detectar Windows revisando `uname -s` para valores `MINGW*`, `MSYS*`, o `CYGWIN*`
**Alternatives considered**: Revisar `OSTYPE` environment variable o `WINDEEN` variable
**Rationale**: `uname -s` es más confiable en entornos MSYS2/MINGW, estos entornos setean esa variable consistentemente

### Decision: Formato de Archivo por Plataforma

**Choice**: Usar `.zip` para Windows (como GoReleaser config), `.tar.gz` para Unix
**Alternatives considered**: Siempre usar `.tar.gz` y renombrar en Windows
**Rationale**: GoReleaser ya genera `.zip` para Windows, mantener consistencia

### Decision: Ruta de Instalación en Windows

**Choice**: Usar `%LOCALAPPDATA%\Programs\{BINARY_NAME}` como primario, `%USERPROFILE%\.local\bin` como fallback
**Alternatives considered**: `%PROGRAMFILES%`, `%APPDATA%`
**Rationale**: `%LOCALAPPDATA%\Programs` es el estándar moderno para installation per-user en Windows 10+

## Data Flow

```
detect_platform() → get_archive_name() → install_binary()
       ↓                    ↓                  ↓
   OS=windows         .zip ext          unzip + cp
                   get_binary_name()       ↓
                        ↓                  ↓
                 gentle-qa.exe         add to PATH
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `scripts/install.sh` | Modified | Agregado soporte Windows |

## Interfaces / Contracts

No hay cambios a interfaces externas. El script mantiene backwards compatibility:
- Mismos flags: `--method`, `--dir`, `-h/--help`
- Mismos exit codes

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Manual | install.sh en Windows (MSYS2/MINGW) | Usuario debe ejecutar |

## Migration / Rollout

No migration requerida. Solo cambios en script de instalación.

## Open Questions

- [ ] ¿El usuario puede testear en Windows?
- [ ] ¿Hay problemas con el path de Windows en PATH environment?

## Notes

Las tareas 1.1-1.5 fueron completadas por el orquestador. Las tareas 2.x requieren:
1. Taggear `v0.1.0`
2. Ejecutar GoReleaser
3. Acceso con token a homebrew-tap y scoop-bucket