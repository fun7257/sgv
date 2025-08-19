[English](README.md) | [简体中文](README_zh.md)
# sgv - Go Version Manager

`sgv` (Simple Go Version) is a lightweight CLI tool for managing multiple Go versions on your system. It supports Go 1.13 and above, and works on macOS and Linux.

## Features

- **Install Go versions**: Download and install any supported Go version.
- **Switch Go versions**: Instantly switch between installed Go versions.
- **Auto switch**: Automatically switch to the required Go version for the current project based on `go.mod`.
- **Get latest**: Install and switch to the latest Go version with one command.
- **Environment variables**: Manage project-specific environment variables per Go version.
- **List installed versions**: View all Go versions installed by sgv, grouped by major version.
- **List available patch versions**: List all available patch versions for a given major version, and see which are installed.
- **Uninstall Go versions**: Remove any installed Go version (except the currently active one).
- **Show sgv version**: Display the sgv build version and commit hash.

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
sgv install <version>
```
- Example: `sgv install 1.22.1`
- Downloads and installs the version, but does not switch

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
sgv uninstall <version>
```
- Example: `sgv uninstall 1.22.1`
- Cannot uninstall the currently active version

### Show sgv Version

```bash
sgv version
```
- Shows the Go version used to build sgv and its commit hash

### Manage Environment Variables

```bash
sgv env                           # List environment variables for current Go version
sgv env -w KEY=VALUE             # Set environment variable for current Go version
sgv env -u KEY                   # Remove environment variable for current Go version
sgv env --shell                  # Output environment variables in shell format
```

**Examples:**
```bash
sgv env -w GOWORK=auto           # Enable Go workspace mode
sgv env -w GODEBUG=gctrace=1     # Enable GC trace debugging
sgv env -u GODEBUG               # Remove GODEBUG setting
```

**Key Features:**
- **Version isolation**: Each Go version has its own environment variables
- **Seamless loading**: Environment variables are automatically applied when switching versions or modifying them
- **Protected variables**: Critical Go variables (GOROOT, GOPATH, etc.) cannot be modified
- **Persistent storage**: Variables are saved per version and restored automatically

---

## Seamless Experience

sgv provides a seamless experience with automatic environment loading:

- **Automatic switching**: When you run `sgv go1.21.0`, environment variables are automatically loaded
- **Environment management**: Changes made with `sgv env -w` or `sgv env -u` are immediately applied to your current shell
- **Auto commands**: `sgv auto` and `sgv latest` automatically load environment variables after version switches
- **No manual steps**: No need to run `eval` commands or restart your terminal

This is achieved through a wrapper function that's automatically installed in your shell configuration.

---

## Environment Variables

- `SGV_DOWNLOAD_URL_PREFIX`  
  Change the Go download source (e.g., for China mainland users: `https://golang.google.cn/dl/`)

```sh
export SGV_DOWNLOAD_URL_PREFIX=https://golang.google.cn/dl/
```
- Set this before running sgv commands, or add to your shell profile for persistence

---

## Notes

- All Go versions are installed under `~/.sgv/versions/`, and the current version is symlinked as `~/.sgv/current`
- After switching, ensure your `GOROOT` and `PATH` are set (the install script handles this automatically)

---

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License. See [LICENSE](./LICENSE) for details.
