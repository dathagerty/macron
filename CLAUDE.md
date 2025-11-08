# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Macron is a Go-based CLI tool for managing launchd cron tasks on macOS. It provides a user-friendly interface to create, edit, enable, disable, and delete launchd agents for scheduled task execution.

## Architecture

### Core Structure

- **main.go**: Entry point that calls `cmd.Execute()`
- **cmd/**: Contains all Cobra command definitions
  - **root.go**: Root command setup with Viper configuration (reads from `$HOME/.macron.yaml`)
  - **create.go**: Creates new launchd cron tasks with NAME, SCRIPT, and INTERVAL
  - **edit.go**: Edits existing launchd tasks
  - **enable.go**: Enables launchd tasks
  - **disable.go**: Disables launchd tasks
  - **delete.go**: Deletes launchd tasks
- **internal/version/**: Version information extracted from build metadata using `debug.ReadBuildInfo()`

### Dependencies

Built with:
- **Cobra**: CLI framework for command structure
- **Viper**: Configuration management (supports config file and environment variables)
- **mitchellh/go-homedir**: Cross-platform home directory detection

### Configuration

The application uses Viper to load configuration from:
1. Config file specified via `--config` flag
2. Default location: `$HOME/.macron.yaml`
3. Environment variables (automatically mapped)

## Development Commands

### Build
```bash
go build -o macron
```

### Run
```bash
go run main.go [command]
```

### Test
```bash
go test ./...
```

### Install Dependencies
```bash
go mod download
```

### Update Dependencies
```bash
go mod tidy
```

## Implementation Notes

### Adding New Commands

New commands should:
1. Be added as separate files in the `cmd/` directory
2. Follow the Cobra command pattern with `Use`, `Short`, and `Long` descriptions
3. Register themselves in `init()` via `rootCmd.AddCommand(yourCmd)`
4. Implement functionality in the `Run` function

### Version Information

Version information is automatically extracted at build time from VCS data:
- `Revision`: Git commit hash
- `Version`: Module version from build info
- `DirtyTree`: Whether the working tree had uncommitted changes
- `LastCommit`: Timestamp of the last commit

The version is exposed via the root command's `Version` field.

### launchd Integration

This tool is designed to work with macOS launchd agents. When implementing the actual functionality, you'll need to:
- Generate `.plist` files for launchd agents
- Place them in `~/Library/LaunchAgents/`
- Use `launchctl load/unload` to enable/disable tasks
- Use `launchctl start/stop` for manual execution
