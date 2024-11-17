# sysinfo

## Overview

`sysinfo` is a command-line tool written in Go that fetches and displays system information. It is designed to provide users with a quick overview of various system metrics and hardware information.

## Features

- Retrieves CPU information.
- Displays memory usage statistics.
- Provides details about the operating system and kernel version.
- Shows disk usage and partition information.
- Network configuration details.

## Installation

### Prerequisites

- Go 1.16 or higher installed on your machine.

### Steps

1. Clone the repository:
   ```sh
   git clone <repository-url>
   ```

2. Navigate to the project directory:
   ```sh
   cd sysinfo
   ```

3. Install dependencies:
   ```sh
   go mod tidy
   ```

4. Build the binary:
   ```sh
   go build -o sysinfo cmd/systeminfo/main.go
   ```

## Usage

Run the `sysinfo` binary to display system information:

```sh
./sysinfo
```

You can use various flags to customize the output:

```sh
./sysinfo --help
```

### Available Flags

```sh
Usage of ./sysinfo:
  -filter value
    comma-separated list of metrics, available: cpu,mem,net,disk,osinf
  -format string
    output format (text, json) (default "text")
  -log-level string
    log level (debug, info, warn, error, fatal, panic) (default "info")
```

### Example Usage

- Display only CPU and memory information:
  ```sh
  ./sysinfo -filter cpu,mem
  ```

- Display output in JSON format:
  ```sh
  ./sysinfo -format json
  ```

- Set the log level to debug:
  ```sh
  ./sysinfo -log-level debug
  ```

## Project Structure

- `cmd/`: Contains the main entry point for the application.
- `internal/`: Internal packages for retrieving system information.
- `pkg/`: Utility packages used by the application.
- `go.mod` and `go.sum`: Go modules files for dependency management.