# Changelog

## 0.1.0 (2026-04-19)

### Initial Release
- Forked from Pachyderm 2.x
- Removed all enterprise and license features
- Removed Bazel build system; use standard Go modules
- Migrated to Standard Go Project Layout:
  - `/cmd`: Main commands
  - `/internal`: Private packages
  - `/pkg`: Public packages
- Renamed project to "Mammoth"
