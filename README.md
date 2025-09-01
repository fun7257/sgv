[English](README.md) | [简体中文](README_zh.md)
# sgv - Go Version Manager

`sgv` (Simple Go Version) is a lightweight CLI tool for managing multiple Go versions on your system. It supports Go 1.13 and above, and works on macOS and Linux.

## Features

- **Install Go versions**: Download and install any supported Go version.
- **Switch Go versions**: Instantly switch between installed Go versions.
- **Auto switch**: Automatically switch to the required Go version for the current project based on `go.mod`.
- **Get latest**: Install and switch to the latest Go version with one command.
- **Per-version environment variables**: Manage project-specific environment variables for each Go version with automatic loading.
- **List installed versions**: View all Go versions installed by sgv, grouped by major version.
- **List available patch versions**: List all available patch versions for a given major version, and see which are installed.
- **Uninstall Go versions**: Remove any installed Go version (except the currently active one).
- **Show sgv version**: Display the sgv build version and commit hash.
- **Seamless shell integration**: Automatic environment variable loading with no manual intervention required.

---

## Installation

### macOS / Linux

```bash
curl -sSL https://raw.githubusercontent.com/fun7257/sgv/main/install.sh | bash
```

- Installs to `/usr/local/bin/sgv`
- Automatically configures `GOROOT` and `PATH` in your `~/.bashrc` or `~/.zshrc`
- After installation, restart your terminal or run `source ~/.bashrc` or `source ~/.zshrc`

---

## Usage

### Switch or Install & Switch Go Version

```bash
sgv <version>
```
- Example: `sgv 1.22.1` or `sgv go1.21.0`
- If not installed, sgv will download and install the version, then switch
- If in a Go project and the requested version is lower than `go.mod` requires, the operation will abort with an error

### Install Only (Do Not Switch)

```bash
sgv <version> --no-switch
```
- Example: `sgv 1.22.1 --no-switch`
- Downloads and installs the version, but does not switch to it.
- If the version is already installed, it will do nothing.

### Auto Switch (Based on go.mod)

```bash
sgv auto
```
- Detects the required Go version from `go.mod` (prefers `toolchain` if present and higher)
- If not installed, prompts to download and install
- If already active, does nothing
- If not in a Go project, prints a message and does nothing

### Get and Switch to Latest Go Version

```bash
sgv latest
```
- Installs the latest Go version if not present, and switches to it

### List Installed Go Versions

```bash
sgv list
```
- Shows all installed versions, grouped by major version
- The current version is marked with `<- current`

### List All Patch Versions for a Major Version

```bash
sgv sub <major_version>
```
- Example: `sgv sub 1.22`
- Lists all available Go 1.22.x versions, with installed ones marked `(installed)`
- Only available for Go 1.13 and above

### Uninstall a Go Version

```bash
sgv rm <version...>
```
- Example 1 (specific versions): `sgv rm 1.22.1 1.21.7`
- Example 2 (major version): `sgv rm 1.22` (removes all installed 1.22.x versions)
- Cannot uninstall the currently active version.

### Show sgv Version

```bash
sgv version
```
- Shows the Go version used to build sgv and its commit hash

### Manage Environment Variables

sgv provides sophisticated per-version environment variable management:

```bash
sgv env                          # List environment variables for current Go version
sgv env -w KEY=VALUE             # Set environment variable for current Go version
sgv env -u KEY                   # Remove environment variable for current Go version
sgv env -a                       # List all Go versions with their environment variables
sgv env --clear                  # Clear all environment variables for current Go version
sgv env --shell                  # Output environment variables in shell format
sgv env --shell --clean          # Output with cleanup of conflicting variables
```

**Examples:**
```bash
sgv env -w GOWORK=auto           # Enable Go workspace mode
sgv env -w GODEBUG=gctrace=1     # Enable GC trace debugging
sgv env -w CGO_ENABLED=0         # Disable CGO for this Go version
sgv env -u GODEBUG               # Remove GODEBUG setting
sgv env -a                       # See all versions and their variables
```

**Key Features:**
- **Version isolation**: Each Go version has its own environment variables stored separately
- **Automatic loading**: Environment variables are automatically applied when switching versions
- **Protected variables**: Critical Go variables (GOROOT, GOPATH, GOPROXY, etc.) cannot be modified through sgv
- **Conflict prevention**: The `--clean` flag prevents variable conflicts between versions
- **Persistent storage**: Variables are saved in `~/.sgv/env/<version>.env` and restored automatically
- **Shell integration**: Changes are immediately applied to your current shell session

---

## Seamless Experience

sgv provides a seamless experience with automatic environment loading through intelligent shell integration:

### Automatic Environment Loading
- **Version switching**: `sgv 1.22.1` or `sgv go1.21.0` automatically loads environment variables
- **Environment changes**: `sgv env -w KEY=VALUE` and `sgv env -u KEY` immediately apply to your current shell
- **Auto commands**: `sgv auto` and `sgv latest` automatically load environment variables after version switches
- **Flexible version format**: Supports both `1.22.1` and `go1.22.1` formats

### How It Works
The installation script creates a wrapper function in your shell that:
1. Executes the actual `sgv` command
2. Detects successful operations that affect environment variables
3. Automatically runs `eval $(sgv env --shell --clean)` to update your shell
4. Prevents conflicts by cleaning variables from other versions

### No Manual Intervention
- No need to run `eval` commands manually
- No need to restart your terminal
- No need to source configuration files
- Environment variables are immediately available after any change

This sophisticated shell integration makes `sgv` feel like a native part of your development environment.

---

## Configuration

### Environment Variables

- `SGV_DOWNLOAD_URL_PREFIX`  
  Change the Go download source (e.g., for China mainland users)

```bash
export SGV_DOWNLOAD_URL_PREFIX=https://golang.google.cn/dl/
```

Set this before running sgv commands, or add to your shell profile for persistence.

### File Structure

sgv organizes files in a predictable way:

- `~/.sgv/versions/` - All installed Go versions (e.g., `~/.sgv/versions/go1.22.1/`)
- `~/.sgv/current` - Symlink to the currently active Go version
- `~/.sgv/env/` - Environment variable files (e.g., `~/.sgv/env/go1.22.1.env`)

### Shell Integration

The installation script automatically adds:
- Environment variable exports (`GOROOT`, `PATH`)
- A wrapper function that provides automatic environment loading
- Session startup code to restore environment variables

This configuration is added to `~/.bashrc` or `~/.zshrc` depending on your shell.

---

## Notes

- **Cross-platform support**: Works on macOS and Linux (Windows is not supported)
- **Go version support**: Supports Go 1.13 and above
- **File organization**: All Go versions are installed under `~/.sgv/versions/`, with the current version symlinked as `~/.sgv/current`
- **Environment isolation**: Each Go version maintains its own set of environment variables
- **Automatic configuration**: The install script handles all `GOROOT` and `PATH` configuration automatically
- **Shell compatibility**: Works with both bash and zsh shells

---

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License. See [LICENSE](./LICENSE) for details.
