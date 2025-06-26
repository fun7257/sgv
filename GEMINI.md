# Gemini Project Brief: sgv

## Project Overview

`sgv` (Simple Go Version) is a command-line tool written in Go for managing multiple Go language versions. It allows users to install, switch between, list, and uninstall different Go versions. The tool can also automatically select a suitable Go version based on the `go.mod` file in the current project.

The project uses the `cobra` library to structure its command-line interface.

## Key Files and Directories

-   `main.go`: The entry point of the application. It initializes and executes the root command.
-   `go.mod` / `go.sum`: Defines the project's module path and manages its dependencies.
-   `README.md`: Contains detailed user-facing documentation, including installation and usage instructions.
-   `Makefile`: Provides convenience targets for building and installing the application (`make install`).
-   `cmd/`: This directory contains all the command definitions for the CLI.
    -   `root.go`: Sets up the main `sgv` root command.
    -   `install.go`, `list.go`, `update.go`, etc.: Each file defines a specific subcommand (e.g., `sgv install`, `sgv list`).
-   `internal/`: Contains the core application logic, separated from the command-line interface.
    -   `installer/installer.go`: Handles the logic for downloading, installing, and managing Go versions.
    -   `config/config.go`: Manages configuration settings.
    -   `version/version.go`: Handles version-related operations and information.

## Core Commands

To interact with the project during development:

-   **Run commands:** Use `go run main.go <subcommand>` to execute specific commands.
    -   Example: `go run main.go list`
    -   Example: `go run main.go auto`
-   **Build and install:** Use `make install` to compile the binary and install it to `$GOPATH/bin`.
-   **Run tests:** Use `go test ./...` to run all tests within the project.

## Development Workflow

1.  **Function and Documentation Synchronization**: When modifying or adding any functionality, the `README.md` file must be updated accordingly to ensure that the documentation remains consistent with the code's functionality.
2.  **Update Project Brief**: If a change significantly alters the project's scope, architecture, or core dependencies, the 'Project Overview' section in `GEMINI.md` should be updated.
