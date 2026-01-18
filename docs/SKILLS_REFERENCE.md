# AI Skills Reference

Quick reference for all available AI skills in this project.

## Available Skills

### Go CLI with Cobra
**Location**: `.agent/skills/go-cli-cobra/SKILL.md`

**What it provides**:
- Root command template
- Resource command template
- Version command template
- Subcommand patterns (e.g., `app resource action`)
- Flag handling best practices

**When to use**:
- Adding new CLI commands
- Restructuring command hierarchy
- Implementing flags and arguments

---

### Go TUI with Bubble Tea
**Location**: `.agent/skills/go-tui-bubbletea/SKILL.md`

**What it provides**:
- Lipgloss style definitions
- Table model template
- Elm Architecture patterns (Model-Update-View)
- Color reference guide
- Keyboard handling patterns

**When to use**:
- Creating interactive terminal UIs
- Adding table/list views
- Styling terminal output

---

### Go PowerShell Integration
**Location**: `.agent/skills/go-powershell-integration/SKILL.md`

**What it provides**:
- ShellExecutor interface pattern
- Mock implementation for testing
- Security best practices (injection prevention)
- JSON output parsing patterns
- Error handling patterns

**When to use**:
- Adding PowerShell commands
- Windows system integration
- Mocking for tests

---

### Go Testing Patterns
**Location**: `.agent/skills/go-testing-patterns/SKILL.md`

**What it provides**:
- Table-driven test template
- Interface-based mocking
- Test helper patterns
- Build tags for integration tests
- GitHub Actions CI configuration

**When to use**:
- Writing unit tests
- Setting up CI/CD
- Mocking dependencies

---

### Go Project Structure
**Location**: `.agent/skills/go-project-structure/SKILL.md`

**What it provides**:
- Standard directory layout
- Makefile template
- GitHub Actions workflows
- README template
- CHANGELOG format

**When to use**:
- Starting new projects
- Restructuring existing projects
- Setting up CI/CD

---

## How Skills Work

1. **Discovery**: AI identifies relevant skills based on your request
2. **Reading**: AI reads the SKILL.md file
3. **Application**: AI applies patterns and templates from the skill

## Skill Triggers

| Request | Skill Applied |
|---------|---------------|
| "Add a command for X" | Go CLI with Cobra |
| "Create a table view" | Go TUI with Bubble Tea |
| "Execute PowerShell command" | Go PowerShell Integration |
| "Add tests for X" | Go Testing Patterns |
| "Set up CI/CD" | Go Project Structure |

## Creating Custom Skills

1. Create directory: `.agent/skills/<skill-name>/`
2. Create `SKILL.md` file:

```markdown
---
name: Skill Name
description: Brief description
---

# Skill Name

## When to Use
Description of use cases...

## Templates
Code templates and patterns...

## Best Practices
Guidelines and recommendations...
```

## Skill Best Practices

1. **Focus**: Each skill should cover one topic
2. **Examples**: Include working code examples
3. **Variations**: Document common variations
4. **References**: Link to external documentation
