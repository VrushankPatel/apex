# Contributing to APEX

This guide explains how to contribute to APEX through GitHub.

## Getting Started

### 1. Fork the Repository

1. Visit the [APEX repository](https://github.com/VrushankPatel/apex) on GitHub
2. Click the "Fork" button in the top-right corner
3. Select your GitHub account as the destination for the fork

### 2. Clone Your Fork

```bash
git clone https://github.com/YOUR_USERNAME/apex.git
cd apex
```

### 3. Set Up Remote Repositories

```bash
# Add the original repository as "upstream"
git remote add upstream https://github.com/VrushankPatel/apex.git

# Verify your remotes
git remote -v
```

## Making Contributions

### 1. Create a Feature Branch

Always create a new branch for your changes:

```bash
# Ensure you're up to date
git fetch upstream
git checkout main
git merge upstream/main

# Create a new branch
git checkout -b feature/your-feature-name
```

### 2. Development Guidelines

1. **Code Style**
   - Follow Go best practices and conventions
   - Use meaningful variable and function names
   - Add comments for complex logic
   - Write tests for new features

2. **Commit Messages**
   - Write clear, descriptive commit messages
   - Start with a verb (e.g., "Add", "Fix", "Update")
   - Reference issue numbers if applicable

### 3. Testing Your Changes

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run
```

### 4. Submitting Changes

1. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request**
   - Go to the [APEX repository](https://github.com/VrushankPatel/apex)
   - Click "Pull Requests" > "New Pull Request"
   - Select "compare across forks"
   - Choose your fork and branch
   - Fill out the PR template
   - Submit the pull request

### 5. PR Best Practices

- Keep changes focused and atomic
- Include tests for new features
- Update documentation if needed
- Respond to review comments promptly
- Rebase on main if requested

## Keeping Your Fork Updated

Regularly sync your fork with the main repository:

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

## Issue Guidelines

### Creating Issues

1. **Search First**
   - Check if the issue already exists
   - Look through closed issues

2. **Issue Types**
   - Bug reports
   - Feature requests
   - Documentation improvements
   - Performance issues

3. **Issue Template**
   - Use the provided issue templates
   - Include all requested information
   - Add steps to reproduce for bugs

## Code Review Process

1. **Initial Review**
   - Maintainers will review your PR
   - Automated tests must pass
   - Code style must meet guidelines

2. **Feedback**
   - Address all review comments
   - Ask questions if unclear
   - Make requested changes

3. **Approval and Merge**
   - PR needs maintainer approval
   - Changes may be requested
   - Maintainers will merge approved PRs

## Community Guidelines

- Be respectful and constructive
- Help others when possible
- Follow the code of conduct
- Participate in discussions

## Getting Help

- Join our community channels
- Ask questions in issues
- Read the documentation
- Contact maintainers if needed

---

Thank you for contributing to APEX! Your efforts help make the project better for everyone.