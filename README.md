# Simonsays

A powerful, extensible Swiss Army knife for file and directory manipulation tasks.

[![Go Report Card](https://goreportcard.com/badge/github.com/GyroZeplix/simonsays)](https://goreportcard.com/report/github.com/GyroZepelix/simonsays)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

Simonsays is a command-line tool written in Go that provides various file management utilities in a modular, extensible framework. It simplifies common file organization tasks like sorting files by date or type, bulk renaming, and more.

## Features

- **Sort by Date**: Organize files into directories based on creation date
- **Sort by Type**: Categorize files into directories based on file extensions
- **Sort by Size**: Group files into small, medium, and large categories
- **Bulk Rename**: Rename multiple files using customizable patterns
- **List Files**: Display files with various sorting options
- **Extensible Framework**: Easily add new commands with minimal boilerplate

## Installation

### Prerequisites

- Go 1.24 or higher
- [Just](https://github.com/casey/just) (optional, for build automation)

### From Source

#### Using Just (recommended)

```bash
# Clone the repository
git clone https://github.com/GyroZepelix/simonsays.git
cd simonsays

# Build the binary
just build

# Optionally, install to GOPATH/bin
just install
```

#### Using Go commands

```bash
git clone https://github.com/GyroZepelix/simonsays.git
cd simonsays
go build -o bin/simonsays cmd/simonsays/main.go
```

### Using Go Install

```bash
go install github.com/GyroZepelix/simonsays@latest
```

## Build Tasks

Run to see all available build tasks

```bash
just -l                # List all available build tasks
```

## Usage

### Basic Syntax

```bash
simonsays [command] [arguments] [options]
```

### Available Commands

#### Sort by Date

Organizes files into directories based on creation date (format: DDMMYYYY).

```bash
simonsays sortbydate /path/to/directory
```

**Example:**
```bash
$ simonsays sortbydate ~/Downloads
Moved /home/user/Downloads/document.pdf to /home/user/Downloads/08032025/document.pdf
Moved /home/user/Downloads/image.jpg to /home/user/Downloads/07032025/image.jpg
```

#### Bulk Rename

Renames multiple files using a pattern.

```bash
simonsays bulkrename /path/to/directory "pattern" [options]
```

**Pattern Variables:**
- `{index}`: Sequential number (starting from 1 or specified by --start)
- `{name}`: Original filename without extension
- `{ext}`: Original file extension

**Options:**
- `--dry-run`: Show what would be renamed without actually renaming files
- `--start [num]`: Starting index for the {index} variable (default: 1)

**Example:**
```bash
$ simonsays bulkrename ~/Photos "vacation_{index}" --start 100
Renamed: /home/user/Photos/IMG_1234.jpg -> /home/user/Photos/vacation_100.jpg
Renamed: /home/user/Photos/IMG_1235.jpg -> /home/user/Photos/vacation_101.jpg
```

#### Sort by Type

Organizes files into directories based on their file extensions.

```bash
simonsays sortbytype /path/to/directory
```

**Example:**
```bash
$ simonsays sortbytype ~/Documents
Moved /home/user/Documents/report.pdf to /home/user/Documents/pdf/report.pdf
Moved /home/user/Documents/data.csv to /home/user/Documents/csv/data.csv
```

#### Sort by Size

Categorizes files into small (<1MB), medium (1-10MB), and large (>10MB) directories.

```bash
simonsays sortbysize /path/to/directory
```

**Example:**
```bash
$ simonsays sortbysize ~/Downloads
Moved /home/user/Downloads/small-image.jpg to /home/user/Downloads/small/small-image.jpg
Moved /home/user/Downloads/presentation.pptx to /home/user/Downloads/medium/presentation.pptx
Moved /home/user/Downloads/video.mp4 to /home/user/Downloads/large/video.mp4
```

#### List Files

Lists files with various sorting options.

```bash
simonsays list /path/to/directory [options]
```

**Options:**
- `--sort [name|size|time]`: Sort files by name, size, or modification time (default: name)
- `--recursive`: List files recursively, including those in subdirectories

**Example:**
```bash
$ simonsays list ~/Documents --sort size --recursive
Files:
- large-document.pdf (Size: 5242880 bytes, Modified: 2025-03-07T15:42:12Z)
- presentation.pptx (Size: 2097152 bytes, Modified: 2025-03-06T09:23:45Z)
- notes.txt (Size: 1024 bytes, Modified: 2025-03-08T10:15:30Z)
```

## Extending Simonsays

Simonsays is designed to be easily extensible with new commands. Here's how to add a new command:

1. Create a new struct that implements the `Command` interface:
```go
type MyNewCommand struct{}

func (c *MyNewCommand) Name() string {
    return "mynewcommand"
}

func (c *MyNewCommand) Description() string {
    return "Description of my new command"
}

func (c *MyNewCommand) Execute(ctx *cli.Context) error {
    // Implementation goes here
    return nil
}
```

2. Register the command in the `main` function:
```go
registry.Register(&MyNewCommand{})
```

3. Add command-specific flags if needed:
```go
switch command.Name() {
case "mynewcommand":
    cliCommand.Flags = []cli.Flag{
        &cli.StringFlag{
            Name:  "flag-name",
            Usage: "Description of the flag",
            Value: "default-value",
        },
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) - CLI framework for Go
- [just](https://github.com/casey/just) - Command runner used for build automation
