# Contributing to Forge

Thank you for your interest in contributing to Forge! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Your environment (Go version, OS, etc.)
- Any relevant logs or error messages

### Suggesting Features

We welcome feature suggestions! Please create an issue with:
- A clear description of the feature
- Use cases and benefits
- Potential implementation approach (if you have ideas)

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Ensure all tests pass**: `make test`
5. **Run the linter**: `make lint`
6. **Format your code**: `make fmt`
7. **Commit your changes** with clear, descriptive commit messages
8. **Push to your fork** and submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Getting Started

```bash
# Clone the repository
git clone https://github.com/entrhq/forge.git
cd forge

# Install development tools
make install-tools

# Run tests
make test

# Run linter
make lint

# Build examples
make examples

# Run example application
make run-example
```

## Coding Standards

### Go Style

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (run `make fmt`)
- Write idiomatic Go code
- Keep functions focused and small
- Use meaningful variable and function names

### Code Organization

- **Public APIs** go in `pkg/` packages
- **Internal implementation** goes in `internal/` packages
- **Example applications** go in `examples/`
- **Documentation** goes in `docs/`

### Testing

- Write table-driven tests where appropriate
- Aim for good test coverage (check with `make test-coverage`)
- Use meaningful test names that describe what is being tested
- Mock external dependencies
- Test edge cases and error conditions

### Documentation

- Add package documentation for all public packages
- Document all exported types, functions, and methods
- Use examples in documentation where helpful
- Keep documentation up to date with code changes

### Commit Messages

Write clear commit messages:
- Use the imperative mood ("Add feature" not "Added feature")
- Keep the first line under 50 characters
- Provide detailed description in the body if needed
- Reference issues and PRs where relevant

Example:
```
Add LLM provider interface

Implement the core Provider interface for pluggable LLM support.
This enables integration with OpenAI-compatible APIs and allows
for custom provider implementations.

Fixes #123
```

## Project Structure

```
forge/
├── pkg/              # Public, importable packages
├── internal/         # Private implementation
├── examples/         # Example applications
├── docs/            # Documentation
├── .github/         # GitHub workflows and config
└── scripts/         # Build and utility scripts
```

## Testing Strategy

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test interaction between components
- **Example tests**: Ensure examples compile and run

Run all tests:
```bash
make test
```

Generate coverage report:
```bash
make test-coverage
```

## Makefile Targets

- `make test` - Run all tests with coverage
- `make lint` - Run linters
- `make fmt` - Format code
- `make examples` - Build example applications
- `make run-example` - Run simple example
- `make clean` - Clean build artifacts
- `make install-tools` - Install development tools
- `make all` - Run formatting, linting, tests, and build examples

## Questions?

If you have questions about contributing, feel free to:
- Open an issue for discussion
- Review existing issues and PRs
- Check the documentation in `docs/`

Thank you for contributing to Forge!