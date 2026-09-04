# Security Policy

The QuickVM team and community take the security of our application, users,
and systems seriously. This document outlines our security policies, supported
versions, and instructions for reporting vulnerabilities responsibly.

---

## Supported Versions

Security patches and bug fixes are prioritized for the latest minor and patch
releases of active versions.

| Version | Supported          | Security Updates |
| ------- | ------------------ | ---------------- |
| 1.4.x   | :white_check_mark: | Active support   |
| 1.x     | :white_check_mark: | Critical patches |
| < 1.0   | :x:                | Unsupported      |

---

## Reporting a Vulnerability

If you discover a security vulnerability in QuickVM, please **do not open a
public issue or discuss it in public channels**. Public disclosure puts other
users at risk before a fix is made available.

### How to Report

Please report vulnerabilities privately using one of the following methods:

1. **GitHub Private Vulnerability Reporting (Preferred)**:
   - Navigate to the
     [QuickVM Security Advisories](https://github.com/hoangtran1411/quickvm/security/advisories)
     page.
   - Click **"Report a vulnerability"** to submit your findings confidentially.

2. **Direct Maintainer Contact**:
   - If GitHub private reporting is unavailable, reach out directly to the
     project lead [Hoang Tran (@hoangtran1411)](https://github.com/hoangtran1411)
     via private messaging on GitHub.

### What to Include in Your Report

To help us investigate and remediate the issue quickly, please include:

- **Summary**: High-level description of the vulnerability and impact.
- **Component**: The affected command or package (e.g., `cmd/rdp.go`,
  `internal/hyperv/export.go`).
- **Reproduction Steps**: Step-by-step instructions or commands.
- **Proof of Concept (PoC)**: Minimal demonstration of the behavior.
- **Environment**: OS version (e.g., Windows 11 23H2), QuickVM and Go version.
- **Remediation**: Any suggestions or potential fixes you may have.

### Response Timeline & SLA

- **Initial Response**: We aim to acknowledge receipt within **48 hours**.
- **Assessment & Triage**: We assess and verify within **5 business days**.
- **Fix & Disclosure**: A fix will be developed, tested, and released as
  quickly as feasible (typically 14–30 days). We adhere to coordinated
  vulnerability disclosure and coordinate release dates with the reporter.

---

## Threat Model & Security Considerations

Because QuickVM manages system virtualization and runs with elevated privileges
on Windows, the following threat vectors are actively defended against:

### 1. Elevated Privileges (Administrator Rights)

QuickVM requires Administrator privileges to communicate with Hyper-V host
management APIs and WMI services.

- QuickVM does not spawn unnecessary external processes with ambient or
  unconstrained privileges.
- Execution is strictly scoped to Hyper-V management workflows.

### 2. PowerShell Command Injection Prevention

QuickVM interfaces with Hyper-V cmdlets (`Get-VM`, `Start-VM`, `Stop-VM`).

- **Policy**: User inputs (such as VM names, snapshot labels, and file paths)
  must **never** be concatenated into raw PowerShell command strings.
- **Enforcement**: Commands are executed using isolated arguments with
  `exec.CommandContext` and strict argument separation to prevent arbitrary
  code execution via shell metacharacters.

### 3. Credential & Sensitive Data Safety

- Commands that accept credentials (such as `quickvm rdp`) must never log
  plaintext passwords, connection tokens, or secrets to terminal outputs, log
  files, or error messages.
- Temporary files created during execution must be secured and cleaned up.

### 4. Path Traversal & Filesystem Safeguards

- Commands that handle export, import, or disk image paths (`quickvm export`,
  `quickvm import`) validate input paths to prevent directory traversal
  attacks (`../`) that could alter or overwrite unauthorized system directories.

### 5. Binary Integrity & Auto-Updates

- QuickVM provides self-updating capabilities (`quickvm update`).
- Release binaries published on GitHub Releases are accompanied by
  cryptographically strong SHA-256 checksums (`.sha256`), allowing verification
  of binary integrity before execution.

---

## Hall of Fame & Recognition

We deeply appreciate the efforts of security researchers and community members
who help keep QuickVM and its users secure. If you responsibly report a
confirmed vulnerability, we will gladly credit you in our release notes and
security advisories (unless you prefer to remain anonymous).
