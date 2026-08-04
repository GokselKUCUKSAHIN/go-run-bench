# Security Policy

## Supported versions

| Version | Supported |
| ------- | --------- |
| latest release | Yes |
| `main` branch | Yes |
| older releases | No |

## Reporting a vulnerability

Please **do not** open a public GitHub issue for security vulnerabilities.

Prefer one of these channels:

1. [GitHub Security Advisories](https://github.com/GokselKUCUKSAHIN/go-run-bench/security/advisories/new) (preferred)
2. Private contact via [@GokselKUCUKSAHIN](https://github.com/GokselKUCUKSAHIN)

Include:

- Description of the issue
- Steps to reproduce
- Impact assessment (if known)
- Suggested fix (optional)

## Response

You should receive an acknowledgment within **7 days**. After triage, we will share next steps and credit reporters who want to be named when a fix ships.

## Scope notes

`go-run-bench` is a local CLI that walks a directory and runs `go test`. Treat untrusted project trees as untrusted input: only point `-wd` at repositories you trust.
