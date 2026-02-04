---
description: Maintain and improve Agent Experience (AX) documentation
---

# Improve Agent Experience (AX)

This workflow helps maintain and enhance AI agent documentation.

## When to Use

- After adding new features
- When AI agents make repeated mistakes
- During periodic maintenance
- After major refactoring

---

## Step 1: Audit Current AX Files

// turbo
```powershell
# List all AGENTS.md files
Get-ChildItem -Path . -Filter "AGENTS.md" -Recurse | Select-Object FullName
```

// turbo
```powershell
# List context files
Get-ChildItem -Path .\.agent\context\ -Filter "*.md" | Select-Object Name
```

---

## Step 2: Check for Outdated Information

Review these areas:
- [ ] Commands in `AGENTS.md` match actual CLI
- [ ] File paths in AGENTS.md files are accurate
- [ ] Code snippets compile and work
- [ ] Known issues are still relevant

// turbo
```powershell
# Verify commands exist
quickvm --help
```

---

## Step 3: Update from Recent Changes

Check recent commits for documentation impact:

// turbo
```powershell
git log --oneline -10 --name-only
```

For each changed file, ask:
- Does this affect any AGENTS.md?
- Should known-issues.md be updated?
- Are new snippets needed?

---

## Step 4: Add Missing AGENTS.md Files

Check for directories without AGENTS.md:

```powershell
# Directories that might need AGENTS.md
# - Any new package under internal/
# - Any new command group under cmd/
```

Template for new AGENTS.md:

```markdown
# AGENTS.md - [Directory Name]

Brief description of this package/directory.

## Purpose

What this code does.

## Key Files

| File | Purpose |
|------|---------|
| file.go | Description |

## Patterns

Common patterns used here.

## Adding New Features

How to extend this code.
```

---

## Step 5: Update Known Issues

Review recent bug fixes and add preventive documentation:

```powershell
git log --oneline --grep="fix" -10
```

Add to `.agent/context/known-issues.md`:
- New gotchas discovered
- Common mistakes made
- Edge cases found

---

## Step 6: Enhance Snippets

Check if frequently-used code patterns are documented:

```powershell
# Look for common patterns in codebase
Select-String -Path ".\**\*.go" -Pattern "func.*error" | Select-Object -First 5
```

Add useful snippets to `.agent/context/snippets.md`.

---

## Step 7: Validate Documentation

// turbo
```powershell
# Check markdown syntax
Get-ChildItem -Path . -Filter "*.md" -Recurse | ForEach-Object { Write-Host $_.FullName }
```

// turbo
```powershell
# Ensure no broken links (basic check)
Select-String -Path ".\docs\*.md" -Pattern "\[.*\]\(.*\.md\)" | Select-Object LineNumber, Line
```

---

## Step 8: Update AI_AGENT.md Index

Ensure `docs/AI_AGENT.md` reflects current structure:
- [ ] All AGENTS.md files listed
- [ ] All skills documented
- [ ] All workflows documented
- [ ] Structure diagram is accurate

---

## Checklist

- [ ] All AGENTS.md files are up-to-date
- [ ] known-issues.md reflects recent discoveries
- [ ] decisions.md has new ADRs if applicable
- [ ] snippets.md has useful new patterns
- [ ] docs/AI_AGENT.md structure is current
- [ ] No broken internal links
- [ ] All code examples work

---

## Maintenance Schedule

Recommended frequency:
- **Weekly**: Quick review of known-issues
- **Per Release**: Full AX audit
- **After Major Refactor**: Complete update

---

## See Also

- [AI Agent Setup](../docs/AI_AGENT.md)
- [Skills Reference](../docs/SKILLS_REFERENCE.md)
- [Workflows Reference](../docs/WORKFLOWS_REFERENCE.md)
