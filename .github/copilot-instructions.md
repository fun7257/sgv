## sgv (Simple Go Version) – AI Contributor Quick Guide

Purpose: Lightweight Go version manager (install, switch, auto-switch, list, uninstall, latest, sub‑version listing) with per-version environment variable management. Focus on predictable file layout, minimal deps, clear CLI UX, and seamless shell integration.

### Architecture & Flow
1. **Entry**: `main.go` -> `cmd/root.go` (Cobra). Root command also acts as implicit "install & switch" when passed a version arg.
2. **Command Layer**: Lives in `cmd/` - each file defines and `init()`-registers a `*cobra.Command`. Keep new commands self‑contained; side effects only in `Run/RunE`. Notable: `env.go` manages per-version environment variables with flags `-w`, `-u`, `--shell`, `--clean`, `-a`.
3. **Core Logic**: Split under `internal/`:
   - `internal/config`: One-time `Init()` (run via `cobra.OnInitialize`) sets paths (`SgvRoot=~/.sgv`, `VersionsDir`, `CurrentSymlink`) and resolves `DownloadURLPrefix` (env override `SGV_DOWNLOAD_URL_PREFIX`, ensure trailing '/').
   - `internal/version`: Discovery (local via dir read, remote via JSON API + cache), switching (symlink), version metadata (build info via `debug.ReadBuildInfo`). Remote fetch creates per (OS, Arch) entries; Windows filtered out.
   - `internal/installer`: Download + extract tarball to `VersionsDir/<goX.Y.Z>` with progress bar, leaves nested `go/` directory intact (symlink points to that subdir).
   - `internal/env`: **Critical component** - manages per-version environment variables stored as `~/.sgv/env/<version>.env` files. Handles protected variables (GOROOT, GOPATH, etc.), atomic file operations, and shell output generation.
4. **Data Flow**: CLI arg -> normalize version (ensure `go` prefix) -> validate support (>=1.13) -> optional `go.mod` compatibility gate -> install if missing -> `version.SwitchToVersion` replaces symlink -> shell wrapper auto-loads env vars.
5. **Shell Integration**: `install.sh` creates a `sgv()` wrapper function that intercepts commands and auto-runs `eval $(sgv env --shell --clean)` after successful version switches, env modifications, or auto/latest commands.

### Environment Variable System (Critical)
* **Storage**: Each Go version gets `~/.sgv/env/<version>.env` with `KEY=VALUE` format.
* **Protection**: Built-in Go vars (GOROOT, GOPATH, GOPROXY, etc.) are protected from modification.
* **Shell Integration**: `sgv env --shell` outputs `export` commands; `--clean` also outputs `unset` for variables from other versions to prevent conflicts.
* **Commands**: `sgv env -w KEY=VALUE` (set), `sgv env -u KEY` (unset), `sgv env --clear` (clear all), `sgv env -a` (list all versions).
* **Atomic Operations**: Uses temp files + rename for safe concurrent access with file mutex.

### Version Handling Patterns
* **Normalization**: Always store with full patch (e.g. `go1.22.1`). Commands accept both `1.22.1` and `go1.22.1` formats.
* **Local Discovery**: Directory names in `VersionsDir` (ONLY directories count). Current version from symlink target parent.
* **Symlink Model**: `CurrentSymlink` -> `<VersionsDir>/<version>/go`; never point directly at binary.
* **Remote Cache**: 10-minute cache in temp dir (`sgv-remote-versions-cache.json`). Stale cache used on API failure.

### Key Conventions
* **Error Handling**: Wrap with context `fmt.Errorf("<action>: %w", err)`. CLI layer prints to stderr + `os.Exit(1)`. Library packages return errors only.
* **Interactive Prompts**: Simple `fmt.Scanln` or `bufio.Reader`; keep terse (see `auto.go`, `uninstall.go`, `env.go --clear`).
* **Version Comparisons**: Always normalize to semver with `normalizeGoVersion` before `semver.Compare`.
* **Mutex Usage**: File operations in `internal/env` use `fileMutex` for concurrent safety.

### Critical Developer Workflows
* **Local Build**: `make local` (builds + copies to `/usr/local/bin` with sudo)
* **Development**: `make build` -> `./sgv <command>` for testing
* **Shell Function**: Test with real shell functions in `install.sh` - regex patterns are tricky
* **Env Testing**: Check `internal/env/env_test.go` for comprehensive test patterns including edge cases

### Shell Integration Patterns
The `install.sh` script generates a `sgv()` wrapper function with these trigger patterns:
```bash
# Version switch: sgv go1.22.1 or sgv 1.22.1
if [ $# -eq 1 ] && [[ "$1" =~ ^(go)?[0-9]+\.[0-9]+(\.[0-9]+)?$ ]]

# Env modification: sgv env -w KEY=VALUE, sgv env -u KEY  
elif [ "$1" = "env" ] && { [ "$2" = "-w" ] || [ "$2" = "--write" ] || [ "$2" = "-u" ] || [ "$2" = "--unset" ]; }

# Auto commands: sgv auto, sgv latest
elif [ "$1" = "auto" ] || [ "$1" = "latest" ]
```

### External Dependencies (Keep Minimal)
* `github.com/spf13/cobra` – CLI framework
* `github.com/samber/lo` – functional helpers  
* `github.com/schollz/progressbar/v3` – download UI
* `golang.org/x/mod/semver` – version comparison

### Adding New Commands
1. Create `cmd/newcmd.go` with `cobra.Command`
2. Register in `init()`: `rootCmd.AddCommand(newCmd)`
3. Business logic in `internal/` if complex
4. Follow env var loading patterns if command affects versions
5. Update shell wrapper in `install.sh` if auto-loading needed

### Critical Pitfalls
* **Never uninstall active version** (guard in `uninstall.go`)
* **Windows unsupported** (explicit checks in `root.go`, `installer.go`)
* **Shell regex escaping** in `install.sh` - test thoroughly
* **Protected env vars** - check `env.IsProtectedVar()` before modifications
* **File operations** - use atomic patterns from `env.SaveEnvVars()`

---
Questions or unclear patterns? Point to specific file + line for context before large changes.
