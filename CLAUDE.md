# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `mono-track`, a Go-based CLI application for tracking and analyzing expenses using AI. The project is currently in early development with basic command structure in place.

## Commands

### Build and Run
- `go run src/main.go` - Run the application directly
- `make build` - Build binary to bin/mono-track using Makefile
- `go build -o bin/mono-track ./src` - Build binary directly
- `./bin/mono-track` - Run the built binary

### Testing
- `make test-unit` - Run unit tests only
- `make coverage-unit` - Run unit tests with coverage
- `make test-e2e` - Run end-to-end tests (builds coverage binary first)
- `make coverage-e2e` - Run e2e tests with coverage reporting
- `make build-e2e-coverage` - Build test binary with coverage instrumentation
- `make clean` - Clean build files and coverage data

### Application Commands
- `mono-track` - Run main expense tracking flow (currently placeholder)
- `mono-track help` - Show help information
- `mono-track version` - Show version information

## Architecture

### Core Structure
- **src/main.go**: Entry point with basic CLI command routing (help, version, main flow)
- **Module**: `mono-track` (Go 1.25.0)
- **Version**: Currently hardcoded as "1.0.0" in main.go

### Test Framework
The project has a sophisticated end-to-end testing framework in `tests/e2e/testutil/`:

- **FlowTestBuilder**: Fluent API for building flow execution tests with configurable timeouts, environment variables, and expectations
- **FlowRunner**: Handles subprocess execution with coverage collection and security validations
- **Test Organization**: Uses coverage tracking with `GOCOVERDIR` environment variable
- **Binary Management**: Automatically builds test binaries when needed (`bin/mono-track-e2e`)

### Key Test Patterns
- Tests expect a test binary at `bin/mono-track-e2e`
- Coverage data collected in `coverage/e2e/` directory
- Uses builder pattern for test configuration: `NewFlowTest(t).WithCommand("list").ExpectSuccess().Run()`
- Supports both flow file execution and direct command testing

### Claude Configuration
The project includes extensive Claude Code configuration in `.claude/`:
- **Agents**: Specialized agents for PR commits, user story creation, PRD interviews, and user story refinement
- **Templates**: Product requirements and user story templates
- **Metadata**: PRD documents for both CLI and webapp components

## Development Notes

### Current State
- Basic CLI structure with command routing implemented
- Sophisticated testing infrastructure ready for development
- Main expense tracking functionality is placeholder (to be implemented)

### Testing Approach
- E2E tests use subprocess execution of built binaries
- Coverage collection integrated into test runs
- Security validations for binary paths and command arguments
- Timeout and environment variable management built-in
