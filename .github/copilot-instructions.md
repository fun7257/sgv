# Copilot Instructions for sgv (Simple Go Version)

## Project Overview
- **sgv** is a CLI tool for managing multiple Go versions, supporting install, switch, auto-switch (via `go.mod`), list, and uninstall operations.
- Main entrypoint: `main.go` delegates to subcommands in `cmd/`.
- Core logic is organized in `internal/`:
  - `internal/version/`: Go version management (fetch, install, switch, list, uninstall)
  - `internal/config/`: Paths, environment, and configuration
  - `internal/installer/`: Download and extraction logic

## Key Patterns & Conventions
- **Go version info** is fetched from the official Go website (see `GetRemoteVersions` in `internal/version/version.go`).
- **Installed versions** are managed in a directory (`config.VersionsDir`), with the active version symlinked by `config.CurrentSymlink`.
- **Switching versions** updates the symlink and environment variables.
- **Auto-switch** uses `go.mod` in the current directory to determine the required Go version.
- **Error handling**: Always wrap errors with context (see `fmt.Errorf("...: %w", err)`).
- **Sorting**: Installed and remote versions are always sorted before display.
- **External dependencies**: Uses `github.com/samber/lo` for functional utilities.

## Developer Workflows
- **Build**: `go build -o sgv main.go`
- **Test**: No explicit test suite found; add tests in `internal/` as needed.
- **Install**: Use `install.sh` (Unix) or `install.ps1` (Windows) for end-to-end install experience.
- **Debug**: Use `fmt.Printf` or `log` for debugging; no custom logger.
- **Release**: Update `sgvVersion` and `sgvCommit` in `internal/version/version.go` as needed.

## Examples
- To add a new subcommand: create a new file in `cmd/`, register it in `root.go`.
- To fetch all available Go versions: use `FetchAllGoVersions()` in `internal/version/version.go`.
- To change install location or symlink logic, update `internal/config/config.go`.

## Integration Points
- Downloads Go binaries from the official Go CDN (see `config.DownloadURLPrefix`).
- Modifies shell profile files to update `PATH` and `GOROOT` on install.

## Project-Specific Advice
- Always use the config constants for paths and URLs; do not hardcode.
- When adding features, follow the subcommand pattern in `cmd/` and keep business logic in `internal/`.
- Keep CLI output user-friendly and concise.

---
For more details, see `README.md` and code comments in `internal/version/version.go` and `cmd/`.
