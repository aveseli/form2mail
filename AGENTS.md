# Agent Guidelines for form2mail

This document provides guidelines for AI coding agents working on the form2mail codebase.

## Project Overview

form2mail is a Go-based HTTP service that handles contact form submissions and sends email notifications. It follows the standard Go project layout with `cmd/` for entry points and `internal/` for private application code.

## Build, Test, and Lint Commands

### Building
```bash
# Build the binary
go build -o form2mail cmd/server/main.go

# Build and run
go run cmd/server/main.go

# Build with Docker
docker build -t form2mail .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run a single test file
go test ./internal/handler/

# Run a specific test function
go test -run TestContactHandler ./internal/handler/

# Run tests with race detector
go test -race ./...
```

### Linting and Formatting
```bash
# Format all code (ALWAYS run before committing)
go fmt ./...

# Run go vet
go vet ./...

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

### Docker Commands
```bash
# Build locally
docker build -t form2mail .

# Run with environment variables
docker run -p 8080:8080 --env-file .env form2mail

# Run with docker-compose
docker-compose up -d
```

## Code Style Guidelines

### Project Structure
```
form2mail/
├── cmd/server/          # Application entry point (main.go only)
├── internal/            # Private application code (cannot be imported externally)
│   ├── config/          # Configuration loading
│   ├── email/           # Email sending functionality
│   └── handler/         # HTTP handlers
```

### Import Ordering
Imports must be organized in three groups, separated by blank lines:
1. Standard library packages
2. External packages
3. Internal packages (from this module)

**Example:**
```go
import (
    "fmt"
    "log"
    "net/http"

    "github.com/external/package"

    "form2mail/internal/config"
    "form2mail/internal/email"
)
```

### Naming Conventions

**Packages:**
- Use short, lowercase, single-word names
- No underscores or camelCase
- Examples: `config`, `email`, `handler`

**Types:**
- Use PascalCase for exported types
- Use camelCase for unexported types
- Examples: `Config`, `ContactHandler`, `emailSender` (if unexported)

**Functions:**
- Use PascalCase for exported functions
- Use camelCase for unexported functions
- Constructor functions should be named `New<Type>` or `New`
- Examples: `NewSender()`, `Load()`, `getEnv()`

**Variables:**
- Use camelCase for local variables
- Use short names for short scopes (e.g., `cfg` for config)
- Use descriptive names for larger scopes
- Avoid single-letter names except for:
  - Loop indexes: `i`, `j`, `k`
  - HTTP handlers: `w` (http.ResponseWriter), `r` (*http.Request)
  - Context: `ctx`

**Constants:**
- Use PascalCase for exported constants
- Use camelCase or UPPER_SNAKE_CASE for unexported constants

### Type Definitions

**Structs:**
- Define struct types before functions
- Use pointer receivers for methods that modify state
- Use value receivers for methods that don't modify state
- Add struct tags for JSON/form parsing

**Example:**
```go
type ContactForm struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Subject string `json:"subject"`
    Message string `json:"message"`
}
```

### Error Handling

**Always handle errors explicitly. Never ignore errors.**

```go
// GOOD: Check and handle errors
if err := someFunction(); err != nil {
    log.Printf("Failed to do something: %v", err)
    return err
}

// BAD: Ignoring errors
_ = someFunction()

// GOOD: Return errors to caller
func process() error {
    if err := validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    return nil
}

// Use log.Printf for non-fatal errors
log.Printf("Failed to send confirmation email: %v", err)

// Use log.Fatal for fatal errors (only in main)
log.Fatal("SMTP_USER must be set")
```

### HTTP Handlers

**Use the http.Handler interface for handlers:**
```go
type ContactHandler struct {
    emailSender *email.Sender
    corsOrigin  string
}

func (h *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

**Check HTTP methods explicitly:**
```go
if r.Method != http.MethodPost {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}
```

### Comments

- Add package-level comments for all packages
- Add comments for all exported types, functions, and constants
- Use complete sentences starting with the item name
- Avoid obvious comments

**Example:**
```go
// Config holds application configuration loaded from environment variables.
type Config struct {
    SMTPHost string
}

// Load reads configuration from environment variables and returns a Config.
func Load() Config {
    // implementation
}
```

### Configuration

- All configuration comes from environment variables
- Use sensible defaults where appropriate
- Document required vs optional variables
- Validate required configuration at startup in `main()`

### Dependency Injection

- Use constructor functions (New*) to create instances
- Pass dependencies explicitly via constructor
- Avoid global state

**Example:**
```go
func NewContactHandler(emailSender *email.Sender, corsOrigin string) *ContactHandler {
    return &ContactHandler{
        emailSender: emailSender,
        corsOrigin:  corsOrigin,
    }
}
```

## Testing Guidelines

When creating tests:
- Place test files next to the code they test
- Name test files with `_test.go` suffix
- Use table-driven tests for multiple cases
- Mock external dependencies (SMTP, HTTP)
- Test error cases, not just happy paths

**Example test structure:**
```go
func TestContactHandler_ServeHTTP(t *testing.T) {
    tests := []struct {
        name       string
        method     string
        body       string
        wantStatus int
    }{
        // test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Git Workflow

- Commit messages should be clear and descriptive
- Reference issue numbers if applicable
- Run `go fmt ./...` before committing
- Run `go mod tidy` if dependencies changed

## Docker Notes

- Multi-stage builds to keep image size small
- Use Alpine Linux for minimal base image
- Expose port 8080 by default
- Configuration via environment variables only (no .env in container)
