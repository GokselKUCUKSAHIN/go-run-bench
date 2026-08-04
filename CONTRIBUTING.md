# Contributing to go-run-bench

Thanks for your interest. This project stays small on purpose.

## Design principles

- **0 dependencies** — Go standard library only. New third-party imports will be rejected.
- **Simple stupid design** — Prefer clear, boring code over clever abstractions.
- **Blazing fast execution** — Keep the happy path thin: find files → run `go test -bench` → save results.

If a change fights these principles, it probably does not belong here.

## How to contribute

1. Open an issue first for non-trivial changes (behavior, flags, output format).
2. Fork and create a branch from `main`.
3. Keep the diff focused. One concern per PR.
4. Make sure the tool still builds:

```sh
go build -o go-run-bench .
```

5. If you touch CLI flags or output behavior, update [README.md](./README.md) in the same PR.

## Benchmark file convention

The tool discovers files named `*_benchmark_test.go` under the working directory.

## Pull request checklist

- [ ] No new dependencies in `go.mod`
- [ ] Design stays simple — no framework, no package sprawl for its own sake
- [ ] `go build` succeeds
- [ ] README updated when user-facing behavior changes
- [ ] You agree to the [Code of Conduct](./CODE_OF_CONDUCT.md)

## Code of Conduct

By participating, you agree to follow the [Code of Conduct](./CODE_OF_CONDUCT.md).

## Security

Do not open a public issue for vulnerabilities. See [SECURITY.md](./SECURITY.md).
