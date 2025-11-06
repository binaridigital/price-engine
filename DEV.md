# Developer Guide - Feature Branch Workflow

This document outlines the development workflow, branching strategy, and best practices for contributing to the Price Engine project.

## 📋 Table of Contents

- [Branch Strategy](#branch-strategy)
- [Feature Branch Workflow](#feature-branch-workflow)
- [Development Process](#development-process)
- [Testing Requirements](#testing-requirements)
- [Code Review Process](#code-review-process)
- [Merging to Main](#merging-to-main)
- [Best Practices](#best-practices)

---

## 🌿 Branch Strategy

### Branch Types

1. **`main`** - Production-ready code only
   - Protected branch (requires PR and approvals)
   - All commits must be tested and validated
   - Never commit directly to `main`

2. **`develop`** (optional) - Integration branch for features
   - Used if you want a staging branch before `main`
   - Can be omitted if working directly to `main` via PRs

3. **`feature/*`** - Feature development branches
   - Format: `feature/description` or `feature/ticket-number-description`
   - Examples:
     - `feature/multi-exchange-support`
     - `feature/coinbase-integration`
     - `feature/order-book-level-2`
     - `feature/FEAT-123-latency-tracking`

4. **`bugfix/*`** - Bug fix branches
   - Format: `bugfix/description` or `bugfix/ticket-number-description`
   - Examples:
     - `bugfix/fix-memory-leak`
     - `bugfix/BUG-456-websocket-reconnect`

5. **`hotfix/*`** - Critical production fixes
   - Format: `hotfix/description`
   - Used for urgent fixes that need to go directly to production
   - Examples: `hotfix/critical-security-patch`

---

## 🔄 Feature Branch Workflow

### Step 1: Create Feature Branch

**Always start from an up-to-date `main` branch:**

```bash
# Ensure you're on main and it's up to date
git checkout main
git pull origin main

# Create and switch to new feature branch
git checkout -b feature/your-feature-name

# Example:
git checkout -b feature/coinbase-exchange-integration
```

**Branch Naming Conventions:**
- Use lowercase letters
- Separate words with hyphens (`-`)
- Be descriptive but concise
- Include ticket number if using issue tracking
- Examples:
  - ✅ `feature/coinbase-integration`
  - ✅ `feature/FEAT-123-latency-tracking`
  - ✅ `feature/order-book-level-2`
  - ❌ `feature/new-stuff`
  - ❌ `feature/update`

### Step 2: Develop Your Feature

**Work on your feature branch:**

```bash
# Make your changes
# ... edit files ...

# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add Coinbase exchange connector

- Implement Coinbase WebSocket connection
- Add trade data ingestion
- Handle reconnection logic
- Add unit tests

Closes #123"

# Push to remote
git push origin feature/your-feature-name
```

**Commit Message Format:**
Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(exchanges): add Coinbase connector
fix(aggregator): resolve memory leak in window management
docs(readme): update installation instructions
test(ingest): add unit tests for Binance connector
```

### Step 3: Keep Branch Updated

**Regularly sync with main to avoid conflicts:**

```bash
# Fetch latest changes
git fetch origin

# Rebase your feature branch on top of main
git checkout feature/your-feature-name
git rebase origin/main

# If conflicts occur, resolve them:
# 1. Fix conflicts in files
# 2. git add <resolved-files>
# 3. git rebase --continue

# Force push (safe on feature branches)
git push origin feature/your-feature-name --force-with-lease
```

**Alternative: Merge instead of rebase (if team prefers):**
```bash
git checkout feature/your-feature-name
git merge origin/main
git push origin feature/your-feature-name
```

### Step 4: Testing & Validation

**Before creating PR, ensure:**

1. **Unit Tests Pass:**
   ```bash
   go test ./...
   ```

2. **Build Successfully:**
   ```bash
   go build ./cmd/aggregator
   ```

3. **Linting:**
   ```bash
   golangci-lint run
   ```

4. **Manual Testing:**
   - Test your feature locally
   - Verify it works as expected
   - Check for edge cases

5. **Update Documentation:**
   - Update README if needed
   - Add/update code comments
   - Update API documentation

### Step 5: Create Pull Request

**Create PR from feature branch to main:**

```bash
# Push your branch (if not already pushed)
git push origin feature/your-feature-name
```

**Then create PR via GitHub/GitLab UI or CLI:**

```bash
# Using GitHub CLI
gh pr create --title "feat: Add Coinbase exchange integration" \
  --body "## Description
  Adds Coinbase exchange connector for real-time price data.

  ## Changes
  - Implement Coinbase WebSocket connection
  - Add trade data ingestion
  - Handle reconnection logic
  - Add comprehensive unit tests

  ## Testing
  - [x] Unit tests pass
  - [x] Manual testing completed
  - [x] Documentation updated

  Closes #123" \
  --base main
```

**PR Template Checklist:**
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass locally
- [ ] No merge conflicts
- [ ] Linked to related issue/ticket

---

## 🧪 Testing Requirements

### Unit Tests

**Every feature must include unit tests:**

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/ingest/...

# Run tests with verbose output
go test -v ./...
```

**Test Coverage:**
- Aim for >80% code coverage
- Critical paths should have >90% coverage
- New code must have tests

### Integration Tests

**For exchange connectors, add integration tests:**

```bash
# Run integration tests (may require API keys)
go test -tags=integration ./...
```

**Note:** Integration tests should be:
- Marked with build tag `//go:build integration`
- Optional (can skip if API keys not available)
- Documented in README

### Manual Testing Checklist

Before submitting PR:
- [ ] Feature works as expected
- [ ] Error handling works correctly
- [ ] Edge cases handled
- [ ] Performance is acceptable
- [ ] No memory leaks
- [ ] Logging is appropriate

---

## 👥 Code Review Process

### For Authors

1. **Self-Review First:**
   - Review your own code before requesting review
   - Check for typos, formatting, logic errors
   - Ensure tests are comprehensive

2. **Request Review:**
   - Assign reviewers (at least 1-2)
   - Add relevant labels
   - Provide context in PR description

3. **Respond to Feedback:**
   - Address all comments
   - Ask questions if unclear
   - Update PR as needed

### For Reviewers

1. **Review Checklist:**
   - [ ] Code follows style guidelines
   - [ ] Logic is correct and efficient
   - [ ] Error handling is appropriate
   - [ ] Tests are adequate
   - [ ] Documentation is updated
   - [ ] No security issues
   - [ ] Performance considerations addressed

2. **Provide Constructive Feedback:**
   - Be specific about issues
   - Suggest improvements
   - Approve when satisfied

3. **Approval:**
   - At least 1 approval required (2 for critical changes)
   - All CI checks must pass
   - No blocking comments

---

## ✅ Merging to Main

### Pre-Merge Checklist

**Before merging, ensure:**

- [ ] All tests pass (unit, integration, e2e)
- [ ] Code review approved (1-2 reviewers)
- [ ] CI/CD pipeline passes
- [ ] No merge conflicts
- [ ] Documentation updated
- [ ] Changelog updated (if applicable)
- [ ] Feature is complete and tested

### Merge Process

**Option 1: Squash and Merge (Recommended)**
- Combines all commits into one
- Cleaner git history
- Use for feature branches

**Option 2: Merge Commit**
- Preserves commit history
- Use if commit history is important

**Option 3: Rebase and Merge**
- Linear history
- Use if team prefers linear history

**After Merge:**
```bash
# Update local main
git checkout main
git pull origin main

# Delete feature branch (local)
git branch -d feature/your-feature-name

# Delete feature branch (remote)
git push origin --delete feature/your-feature-name
```

---

## 📝 Best Practices

### Code Quality

1. **Follow Go Conventions:**
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
   - Use `golint` or `golangci-lint`

2. **Error Handling:**
   - Always handle errors explicitly
   - Provide context in error messages
   - Use `fmt.Errorf` with `%w` for wrapping

3. **Documentation:**
   - Comment exported functions/types
   - Add package-level documentation
   - Update README for user-facing changes

4. **Performance:**
   - Profile before optimizing
   - Consider memory allocations
   - Use appropriate data structures

### Git Best Practices

1. **Commit Often:**
   - Small, logical commits
   - Each commit should be functional
   - Use meaningful commit messages

2. **Don't Commit:**
   - API keys or secrets
   - Large binary files
   - Generated files (unless necessary)
   - IDE-specific files

3. **Use .gitignore:**
   - Keep `.gitignore` updated
   - Don't commit temporary files

### Branch Management

1. **Keep Branches Small:**
   - One feature per branch
   - Easier to review and test
   - Faster to merge

2. **Delete Old Branches:**
   - Delete merged branches
   - Keep repository clean
   - Use naming conventions

3. **Sync Regularly:**
   - Rebase/merge from main frequently
   - Avoid large merge conflicts
   - Keep branches up to date

---

## 🚨 Emergency Hotfix Process

**For critical production issues:**

```bash
# Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-fix

# Make fix
# ... fix the issue ...

# Commit and push
git commit -m "hotfix: fix critical security issue"
git push origin hotfix/critical-fix

# Create PR to main (fast-track review)
# After merge, tag release
git tag -a v1.0.1 -m "Hotfix: Critical security patch"
git push origin v1.0.1
```

---

## 📚 Additional Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/)

---

## ❓ FAQ

**Q: Can I commit directly to main?**
A: No, all changes must go through feature branches and PRs.

**Q: How long should a feature branch live?**
A: Keep branches short-lived (days to weeks). Long-lived branches increase merge conflicts.

**Q: What if I need to update a merged feature?**
A: Create a new feature branch or bugfix branch for the update.

**Q: Can I work on multiple features in one branch?**
A: No, keep features separate. It makes review and testing easier.

**Q: What if my branch has conflicts with main?**
A: Rebase or merge from main, resolve conflicts, and continue.

---

*Last Updated: [Current Date]*
*Version: 1.0*

