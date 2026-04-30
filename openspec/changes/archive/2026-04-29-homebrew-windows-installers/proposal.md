# Proposal: Homebrew + Windows Installers

## Intent

Habilitar instalación de Gentle-QA en macOS, Linux y Windows mediante:
- Homebrew tap (macOS + Linux)
- Scoop bucket (Windows)
- Script instalador actualizado (todos los OS)

El usuario tiene `HOMEBREW_TAP_TOKEN` configurado y quiere un instalador funcional desde Homebrew.

## Scope

### In Scope
- Agregar soporte Windows a `scripts/install.sh` (detectar MINGW/MSYS/CYGWIN)
- Actualizar lógica de descarga para Windows (.zip + .exe vs .tar.gz)
- Ejecutar GoReleaser release inicial (v0.1.0) para generar binaries
- Verificar fórmula Homebrew generada en `homebrew-tap/Formula/`
- Configurar Scoop bucket con JSON para Windows
- Testear instalación en macOS/Linux

### Out of Scope
- GitHub Actions workflow (puede agregarse después)
- Build from source en Homebrew (solo binary)
- Actualización automática de versiones (Homebrew maneja esto)

## Approach

1. **Windows support en install.sh**: Modificar `detect_platform()` para reconocer Windows, usar `.zip` y `.exe`
2. **GoReleaser release**: Taggear v0.1.0 y ejecutar `goreleaser release --clean`
3. **Homebrew**: GoReleaser genera fórmula automáticamente en `homebrew-tap`
4. **Scoop**: GoReleaser genera JSON o crear manualmente en `scoop-bucket`

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `scripts/install.sh` | Modified | Agregar Windows detection, .zip/.exe handling |
| `.goreleaser.yaml` | Verified | Ya configurado correctamente |
| ` EduardoVeraE/homebrew-tap` | New | Fórmula auto-generada por GoReleaser |
| ` EduardoVeraE/scoop-bucket` | New | JSON auto-generado por GoReleaser |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| GoReleaser no tiene acceso push a homebrew-tap | Medium | Generar fórmula manualmente con `brew create` |
| Token no tiene permisos correctos | Low | Verificar permisos antes de ejecutar |
| Windows binary no funciona en GitHub Actions | Low | Testear con matrix de OS |

## Rollback Plan

- Revertir cambios en `install.sh` con git
- Eliminar tags de GitHub (`git tag -d v0.1.0`, `git push origin :refs/tags/v0.1.0`)
- Eliminar fórmula de homebrew-tap manualmente

## Dependencies

- `HOMEBREW_TAP_TOKEN` configurado en GitHub (CONFIRMADO por usuario)
- GoReleaser instalado (`brew install goreleaser`)

## Success Criteria

- [ ] `brew install gentle-qa` funciona en macOS
- [ ] `brew install gentle-qa` funciona en Linux
- [ ] Scoop install funciona en Windows
- [ ] `install.sh` detecta y instala en Windows
- [ ] Todos los binaries (darwin/linux/windows + amd64/arm64) disponibles en GitHub Releases