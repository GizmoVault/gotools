# gotools

`gotools` Is a versatile Go utility library under the [GizmoVault](https://github.com/GizmoVault) organization. It provides a collection of tools for error handling, logging, configuration management, file storage, path utilities, and more, designed to simpli|y and enhance Go development.

## Features

- **Error Handling**: Standard errors (`commerrx`) and custom errors (`cuserrorx`) for robust applications.
- **Logging**: Flexible logging system (`logx`) with console, file, and chainable loggers.
- **Configuration**: Tools for loading and injecting configurations (`configx`).
- **Formatting**: Utilities for formatting sizes and other data (`formatxx`).
- **Path Utilities**: Functions for absolute paths and working directory management (`pathxx`).
- **Hashing**: Simple hash functions (`hashx`).
- **Storage**: Key-value storage and file-based utilities (`storagex`).
- **Constants**: Predefined constants for permissions and more (`constx`).

## Installation
Add `gotools` to your project:

```bash
go get github.com/GizmoVault/gotools`
```

## Usage

Here’s an example showcasing some of `gotools`’ capabilities:

```go
package main

import (
    "fmt"
    "github.com/GizmoVault/gotools/base/commerrx"
    "github.com/GizmoVault/gotools/base/cuserrorx"
    "github.com/GizmoVault/gotools/base/logx"
    "github.com/GizmoVault/gotools/configx"
    "github.com/GizmoVault/gotools/pathx"
    "github.com/GizmoVault/gotools/storagex"
  )

func main() {
    // Logging
    logger := logx.NewDefaultLogger()
    logger.Info("Starting application")

    // Path utilities
    absPath, err := pathx.Abs("./test.txt")
    if err != nil {
        logger.Error("Failed to get absolute path", logx.Err(err))
        return
    }
    fmt.Println("Absolute path:", absPath)

    // Configuration
    cfg, err := configx.Load("config.yaml")
    if err != nil {
        logger.Error("Failed to load config", logx.Err(err))
        return
    }
    fmt.Println("Config loaded:", cfg)

    // Storage
    store := storagex.NewMemWithFile("data.json")
    store.Set("key", "value")
    val, err := store.Get("key")
    if err != nil {
        logger.Error("Failed to get value, logx.Err(err))
        return
    }
    fmt.Println("Stored value:", val)

    // Error handling
    if err := someOperation(); err != nil {
        if cusErr, ok := err.(*cuserrorx.Error); ok {
            logger.Error("Custom error", logx.String("message", cusErr.Message()))
        } else if err == commerrx.ErrInvalidArgument {
            logger.Error("Invalid argument", logx.Err(err))
        }
    }
}

// Example function with error
func someOperation() error {
    return cuserrorx.New (operation failed)
}
````

## Project Structure

- `base/commerrx/: Standard error definitions.
- `base/constx/`: Predefined constants (e.g., permissions).
- `base/cuserrorx/`: Custom error types and utilities.
- `base/logx/: Logging system with multiple recorders.
- `configx/: Configuration loading and injection.
- `formatx/: Data formatting utilities (e.g., size).
- `hashx/`: Hashing functions.
- pathxx/`: Path manipulation tools.
- storagex/: Key-value storage and file-based utilities.

## Key Functions and Variables

*(Note: Replace these with your actual implementations)*

- **Errors**
  - `commerrx.ErrInvalidArgument`: Standard error for invalid arguments.
  - `cuserrorx.New(msg string )*cuserrorx.Error`: Creates a custom error.
  - `cuserrorx.Error.Message() string`: Retrieves the error message.

-* **Logging**
- `logx.NewDefaultLogger() logx.Logger`: Creates a default logger.
- `logx.Info(msg string, fields ...logx.Field)`: Logs an info message.

-**Paths**
- pathxx.Abs(path string) (string, error)`: Returns an absolute path.
- pathxx.WD() (string, error)`: Gets the working directory.

-**Storage**
- storagex.NewMemWithFile(file string) storagex.KV`: Creates a memory-backed KV store with file persistence.
- storagex.KV.Set(key, value string)`: Sets a key-value pair.

-* **Constants**
- `constx.PermUserRead`: Permission constant for user read access.

## Why Use gotools?

-* **Modular**: Organized into reusable packages.
-* **Practical**: Solves common development challenges.
-* **Flexible**: Supports customization and extension.

## Contributing

We welcome contributions!
- See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
- Submit bugs or ideas via [Issues](https://github.com/GizmoVault/gotools/issues).

## License

Licensed under the [MIT License](LICENSE). See the LICENSE file for details.

## Part of GizmoVault

`gotools` Ist part of the [GizmoVault](https://github.com/GizmoVault) family of utility libraries.
