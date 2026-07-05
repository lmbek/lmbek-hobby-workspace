# Contributing to LMBEK Hobby Workspace

Thank you for considering contributing! This document explains how to get started and what we expect from contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)

---

## Code of Conduct

Be respectful, constructive, and professional. We are building software together — treat every contributor the way you would want to be treated.

---

## Getting Started

1. Fork the repository (or create a feature branch if you have write access).
2. Run `make doctor` to verify your environment.
3. Run `make clone` to materialise all managed repositories.
4. Make your changes in the appropriate repository.

---

## Development Workflow

1. **Create a branch** from `main`:
   ```
   git checkout -b feature/short-description
   ```
2. **Make small, focused commits** — one logical change per commit.
3. **Write or update tests** for any code changes in `git-controller/`.
4. **Run tests locally** before pushing:
   ```
   cd git-controller && go test ./...
   ```
5. **Push and open a Pull Request** against `main`.

---

## Pull Request Process

1. Fill in the PR template completely.
2. Ensure all CI checks pass (build, test, lint).
3. Request a review from at least one maintainer.
4. Address review feedback promptly.
5. Squash-merge is preferred for clean history.

### PR Size Guidelines

| Size       | Lines Changed | Review Time                               |
|------------|---------------|-------------------------------------------|
| Small      | < 200         | 2 weeks                                   |
| Medium     | 200–2000      | 4 weeks                                   |
| Large      | 2000+         | >2 months (full end to end test required) |


---

## Coding Standards

- Follow the existing code style in each repository.
- Go code must pass `go vet` and `go test` with no failures.
- Use meaningful variable and function names.
- Keep functions short and focused (single responsibility).
- See [GUIDELINES.md](./GUIDELINES.md) for the master guidelines.

---

## Commit Messages

Use the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <short description>

<optional body>

<optional footer>
```

### Types

| Type       | When to use                                      |
|------------|--------------------------------------------------|
| `feat`     | A new feature                                    |
| `fix`      | A bug fix                                        |
| `docs`     | Documentation only changes                       |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test`     | Adding or updating tests                         |
| `chore`    | Maintenance tasks (deps, CI, tooling)            |
| `other`    | Changes that don't fit any other type            |

### Examples

```
feat(status): add repository status dashboard command
fix(clone): handle empty repository URL gracefully
docs(readme): add architecture diagram section
chore(ci): upgrade Go version in workflow to 1.26
```
