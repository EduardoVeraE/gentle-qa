# Tasks: homebrew-windows-installers

## Phase 1: Install Script Updates ✅ COMPLETADO

- [x] 1.1 Modificar `detect_platform()` en `scripts/install.sh` para reconocer Windows (MINGW/MSYS/CYGWIN)
- [x] 1.2 Actualizar `get_archive_name()` para usar `.zip` en Windows
- [x] 1.3 Actualizar lógica de descarga para Windows (.zip + .exe)
- [x] 1.4 Actualizar lógica de instalación para rutas Windows
- [x] 1.5 Agregar verificación de instalación para Windows
- [x] 1.6 Crear `install.ps1`para Windows (PowerShell installer)
- [x] 1.7 Verificar README.md con instrucciones Scoop y Homebrew

## Phase 2: GoReleaser Release ✅ COMPLETADO

- [x] 2.1 Taggear release inicial v0.1.0
- [x] 2.2 Ejecutar `goreleaser release --clean`
- [x] 2.3 Verificar binaries subidos a GitHub Releases

## Phase 3: Homebrew Verification ✅ COMPLETADO (v1.22.0)

- [x] 3.1 Verificar fórmula generada en homebrew-tap/Formula/
- [x] 3.2 Testear `brew install gentle-qa` en macOS (o documentar que el usuario debe hacerlo)

## Phase 4: Scoop Setup ✅ COMPLETADO (v1.22.0)

- [x] 4.1 Verificar JSON generado en scoop-bucket/bucket/
- [x] 4.2 Testear `scoop install gentle-qa` en Windows (o documentar que el usuario debe hacerlo)