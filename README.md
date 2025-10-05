# FAAS - Function as a Service Framework

FAAS is a flexible function executor framework that provides a generic execution environment for common automation tasks. It implements a consistent interface pattern that makes it easy to add, configure, and execute various types of functions.

## Architecture

The framework follows a clean interface-based architecture where all functions implement the same `intf.Function` interface with standardized methods:

- `GetConfig()`: Returns function configuration and metadata
- `ParsePayload()`: Parses input data into structured format
- `Validate()`: Validates input parameters and requirements
- `Execute()`: Performs the actual function execution

## Available Functions

The framework includes built-in support for the following function types:

### **Email** (`email`)

Send emails through SendGrid API integration. Supports both plain text and HTML content with customizable sender/recipient information.

### **SMS** (`sms`)

Send SMS and MMS messages through Twilio API integration. Supports text messages and multimedia messaging with media URL attachments.

### **Slack** (`slack`)

Send messages to Slack channels using the Slack API. Perfect for notifications and team communication automation.

### **Docker Registry** (`docker_registry`)

Execute Docker containers from various registries including private registries with authentication support.

### **HTTP** (`http`)

Make HTTP requests (GET, POST, etc.) to external APIs and services. Handles request/response processing automatically.

### **Logger** (`logger`)

Log messages and data for debugging, monitoring, and audit purposes.

### **GitHub** (`github`)

GitHub integration for repository operations and automation (implementation in progress).

## Features

- ✅ **Consistent Interface**: All functions follow the same execution pattern
- ✅ **Type Safety**: Structured input validation and error handling
- ✅ **Extensible**: Easy to add new function types
- ✅ **Well Tested**: Comprehensive test suite with 50+ tests
- ✅ **Modern Go**: Uses latest Go patterns and best practices
- ✅ **Production Ready**: Error handling, logging, and monitoring support

## Development

### Prerequisites

- Go 1.24.1 or later
- Make (for build automation)

### Quick Start

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Build the project
make build

# Format and lint code
make check

# View all available commands
make help
```

### Adding New Functions

1. Create a new file in `faas/functions/`
2. Implement the `intf.Function` interface
3. Add constructor function (`NewYourFunction()`)
4. Register in `faas.go`
5. Add comprehensive tests

## Project Structure

```
faas/
├── faas/
│   ├── faas.go              # Main FAAS framework
│   ├── intf/
│   │   └── function.go      # Function interface definition
│   ├── functions/           # Function implementations
│   └── helpers/             # Utility packages
├── Makefile                 # Build automation
└── README.md               # This file
```
