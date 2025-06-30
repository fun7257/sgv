# sgv - Go Version Manager

sgv (Simple Go Version) is a lightweight command-line tool for managing multiple Go versions on your system. It allows you to easily install, switch between, list, and uninstall different Go versions.

**Note:** sgv only supports Go versions 1.13 and later.

## Features

*   **Install Go Versions:** Download and install specific Go versions.
*   **Switch Go Versions:** Easily switch between installed Go versions.
*   **Auto-switch Go Versions:** Automatically switch to the most suitable Go version for the current Go project based on `go.mod`.
*   **Update to Latest Go:** Check for the latest Go version, install it if needed, and switch to it.
*   **List Local Versions:** View all Go versions installed by sgv.
*   **Uninstall Go Versions:** Remove installed Go versions.
*   **Display Version Information:** Show the build Go version and commit hash of sgv itself.

## Installation (Linux/macOS)

You can install `sgv` on Linux and macOS with a single command using the installation script. This script will automatically detect your operating system and architecture, download the latest pre-compiled binary, and set up the necessary environment variables.

```bash
curl -sSL https://raw.githubusercontent.com/fun7257/sgv/main/install.sh | bash
```

The script will install `sgv` to `/usr/local/bin` and will prompt for `sudo` access if required. It will also update your shell profile (`~/.bashrc` or `~/.zshrc`) to configure the `GOROOT` and `PATH` environment variables.

After the installation is complete, restart your terminal or run `source ~/.bashrc` (or `source ~/.zshrc`) to apply the changes.

## Installation (Windows)

To install `sgv` on Windows, open PowerShell **as an Administrator** and run the following command. This script will download the latest pre-compiled binary, install it, and configure the necessary system environment variables.

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor [System.Net.SecurityProtocolType]::Tls12; Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/fun7257/sgv/main/install.ps1'))
```

After the installation is complete, please restart your terminal or command prompt for the environment variable changes to take effect.

## Usage

### Install or Switch Go Version

To install a new Go version or switch to an already installed one:

```bash
sgv <version>
```

Examples:

*   `sgv 1.22.1` (installs or switches to Go 1.22.1)
*   `sgv go1.21.0` (installs or switches to Go 1.21.0)

**Note:** If the current directory is a Go project and the requested version is lower than the `go.mod` requirement, the operation will be aborted with an error.

### Install Go Version (without switching)

To download and install a specific Go version without making it the active one:

```bash
sgv install <version>
```

Example:

*   `sgv install 1.22.1`

### Auto-switch Go Version

To automatically switch to the most suitable Go version for the current Go project:

```bash
sgv auto
```

This command will:

1.  Check if the current directory contains a `go.mod` file.
2.  If a `go.mod` file is found, it will read the `go` and `toolchain` versions, prioritizing the `toolchain` version if it is higher. This determined version is the *exact* target version.
    *   **Version Parsing Note:** For Go 1.21 and later, if `go.mod` specifies a version like `go1.21`, `sgv` will interpret it as `go1.21.0`. For versions prior to Go 1.21 (e.g., `go1.20`), `go1.20` will continue to be interpreted as `go1.20` (without an appended `.0`).
3.  It will then check if this exact target version is installed locally.
4.  If the target version is not the currently active version, it will prompt you to switch.
5.  If you confirm, it will switch to the target version, downloading and installing it if it's not already present locally.
6.  If the current active Go version is already the target version, no switch will occur and no output will be displayed.
7.  If the current directory is not a Go project (no `go.mod` found), it will inform you and do nothing.

Example of the interactive prompt:
```
go.mod requires Go version: go1.22.1
Found suitable version: go1.22.1. (Will download and install)
Switch to this version? (y/n): y
```

### Update to the Latest Go Version

To check for the latest Go version, install it if it's not already installed, and switch to it:

```bash
sgv update
```

This command will:
1. Fetch the latest available Go version.
2. If the latest version is not already installed, it will be downloaded and installed.
3. Switch the active Go version to the latest version.

### List Installed Go Versions

To see all Go versions you have installed with `sgv`:

```bash
sgv list
```

This will show a list of installed versions, grouped by major version, with the currently active version marked.

Example output:
```
Installed Go versions:
go1.21:
  go1.21.0
  go1.21.9
go1.22:
  go1.22.1 <- current
  go1.22.4
```

### List Minor Versions of a Go Version

To see all minor versions of a specific Go version:

```bash
sgv sub <major_version>
```

Example:

*   `sgv sub 1.22`

This will show a list of all available Go 1.22.x versions, with currently installed versions marked.

**Note:** This command is only available for Go versions 1.13 and higher.

### Uninstall a Go Version

To remove an installed Go version:

```bash
sgv uninstall <version>
```

Example:

*   `sgv uninstall 1.22.1`

**Note:** You cannot uninstall the currently active Go version. You must switch to another version first.

### Display SGV Version

To display the Go version used to build `sgv` and its commit hash:

```bash
sgv version
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details. (You might want to create a LICENSE file if you haven't already.)
