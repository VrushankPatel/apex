# Contributing to APEX

We love your input! We want to make contributing to APEX as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `master`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Issue that pull request!

## Pull Request Process

1. Update the README.md with details of changes to the interface, if applicable
2. Update the docs/ with any new documentation or changes to existing docs
3. The PR will be merged once you have the sign-off of at least one other developer
4. If you haven't already, complete the Contributor License Agreement ("CLA")

## Any contributions you make will be under the MIT Software License

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using GitHub's [issue tracker](https://github.com/VrushankPatel/apex/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/VrushankPatel/apex/issues/new).

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Code Style Guidelines

### Go

Follow the standard Go style guidelines:

- Use `gofmt` to format your code
- Follow [Effective Go](https://golang.org/doc/effective_go.html) principles
- Document your code using [GoDoc](https://blog.golang.org/godoc-documenting-go-code) conventions

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./pkg/detector
```

### Writing Tests

- Write test cases for new features
- Update tests when modifying existing features
- Aim for high test coverage
- Use table-driven tests when appropriate

## Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Include examples in documentation when helpful
- Keep documentation up to date with code changes

### Project Documentation

- Update README.md for significant changes
- Maintain API documentation
- Add/update guides in docs/ directory

## Development Setup

1. Install Go 1.18 or higher
2. Clone your fork of the repo
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Set up pre-commit hooks:
   ```bash
   cp hooks/pre-commit .git/hooks/
   chmod +x .git/hooks/pre-commit
   ```

## Project Structure

```
apex/
├── cmd/                    # Command line tools
├── pkg/                    # Library code
│   ├── config/            # Configuration handling
│   ├── detector/          # Arbitrage detection logic
│   ├── exchange/          # Exchange integrations
│   ├── models/            # Data models
│   ├── server/            # Web server
│   └── util/              # Utilities
├── web/                   # Web interface
├── docs/                  # Documentation
└── tests/                 # Integration tests
```

## Adding New Features

1. **Exchange Integration**
   - Implement the Exchange interface in `pkg/exchange`
   - Add configuration in `pkg/config`
   - Add tests in `pkg/exchange/tests`

2. **Arbitrage Strategies**
   - Add new strategy in `pkg/detector`
   - Implement required models in `pkg/models`
   - Add configuration options if needed

3. **Web Interface**
   - Add new components in `web/components`
   - Update API endpoints if required
   - Add new WebSocket message types if needed

## Review Process

1. **Code Review**
   - All code changes require review
   - Address review comments promptly
   - Keep discussions focused and professional

2. **Testing Review**
   - Ensure all tests pass
   - Check test coverage
   - Verify integration tests

3. **Documentation Review**
   - Check documentation accuracy
   - Verify API documentation
   - Review guide updates

## Community

- Join our [Discord](https://discord.gg/apex-arbitrage) for discussions
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)
- Participate in issue discussions
- Help others in the community

## Recognition

We believe in recognizing contributions:

- Contributors are listed in CONTRIBUTORS.md
- Significant contributions may lead to maintainer status
- We highlight important contributions in release notes

## Questions?

Don't hesitate to ask questions:

- Open a [GitHub Discussion](https://github.com/VrushankPatel/apex/discussions)
- Join our Discord community
- Contact the maintainers directly

Thank you for contributing to APEX! 