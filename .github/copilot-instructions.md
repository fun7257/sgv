## sgv (Simple Go Version) – AI Contributor Quick Guide

Purpose: Lightweight Go version manager (install, switch, auto-switch, list, uninstall, latest, sub‑version listing). Focus on predictable file layout, minimal deps, clear CLI UX.

### Architecture & Flow
1. Entry: `main.go` -> `cmd/root.go` (Cobra). Root command also acts as implicit "install & switch" when passed a version arg.
2. Command set lives in `cmd/`: each file defines and `init()`-registers a `*cobra.Command` (e.g. `auto.go`, `install.go`, `list.go`, `sub.go`, `uninstall.go`, `version.go`). Keep new commands self‑contained; side effects only in `Run/RunE`.
3. Core logic split under `internal/`:
  - `internal/config`: one-time `Init()` (run via `cobra.OnInitialize`) sets paths (`SgvRoot=~/.sgv`, `VersionsDir`, `CurrentSymlink`) and resolves `DownloadURLPrefix` (env override `SGV_DOWNLOAD_URL_PREFIX`, ensure trailing '/'). Never duplicate these path computations elsewhere.
  - `internal/version`: discovery (local via dir read, remote via JSON API + cache), switching (symlink), version metadata (build info via `debug.ReadBuildInfo`). Remote fetch creates per (OS, Arch) entries; Windows filtered out.
  - `internal/installer`: download + extract tarball to `VersionsDir/<goX.Y.Z>` with a progress bar, then leaves nested `go/` directory intact (symlink points to that subdir).
4. Data flow (install & switch path): CLI arg -> normalize version (ensure `go` prefix) -> validate support (>=1.13) -> optional `go.mod` compatibility gate (semver compare via `golang.org/x/mod/semver`) -> install if missing (`installer.Install`) -> `version.SwitchToVersion` replaces `~/.sgv/current` symlink -> user shell (preconfigured by `install.sh`) picks new `GOROOT` binaries.
5. Auto mode (`auto.go`): parses `go.mod` lines for `go` and `toolchain` directives (prefers higher); normalizes `>=1.21` minor to `.0`; prompts user before switching; reuses root command programmatically.

### Key Conventions / Patterns
* Always store versions with full patch: e.g. `go1.22.1` (auto-normalization ensures this for >=1.21 when coming from `go.mod`).
* Local versions = directory names inside `VersionsDir` (ONLY directories count). Current version derived from parent directory of symlink target. Don't parse `GOROOT`.
* Symlink model: `CurrentSymlink` -> `<VersionsDir>/<version>/go`; never point directly at binary.
* Remote versions API: `GET <DownloadURLPrefix>?mode=json&include=all`; cached (temp file) for 10 minutes (`internal/version/cache.go`). On fetch error: try stale cache before failing.
* Platform filter: ignore `windows`; user platform determined at runtime for compatibility display (see `cmd/sub.go`).
* Error style: wrap with context `fmt.Errorf("<action>: %w", err)`; user-facing errors print to stderr then `os.Exit(1)` in command layer only (library packages return errors). Avoid exiting inside `internal/*`.
* Version comparisons: always normalize to semver with `normalizeGoVersion` before `semver.Compare`.
* Interactive confirmations: simple `fmt.Scanln` or `bufio.Reader`; keep prompts terse (see `auto.go`, `uninstall.go`).

### External Dependencies (Keep Minimal)
* `github.com/spf13/cobra` – CLI framework.
* `github.com/samber/lo` – simple functional helpers (used for filtering stable versions).
* `github.com/schollz/progressbar/v3` – download progress UI.
* `golang.org/x/mod/semver` – semver comparison.
Avoid adding heavy deps; prefer stdlib where possible.

### Developer Workflows
* Build (local dev): `go build -o sgv .` or `make build` (uses root module; binary name must stay `sgv`).
* Install to PATH (dev): `go install .` or run `install.sh` for full user setup (adds env lines to shell config). Release tarballs expected as `sgv_<version>_<os>_<arch>.tar.gz` consumed by `install.sh`.
* Releasing: bump Git tag; runtime self-report uses `debug.ReadBuildInfo` (no manual version constant editing needed unless changing default placeholders). Ensure CI build injects proper VCS info.
* Debug: add temporary `fmt.Printf` statements; no structured logger present.
* Cache invalidation (manual): delete temp file returned by `VersionCache.GetCacheInfo()` (location: system temp dir, filename `sgv-remote-versions-cache.json`).

### Adding a New Command (Example)
1. Create `cmd/foo.go` defining `var fooCmd = &cobra.Command{ ... }`.
2. In its `init()`, call `rootCmd.AddCommand(fooCmd)`.
3. Put business logic in a new/internal package if it grows; keep CLI layer narrow (arg parsing, I/O, exit codes).
4. Follow normalization rules if accepting version input (ensure prefix `go`, validate support, maybe semver compare if interacting with `go.mod`).

### Common Pitfalls / Edge Cases
* Do NOT attempt to uninstall the active version (guard exists in `uninstall.go`). If changing logic, preserve this check before deletion.
* Windows must remain unsupported (explicit exits in `root.go` and `installer.Install`). Any platform expansion should centralize gate logic in `checkPlatformSupport()`.
* `auto` command should emit no output when already on correct version (intentional quiet success). Preserve this to avoid noisy shells.
* When modifying extraction, keep per-file close semantics to avoid "too many open files" (see comment in `extractTarGz`).
* Maintain trailing slash enforcement for custom `SGV_DOWNLOAD_URL_PREFIX`.

### Where To Look For Examples
* Version normalization & compatibility: `cmd/util.go`.
* Remote fetch + cache strategy: `internal/version/version.go` + `cache.go`.
* Installation & extraction details: `internal/installer/installer.go`.
* Grouped listing + coloring: `cmd/list.go` and `cmd/sub.go` (ANSI + `fatih/color`).

### Safe Extension Ideas
* Add a `cache` subcommand to show/clear remote version cache (reusing `VersionCache` methods).
* Add a `which` command to print active Go binary path (`CurrentSymlink/bin/go`).

Keep changes small, reuse existing helpers, and surface user-visible alterations via clear CLI messages.

---
Questions or unclear patterns? Point to the specific file + line context and propose an adjustment before large refactors.
