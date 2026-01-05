# Release Workflow Fix - 2026-01-05

## Issue

The GitHub Actions release workflow was failing with a **403 Forbidden** error when attempting to create a release:

```
⚠️ GitHub release failed with status: 403
undefined
❌ Too many retries. Aborting...
```

## Root Cause

The workflow was missing the `permissions` block that grants write access to repository contents. By default, the `GITHUB_TOKEN` only has read permissions:

```
Contents: read
Metadata: read
Packages: read
```

However, creating a GitHub release requires `contents: write` permission.

## Solution

Added the `permissions` block to the release workflow (`.github/workflows/release.yml`):

```yaml
on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write  # ← Added this

jobs:
  create-release:
    name: Create Release
    runs-on: windows-latest
    # ...
```

## Status

✅ **Fixed** - The workflow has been updated and pushed to the `main` branch.

## Next Steps

To create a successful release with the fixed workflow:

### Option 1: Create a New Tag (Recommended)
```powershell
# Delete the old tag locally and remotely
git tag -d v1.1.0
git push origin :refs/tags/v1.1.0

# Create a new tag on the latest commit (which includes the fix)
git tag v1.1.0
git push origin v1.1.0
```

### Option 2: Create a New Version
```powershell
# Create a new version tag
git tag v1.1.1
git push origin v1.1.1
```

The release workflow will automatically run when the tag is pushed and should now succeed in creating the GitHub release with all the build artifacts.

## Related Files
- `.github/workflows/release.yml` - Fixed workflow file
- Tag: `v1.1.0` - Current tag (needs to be recreated or replaced)
