# AGENTS.md - Documentation (docs/)

This directory contains all project documentation in Markdown format.

## Documentation Structure

```text
docs/
├── README.md              # Documentation index (this file's parent)
├── GETTING_STARTED.md     # Installation, commands, quick start
├── DEMO.md                # Real-world examples and use cases
├── DEVELOPER.md           # Architecture, development guide
├── CONTRIBUTING.md        # How to contribute
├── FEATURE_ROADMAP.md     # Planned features, design philosophy
├── AUTO_UPDATE.md         # Auto-update functionality
├── GPU_PASSTHROUGH.md     # GPU-P feature guide
├── AI_AGENT.md            # AI agent setup and configuration
├── SKILLS_REFERENCE.md    # Quick reference for all AI skills
└── WORKFLOWS_REFERENCE.md # Quick reference for all AI workflows
```

## Documentation Standards

### Language & Style

- **Language**: English only (EN-US)
- **Tone**: Clear, concise, professional
- **Audience**: Developers and power users

### Formatting Guidelines

```markdown
# Document Title

Brief introduction paragraph.

## Section Header

Content with proper formatting.

### Subsection

- Bullet points for lists
- Use `inline code` for commands/paths
- Use code blocks for examples

## Code Examples

```powershell
quickvm start 1
```

## Tables

| Column 1 | Column 2 |
|----------|----------|
| Data     | Data     |

### Must-Have Elements

1. **Title** - Clear, descriptive `# Header`
2. **Introduction** - What this doc covers
3. **Prerequisites** (if applicable)
4. **Step-by-step instructions** (for guides)
5. **Examples** - Real, working examples
6. **See Also** - Links to related docs

## Adding New Documentation

1. Create `docs/NEW_FEATURE.md`
2. Follow naming: `UPPERCASE_WITH_UNDERSCORES.md`
3. Add to `docs/README.md` index table
4. Cross-link from related documents
5. Include in `docs/AI_AGENT.md` if AI-relevant

## Updating Existing Docs

When code changes affect documentation:

1. Update relevant docs immediately
2. Check `GETTING_STARTED.md` command lists
3. Update `DEVELOPER.md` if architecture changes
4. Bump "Last Updated" date if present

## Emoji Standards

Use consistent emoji for visual cues:

| Emoji | Usage |
|-------|-------|
| 📚 | Documentation sections |
| 🎯 | Quick navigation / goals |
| ⚡ | Performance / speed |
| 🔧 | Configuration / tools |
| ✅ | Success / completed |
| ❌ | Error / don't do |
| ⚠️ | Warning / caution |
| 💡 | Tips / hints |
| 🚀 | Getting started / launch |

## Command Documentation Format

When documenting CLI commands:

```markdown
### `quickvm <command>`

**Usage**: `quickvm <command> [options] <args>`

**Description**: Brief description of what it does.

**Arguments**:
- `<arg1>` - Description (required)
- `[arg2]` - Description (optional)

**Options**:
- `-f, --force` - Force operation
- `-v, --verbose` - Verbose output

**Examples**:
```powershell
# Basic usage
quickvm start 1

# With options
quickvm start 1 --force
```

## Screenshots & Media

If adding screenshots:

1. Save to `docs/assets/` directory
2. Use descriptive filenames: `tui-vm-list.png`
3. Reference with relative paths: `![TUI](assets/tui-vm-list.png)`
4. Keep file sizes reasonable (<500KB)

## Version Documentation

When releasing new versions:

1. Update `CHANGELOG.md` in root
2. Update version in `GETTING_STARTED.md`
3. Add migration notes if breaking changes

## AI Agent Relevance

This documentation is consumed by:

- Human developers
- AI coding agents (via AGENTS.md context)

Keep documentation **structured and scannable** for both audiences.
