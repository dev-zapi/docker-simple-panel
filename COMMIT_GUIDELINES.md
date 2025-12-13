# Git Commit Message Guidelines

This document outlines the commit message conventions for the Docker Simple Panel project.

## General Principles

- Write clear, concise commit messages that explain **what** changed and **why**
- Use the imperative mood in the subject line (e.g., "Add feature" not "Added feature")
- Keep the subject line under 72 characters
- Separate subject from body with a blank line
- Use the body to explain what and why vs. how

## Commit Message Format

```
<type>: <subject>

[optional body]

[optional footer]
```

### Subject Line

The subject line should be structured as follows:

**Format:** `<type>: <description>`

#### Types

Use one of the following types to categorize your commit:

- **feat**: A new feature for the user
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (formatting, missing semi-colons, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvements
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes to build system or external dependencies
- **ci**: Changes to CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit

#### Description

- Use imperative, present tense: "change" not "changed" nor "changes"
- Don't capitalize the first letter
- No period (.) at the end
- Be specific and descriptive

**Examples:**
- ✅ `feat: add container restart endpoint`
- ✅ `fix: resolve JWT token expiration issue`
- ✅ `docs: update API documentation for volumes endpoint`
- ❌ `Added new feature` (uses past tense, not imperative)
- ❌ `Fix bug` (not descriptive enough)
- ❌ `Updates` (too vague, not imperative)

### Body (Optional)

- Use the body to explain **what** and **why**, not **how**
- Wrap lines at 72 characters
- Can include multiple paragraphs separated by blank lines
- Can use bullet points (use `-` or `*`)

**Example:**
```
fix: prevent race condition in container status updates

The container status polling was causing race conditions when multiple
requests were made simultaneously. This fix:

- Adds mutex locks around container state reads
- Implements read/write lock pattern for better concurrency
- Ensures atomic updates to container status cache

This resolves intermittent 500 errors when rapidly refreshing the UI.
```

### Footer (Optional)

Use the footer for:
- **Breaking changes**: Start with `BREAKING CHANGE:` followed by description
- **Issue references**: Use `Fixes #123`, `Closes #123`, or `Refs #123`
- **Co-authors**: Use `Co-authored-by: Name <email>`

**Example:**
```
feat: change authentication to use JWT refresh tokens

BREAKING CHANGE: The login endpoint now returns both access and refresh
tokens. Clients must be updated to handle the new response format.

Fixes #45
Co-authored-by: John Doe <john@example.com>
```

## Examples

### Simple Feature Addition
```
feat: add volume deletion endpoint

Implements DELETE /api/volumes/{name} endpoint with proper
authentication and error handling.
```

### Bug Fix
```
fix: correct container health status display

Health status was showing as "unknown" for healthy containers due to
incorrect field mapping in the Docker API response.

Fixes #67
```

### Documentation Update
```
docs: add WebSocket log streaming examples

Added JavaScript and curl examples for connecting to the container
log streaming endpoint.
```

### Refactoring
```
refactor: extract Docker client initialization into factory

Moves Docker client setup logic into a dedicated factory function to
improve testability and reduce code duplication across handlers.
```

### Breaking Change
```
feat: migrate to Svelte 5 runes syntax

BREAKING CHANGE: The frontend now requires Svelte 5. All components
have been updated to use runes ($state, $derived, $effect) instead of
the old reactivity system.

Refs #89
```

### Multi-paragraph Body
```
fix: resolve WebSocket connection timeout on slow networks

The WebSocket upgrade was failing on networks with >500ms latency due
to overly aggressive timeout settings.

Changes:
- Increase WebSocket handshake timeout from 5s to 30s
- Add retry logic with exponential backoff
- Improve error messages for connection failures

This should improve reliability for users on slower connections or
behind corporate proxies.

Fixes #102
```

## Best Practices

### DO:
- ✅ Use the imperative mood ("add" not "added")
- ✅ Be specific and descriptive
- ✅ Reference issues when applicable
- ✅ Explain the "why" in the body for non-trivial changes
- ✅ Keep commits focused on a single logical change
- ✅ Add Co-authored-by for pair programming or collaborated work

### DON'T:
- ❌ Write vague messages like "fix bug" or "update code"
- ❌ Include multiple unrelated changes in one commit
- ❌ Use past tense or gerunds ("fixed", "fixing")
- ❌ End subject line with a period
- ❌ Exceed 72 characters in the subject line
- ❌ Leave the commit message empty or use placeholders like "WIP"

## Special Cases

### Initial Commits
```
chore: initial commit
```

### Merge Commits
For merge commits, GitHub will auto-generate messages. If writing manually:
```
Merge pull request #123 from user/feature-branch

feat: add container grouping by labels
```

### Revert Commits
```
revert: feat: add experimental caching layer

This reverts commit abc123def456.

The caching layer was causing data inconsistencies in production.
```

## Tools and Validation

### Commit Message Template

You can optionally configure a commit message template to help remind you of the format:

```bash
git config commit.template .gitmessage
```

Then create a `.gitmessage` file in your home directory or project root:
```
<type>: <subject>

# Why is this change needed?
# 

# How does it address the issue?
# 

# What are the side effects?
# 

# Fixes #
```

### Pre-commit Hooks

Consider using tools like [commitlint](https://commitlint.js.org/) to enforce these conventions automatically.

## References

These guidelines are inspired by:
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit)
- [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)

## Questions?

If you're unsure about how to format a commit message, look at recent commits in the repository for examples or ask in the pull request discussion.
