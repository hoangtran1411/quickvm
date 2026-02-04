# AI Agent Setup

This document explains how to use the AI agent configuration, skills, and workflows in this project.

## Overview

The `.agent/` directory contains configuration files that enhance AI coding assistants (like GitHub Copilot, Cursor, Claude, OpenAI Codex, Google Gemini, and others) with project-specific knowledge and automation.

This project follows the [AGENTS.md standard](https://agents.md/) - used by 60k+ open-source projects.

```
quickvm/
├── AGENTS.md                 # Root context for AI agents
├── cmd/AGENTS.md             # CLI commands context
├── internal/hyperv/AGENTS.md # Hyper-V logic context
├── ui/AGENTS.md              # TUI context
└── .agent/
    ├── context/              # Additional AI context
    │   ├── known-issues.md   # Edge cases and gotchas
    │   ├── decisions.md      # Architecture Decision Records
    │   └── snippets.md       # Copy-paste code patterns
    ├── rules/                # Always-on rules and guidelines
    │   ├── go-style-guide.md # Go coding standards
    │   └── go-idioms-reference.md # Detailed idioms
    ├── skills/               # Reusable patterns and templates
    │   ├── go-cli-cobra/
    │   ├── go-tui-bubbletea/
    │   ├── go-powershell-integration/
    │   ├── go-testing-patterns/
    │   ├── go-project-structure/
    │   └── git-release-management/
    └── workflows/            # Step-by-step task automation
        ├── add-command.md
        ├── add-hyperv-feature.md
        └── ...
```

## Rules

Rules in `.agent/rules/` are **always active** and automatically applied by the AI agent. They define:

- Code style guidelines
- Architecture patterns
- Security requirements
- Best practices

### Current Rules

| Rule | Description |
|------|-------------|
| `go-style-guide.md` | Go coding standards, error handling, Charmbracelet UI patterns |

## Skills

Skills are **reusable knowledge modules** that the AI can reference when relevant. Each skill provides:

- Templates and code patterns
- Best practices
- Implementation guidelines

### Available Skills

| Skill | Description | Use Case |
|-------|-------------|----------|
| **Go CLI with Cobra** | CLI application patterns | Adding commands, flag handling |
| **Go TUI with Bubble Tea** | Terminal UI patterns | Interactive interfaces |
| **Go PowerShell Integration** | PowerShell execution patterns | Windows automation |
| **Go Testing Patterns** | Testing best practices | Unit tests, mocking, CI |
| **Go Project Structure** | Project layout standards | New projects, restructuring |

### How Skills Work

When you ask the AI to implement something, it will:
1. Identify relevant skills
2. Read the skill's `SKILL.md` file
3. Apply the patterns and templates from the skill

**Example prompts that trigger skills:**

```
"Add a new command to list snapshots"
→ Triggers: Go CLI with Cobra skill

"Create a beautiful table in the TUI"
→ Triggers: Go TUI with Bubble Tea skill

"Add tests for the export function"
→ Triggers: Go Testing Patterns skill
```

## Workflows

Workflows are **step-by-step guides** for common tasks. They can be invoked using slash commands.

### Available Workflows

| Command | Description | Auto-run |
|---------|-------------|----------|
| `/dev-cycle` | Format, test, build cycle | ✅ All steps |
| `/add-command` | Add new CLI command | Partial |
| `/add-hyperv-feature` | Add Hyper-V functionality | Partial |
| `/add-tui-feature` | Add TUI features | Partial |
| `/test-coverage` | Generate coverage reports | ✅ All steps |
| `/refactor` | Safe refactoring with tests | ✅ All steps |
| `/fix-bug` | Debug and fix bugs | Partial |
| `/release` | Create new release | Partial |
| `/update-deps` | Update dependencies | ✅ All steps |
| `/review-pr` | Code review checklist | ✅ All steps |

### Auto-run Annotations

Workflows use special annotations to control automation:

- `// turbo` - Auto-run the next command without asking
- `// turbo-all` - Auto-run all commands in the workflow

### How to Use Workflows

Simply type the slash command in your AI chat:

```
/dev-cycle
```

The AI will execute the workflow steps automatically.

## Creating New Skills

To add a new skill:

1. Create directory: `.agent/skills/<skill-name>/`
2. Create `SKILL.md` with this format:

```markdown
---
name: Skill Name
description: Brief description for skill discovery
---

# Skill Name

Detailed instructions, templates, and patterns...
```

## Creating New Workflows

To add a new workflow:

1. Create file: `.agent/workflows/<workflow-name>.md`
2. Use this format:

```markdown
---
description: Brief description for workflow discovery
---

# Workflow Title

## Step 1: Description

// turbo
\`\`\`powershell
command-to-run
\`\`\`

## Step 2: Description

...
```

## Best Practices

### For Skills
- Keep skills focused on one topic
- Include working code examples
- Document common variations
- Reference external resources when helpful

### For Workflows
- Use clear step numbers
- Include verification steps
- Add checklists at the end
- Use `// turbo` for safe, repeatable commands

### For Rules
- Keep rules concise
- Focus on project-specific requirements
- Use clear examples of do's and don'ts

## Troubleshooting

### Skill not being applied
- Ensure the skill file is named `SKILL.md`
- Check the YAML frontmatter is valid
- Verify the skill directory is under `.agent/skills/`

### Workflow not found
- Ensure the file is in `.agent/workflows/`
- Check the YAML frontmatter has a `description` field
- Verify file extension is `.md`

### Commands not auto-running
- Check for `// turbo` annotation before the code block
- For all commands, add `// turbo-all` anywhere in the file

---

## Quick References

For quick lookup during development:

- **[Skills Reference](SKILLS_REFERENCE.md)** - All skills in one page
- **[Workflows Reference](WORKFLOWS_REFERENCE.md)** - All workflows in one page

---

## See Also

- [Developer Guide](DEVELOPER.md) - Project architecture
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Main README](../README.md) - Project overview
