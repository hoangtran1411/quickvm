# AI Workflows Quick Reference

Quick reference for all available AI workflows in this project.

## Development Workflows

### `/dev-cycle` - Standard Development Cycle
**Auto-run**: ✅ All steps

Fast iteration cycle:
```
1. Format code (go fmt)
2. Run tests
3. Build binary
```

### `/add-command` - Add New CLI Command
**Auto-run**: Partial

Steps:
1. Create command file in `cmd/`
2. Add business logic to `hyperv/` (if needed)
3. Write unit tests
4. Update documentation
5. Verify with `--help`

### `/add-hyperv-feature` - Add Hyper-V Feature
**Auto-run**: Partial

Steps:
1. Research PowerShell cmdlet
2. Create feature file in `hyperv/`
3. Add mock-based tests
4. Create CLI commands
5. Update documentation

### `/add-tui-feature` - Add TUI Feature
**Auto-run**: Build step only

Steps:
1. Identify change type (keyboard, column, style)
2. Update Model (if needed)
3. Add keyboard shortcut
4. Update help text
5. Manual testing

---

## Quality Workflows

### `/test-coverage` - Generate Coverage Reports
**Auto-run**: ✅ All steps

Steps:
1. Run tests with coverage
2. Generate coverage profile
3. Create HTML report
4. View summary

### `/refactor` - Safe Refactoring
**Auto-run**: ✅ Test/build steps

Steps:
1. Verify existing tests pass
2. Note coverage percentage
3. Make small changes
4. Run tests after each change
5. Verify no regression

### `/fix-bug` - Debug and Fix Bugs
**Auto-run**: Partial

Steps:
1. Reproduce the bug
2. Locate the issue
3. Write failing test
4. Fix the code
5. Verify all tests pass

### `/review-pr` - Code Review
**Auto-run**: ✅ Check steps

Checklist:
- Code style and formatting
- Architecture patterns
- Security (PowerShell injection)
- Testing coverage
- Documentation

---

## Maintenance Workflows

### `/release` - Create New Release
**Auto-run**: Partial

Steps:
1. Update version in `cmd/version.go`
2. Update CHANGELOG.md
3. Run all checks
4. Build release binaries
5. Create and push git tag
6. Verify GitHub Release

### `/update-deps` - Update Dependencies
**Auto-run**: ✅ All steps

Steps:
1. Check current dependencies
2. Check for updates
3. Update dependencies
4. Run `go mod tidy`
5. Run tests
6. Manual verification

### `/improve-ax` - Improve Agent Experience
**Auto-run**: ✅ Most steps

Steps:
1. Audit current AGENTS.md files
2. Check for outdated information
3. Update from recent changes
4. Add missing AGENTS.md files
5. Update known-issues.md
6. Enhance snippets.md
7. Validate documentation
8. Update AI_AGENT.md index

When to use:
- After adding new features
- When AI agents make repeated mistakes
- During periodic maintenance
- After major refactoring

---

## Workflow Annotations

| Annotation | Effect |
|------------|--------|
| `// turbo` | Auto-run the next command |
| `// turbo-all` | Auto-run all commands in workflow |

## Usage Tips

1. **Type the slash command** in AI chat to start a workflow
2. **Turbo commands run automatically** - no confirmation needed
3. **Non-turbo commands** will ask for confirmation
4. **Checklists** at end of workflows help verify completion

## Customizing Workflows

Workflows are in `.agent/workflows/`. To customize:

1. Edit the `.md` file
2. Add `// turbo` before commands you trust
3. Add `// turbo-all` at top for full automation
