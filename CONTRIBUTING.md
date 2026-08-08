# Contributing to Mediax

Thank you for your interest in contributing to Mediax! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, FFmpeg version)
- Screenshots or logs if applicable

### Suggesting Enhancements

Enhancement suggestions are welcome! Please:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful
- Provide examples of how the enhancement would be used

### Pull Requests

1. Fork the repository
2. Create a branch for your feature (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Coding Standards

- Follow Go best practices and conventions
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused
- Write tests for new features
- Ensure code passes `go fmt` and `go vet`

### Development Setup

```bash
# Clone your fork
git clone https://github.com/robert-sarah/mediax.git
cd mediax

# Install dependencies
go mod download

# Run tests
go test ./...

# Build the project
go build -o mediax
```

## Project Structure

```
mediax/
├── cmd/           # CLI command definitions
├── internal/
│   ├── tui/       # Terminal user interface
│   └── verbs/     # Verb implementations
├── assets/        # Logo and branding
├── main.go        # Entry point
└── go.mod         # Go module file
```

## Adding New Verbs

To add a new verb:

1. Create a new file in `internal/verbs/` (e.g., `newverb.go`)
2. Implement the `Verb` interface:
   ```go
   type Verb interface {
       Name() string
       Description() string
       Usage() string
       Execute(args []string, flags map[string]string) error
   }
   ```
3. Register the verb in the file's `init()` function
4. Add the verb to the appropriate category in `internal/tui/mainmenu.go`
5. Update documentation in `cmd/root.go` and `README.md`

## Testing

Run tests with:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Documentation

- Keep README.md up to date
- Add godoc comments to exported functions
- Update CHANGELOG.md for significant changes

## Questions?

Feel free to open an issue for any questions about contributing to Mediax.
