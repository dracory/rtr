# Contributing to RTR Router

Thank you for considering contributing to the RTR Router project! We welcome all forms of contributions, including bug reports, feature requests, documentation improvements, and code contributions.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)
- [Feature Requests](#feature-requests)
- [Documentation](#documentation)
- [License](#license)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally
   ```bash
   git clone https://github.com/dracory/rtr.git
   cd rtr
   ```
3. Set up the development environment:
   ```bash
   go mod download
   ```
4. Run the tests to verify your setup:
   ```bash
   go test ./...
   ```

## Development Workflow

1. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b bugfix/issue-number-description
   # or
   git checkout -b docs/your-documentation-change
   # or
   git checkout -b change/your-change-description
   ```

2. Make your changes following the code style guidelines

3. Run the tests:
   ```bash
   go test ./...
   ```

4. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```

5. Push your changes to your fork:
   ```bash
   git push origin your-branch-name
   ```

6. Open a pull request against the main branch

## Code Style

We follow the standard Go code style. Please run the following before committing:

```bash
gofmt -s -w .
goimports -w .
```

### Guidelines

- Use descriptive variable and function names
- Keep functions small and focused
- Document all exported functions and types
- Write tests for new functionality
- Update documentation when making changes

## Testing

We aim for high test coverage. Please ensure:

1. New features include appropriate tests
2. Bug fixes include regression tests
3. Tests pass on all supported Go versions

Run tests with:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Run benchmarks
go test -bench=.
```

## Pull Request Process

1. Ensure your fork is up to date with the main branch
2. Rebase your changes on top of the latest main branch
3. Make sure all tests pass
4. Update documentation as needed
5. Submit the pull request with a clear description of the changes

### PR Guidelines

- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update relevant documentation
- Reference any related issues

## Reporting Issues

When reporting issues, please include:

1. The version of Go you're using
2. The version of the router
3. Steps to reproduce the issue
4. Expected behavior
5. Actual behavior
6. Any relevant error messages or logs

## Feature Requests

We welcome feature requests! Please:

1. Check if the feature has already been requested
2. Clearly describe the feature and its benefits
3. Include any relevant use cases

## Documentation

Good documentation is crucial. When making changes:

1. Update relevant documentation
2. Add examples for new features
3. Keep the README up to date

## License

By contributing, you agree that your contributions will be licensed
under the project's [LICENSE](LICENSE) file.
