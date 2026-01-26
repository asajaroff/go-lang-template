# Contributing Guide

Thank you for contributing to this project! This guide will help you understand our development workflow and conventions.

## Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/asajaroff/go-lang-template.git
   cd go-lang-template
   ```

2. **Install pre-commit hooks**
   ```bash
   make setup-hooks
   ```

   This will install git hooks that automatically:
   - Validate commit message format
   - Run code formatting
   - Run static analysis
   - Run tests before committing

3. **Install development tools**
   ```bash
   make install-tools
   ```

## Commit Message Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/) specification. All commit messages must follow this format:

```
<type>: <description>

[optional body]

[optional footer]
```

### Commit Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (white-space, formatting, etc)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvements
- **test**: Adding or correcting tests
- **chore**: Changes to build process or auxiliary tools
- **ci**: Changes to CI configuration files and scripts
- **build**: Changes that affect the build system or external dependencies

### Examples

Good commit messages:
```
feat: add user authentication endpoint
fix: resolve nil pointer dereference in config loader
docs: update README with installation instructions
test: add unit tests for health check handler
refactor: simplify error handling in database layer
chore: update dependencies to latest versions
```

Bad commit messages:
```
fixed stuff
WIP
updated code
asdf
```

### Breaking Changes

For breaking changes, add `!` after the type or add `BREAKING CHANGE:` in the footer:

```
feat!: change API response format to include metadata

BREAKING CHANGE: API responses now include a metadata field
```

## Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/) (SemVer):

- **MAJOR** version (X.0.0): Incompatible API changes (breaking changes)
- **MINOR** version (0.X.0): Add functionality in a backwards-compatible manner
- **PATCH** version (0.0.X): Backwards-compatible bug fixes

Commit types map to version bumps:
- `feat:` → MINOR version bump
- `fix:` → PATCH version bump
- `feat!:` or `BREAKING CHANGE:` → MAJOR version bump
- Other types → No version bump (but included in changelog)

## Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feat/my-new-feature
   ```

2. **Make your changes**
   - Write clean, well-documented code
   - Follow Go best practices
   - Add tests for new functionality

3. **Run checks locally**
   ```bash
   make test    # Run tests
   make lint    # Run linting
   make build   # Build the project
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

   The pre-commit hooks will automatically:
   - Validate commit message format
   - Format code with `go fmt`
   - Run `go vet` for static analysis
   - Run all tests

5. **Push and create a pull request**
   ```bash
   git push origin feat/my-new-feature
   ```

## Testing

- Write tests for all new functionality
- Maintain test coverage above 80%
- Run tests with: `make test`
- Check coverage: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

## Code Style

- Use `go fmt` for formatting (automatic with pre-commit hooks)
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Keep functions small and focused
- Write clear, descriptive variable names
- Add comments for complex logic

## Pre-commit Hook Bypass

In rare cases where you need to bypass pre-commit hooks (not recommended):

```bash
git commit --no-verify -m "your message"
```

**Note**: This should only be used in exceptional circumstances and will be reviewed carefully in pull requests.

## Questions?

If you have questions or need help, please open an issue or reach out to the maintainers.
