# Contributing to kokaq-core

First off, thank you for taking the time to contribute! 🎉  
`kokaq-core` is the heart of the [kokaq](https://github.com/kokaq) platform — a distributed, cloud-native priority queue system.

This repository contains all core logic: scheduling, storage interfaces, task models, and deterministic prioritization logic — without any network bindings.

## 📜 Ground Rules

- **No networking code** in this repo. That belongs in [kokaq-server](https://github.com/kokaq/kokaq-server).
- Always write **unit tests** and run `go test -race ./...` before submitting a PR.
- Follow **Go naming and formatting** conventions. Run `gofmt` and `golangci-lint`.
- Be mindful of **public API surface**. Any breaking changes must be clearly discussed.
- Keep dependencies **minimal and standard**. No third-party packages unless justified.

## 🔧 Development Setup

```bash
git clone https://github.com/kokaq/kokaq-core.git
cd kokaq-core
go mod tidy
go test ./...
```

To enable coverage, race detector, and static checks:

```bash
make test
```

> You’ll need Go 1.20 or later.

## 📁 Project Structure

```bash
pkg/             # Public, importable packages (storage, scheduler, etc.)
internal/        # Internal helpers (clocks, profiler, metrics)
test/            # Shared test cases, fuzzing, and benchmarks
docs/            # Architecture references and diagrams
```

## 🧪 Testing

All logic must be covered by:
- Unit tests using testing package
- Optional fuzzing for unmarshaling/serialization
- Benchmarks for schedulers or hot paths

Run all tests:

```bash
go test ./...
```

Run fuzzing:

```bash
go test -fuzz=Fuzz -fuzztime=30s ./test/fuzz/...
```

## ✅ Commit Style

Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to allow automatic changelog generation:

```makefile
feat(scheduler): add weighted fair implementation
fix(storage): correct dequeue starvation logic
test: add fuzzing for wireframe decoding
```

## 🙋 Support & Questions

- For bugs or issues, open a [GitHub issue](https://github.com/kokaq/kokaq-core/issues)
- For design questions or roadmap discussion, open a [GitHub discussion](https://github.com/orgs/kokaq/discussions)

## 🙌 Thank You!

Your contributions help make `kokaq` a fast, flexible, and developer-first queuing platform.
