---
name: git-release-management
description: Manages professional Git commits using Conventional Commits and creates well-formatted release tags. Use this skill when committing code, creating releases, or tagging versions.
---

# Git Release Management Skill 🚀

This skill defines the standards for professional Git commits and automated release tagging. Following these patterns ensures a clean, searchable history and professional-grade release notes for users and contributors.

## 📝 Conventional Commits

Always use the **Conventional Commits** specification for commit messages. This makes the history readable and allows for automated changelog generation.

### Format
`<type>(<scope>): <description>`

### Types
- `feat`: A new feature (e.g., `feat(ui): add dark mode support`)
- `fix`: A bug fix (e.g., `fix(updater): resolve permission denied on Windows`)
- `docs`: Documentation only changes (e.g., `docs: translate README to English`)
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit

### Guidelines
- Use the **imperative mood** ("add", not "added").
- Do not capitalize the first letter of the description.
- No period (.) at the end of the description.

---

## 🏷️ Professional Git Tagging

When creating a release tag, provide a structured and descriptive message that acts as a human-readable changelog.

### Versioning
- Use **Semantic Versioning** (`vMajor.Minor.Patch`).
- Always prefix with `v` (e.g., `v2.1.0`).

### Tag Message Template
Use emojis and logical sections to categorize changes:

```text
vX.Y.Z - [Highlight Title]

🚀 [Impactful Changes/Features]
- Point 1 explaining the 'Why' and 'What'.
- Point 2...

🛠️ [Improvements & Fixes]
- Fixed issue X by doing Y.
- Optimized Z for better performance.

📚 [Documentation & Meta]
- Updated README with instructions.
- Added contributing guidelines.

🤝 [Contributors/Community]
- Shout out to community members if applicable.
```

---

## 💻 Practical Usage Commands

### Professional Commit
```powershell
git add .
git commit -m "feat(api): optimize response time by caching results"
```

### Professional Tagging
```powershell
git tag -a v2.1.0 -m "v2.1.0 - Premium UI Revamp Release

✨ Major UI Update:
- Complete frontend redesign featuring a 3-column dashboard layout.
- Premium Glassmorphism theme with dynamic backgrounds.

🛠️ Improvements:
- CI/CD linting fixes for cross-platform support.
- Optimized build scripts for faster deployment."

# Push to origin
git push origin main --tags
```

---

## 💎 Premium Aesthetics in Release Notes
- Use **bold text** for keywords.
- Use **bullet points** for readability.
- Add a **summary line** at the top of the tag message.
- Group related changes together to tell a story of "Evolution".
