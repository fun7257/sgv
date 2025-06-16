# sgv - Go Version Manager

sgv (Simple Go Version) is a lightweight command-line tool for managing multiple Go versions on your system. It allows you to easily install, switch between, list, and uninstall different Go versions.

## Features

*   **Install Go Versions:** Download and install specific Go versions.
*   **Switch Go Versions:** Easily switch between installed Go versions.
*   **List Local Versions:** View all Go versions installed by sgv.
*   **List Remote Versions:** See available Go versions from the official Go website.
*   **Uninstall Go Versions:** Remove installed Go versions.

## Installation

To install `sgv`, make sure you have Go installed (any version will do for building `sgv` itself).

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/sgv.git # Replace with your actual repo URL
    cd sgv
    ```

2.  **Build and install:**
    ```bash
    make install
    ```
    This will compile `sgv` and place the executable in your `$GOPATH/bin` (or `$GOBIN`) directory, making it available in your PATH.

## Usage

### Install or Switch Go Version

To install a new Go version or switch to an already installed one:

```bash
sgv <version>
```

Examples:

*   `sgv 1.22.1` (installs or switches to Go 1.22.1)
*   `sgv go1.21.0` (installs or switches to Go 1.21.0)

### List Installed Go Versions

To see all Go versions you have installed with `sgv`:

```bash
sgv list
```

This will show a list of installed versions, with the currently active version marked.

### List Available Remote Go Versions

To see all Go versions available for download from the official Go website:

```bash
sgv list-remote
```

### Uninstall a Go Version

To remove an installed Go version:

```bash
sgv uninstall <version>
```

Example:

*   `sgv uninstall 1.22.1`

**Note:** You cannot uninstall the currently active Go version. You must switch to another version first.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details. (You might want to create a LICENSE file if you haven't already.)
