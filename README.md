# sgv - Go Version Manager

sgv (Simple Go Version) is a lightweight command-line tool for managing multiple Go versions on your system. It allows you to easily install, switch between, list, and uninstall different Go versions.

**Note:** sgv only supports Go versions 1.13 and later.

## Features

*   **Install Go Versions:** Download and install specific Go versions.
*   **Switch Go Versions:** Easily switch between installed Go versions.
*   **Auto-switch Go Versions:** Automatically switch to the most suitable Go version for the current Go project based on `go.mod`.
*   **List Local Versions:** View all Go versions installed by sgv.
*   **Uninstall Go Versions:** Remove installed Go versions.
*   **Display Version Information:** Show the build Go version and commit hash of sgv itself.

## Installation

To install `sgv`, make sure you have Go installed (any version will do for building `sgv` itself).

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/fun7257/sgv.git
    cd sgv
    ```

2.  **Build and install:**
    ```bash
    make install
    ```
    This will compile `sgv` and place the executable in your `$GOPATH/bin` (or `$GOBIN`) directory, making it available in your PATH.

### Environment Variables Setup

After installing `sgv` and your first Go version, you need to set up some environment variables to make `go` commands work correctly. Add the following lines to your shell profile (e.g., `~/.zshrc`, `~/.bashrc`, `~/.profile`):

```bash
export SGV_ROOT="$HOME/.sgv"
export GOROOT="$SGV_ROOT/current"
export PATH="$GOROOT/bin:$HOME/go/bin:$PATH"
unset GOPATH
```

**For `~/.zshrc` users, you can append these lines directly by running:**

```bash
cat << 'EOF' >> ~/.zshrc
export SGV_ROOT="$HOME/.sgv"
export GOROOT="$SGV_ROOT/current"
export PATH="$GOROOT/bin:$HOME/go/bin:$PATH"
unset GOPATH
EOF
```

After adding these lines, remember to `source` your shell profile to apply the changes:

```bash
source ~/.zshrc # Or ~/.bashrc, ~/.profile, etc.
```

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
2.  If a `go.mod` file is found, it will read the required Go version.
3.  It will then find the smallest installed Go version that meets or exceeds the `go.mod` requirement.
4.  If a suitable version is found and it's not already the active version, `sgv` will switch to it.
5.  If no suitable installed version is found, it will prompt you to install a compatible version.
6.  If the current directory is not a Go project (no `go.mod` found), it will inform you and do nothing.

### List Installed Go Versions

To see all Go versions you have installed with `sgv`:

```bash
sgv list
```

This will show a list of installed versions, with the currently active version marked.

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
