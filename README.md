[English](README.md) | [简体中文](README_zh.md)
# sgv - Go Version Manager

`sgv` (Simple Go Version) is a lightweight CLI tool for managing multiple Go versions on your system. It supports Go 1.13 and above, and works on macOS, Linux, and Windows.

## Features

- **Install Go versions**: Download and install any supported Go version.
- **Switch Go versions**: Instantly switch between installed Go versions.
- **Auto switch**: Automatically switch to the required Go version for the current project based on `go.mod`.
- **Get latest**: Install and switch to the latest Go version with one command.
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

### Windows

Open PowerShell as Administrator and run:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor [System.Net.SecurityProtocolType]::Tls12; Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/fun7257/sgv/main/install.ps1'))
```

- Installs to `C:\Program Files\sgv\sgv.exe`
- Automatically configures environment variables
- After installation, restart your terminal

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
