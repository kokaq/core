<div align="center">
  <img height="300" src="https://github.com/kokaq/.github/blob/main/kokaq-core.png?raw=true" alt="cute quokka as kokaq logo"/>
</div>

`core` is the heart of `kokaq`, a distributed, cloud-native priority queue implementation. It provides the foundational logic and data structures for enabling true, weight-based prioritization across distributed systems.


<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/kokaq/core.svg)](https://pkg.go.dev/github.com/kokaq/core) -->
<!-- [![Tests](https://github.com/kokaq/core/actions/workflows/test.yml/badge.svg)](https://github.com/kokaq/core/actions/workflows/test.yml) -->

## 🔍 What is core?

`core` contains all domain logic for the `kokaq` platform:
- Pluggable storage backends for task persistence
- Scheduler strategies for priority and fairness
- Task and message wireframe implementations
- Profiler utilities for CPU/mem/IO usage
- Proxy/middleware decorators for task flow interception

> ⚠️ This module is **network-agnostic**. No sockets, gRPC, HTTP, or Redis clients.

## ✨ Features

- 🧩 **Composable**: Cleanly separated storage, scheduling, and message logic
- ⚖️ **Priority-aware**: Built-in weighted fair and aging schedulers with deterministic dequeueing
- 🧪 **Fully testable**: Deterministic, goroutine-safe, and race-free
- 📦 **Tiny API surface**: Focused only on core primitives
- 📊 **Profiler-ready**: Built-in hooks for benchmarking and introspection
- ⚙️ **Lightweight core**: Pure logic with no external dependencies
- 🧵 **Concurrency-aware**: Built to scale across distributed workers
- 📦 **Modular**: Can be embedded into larger systems or composed into services


## 🚀 Getting Started

```bash
go get github.com/kokaq/core
```
Import and use in your server:

```go
import (
    "github.com/kokaq/core/queue"
)
queueNs, _ := queue.NewDefaultKokaq(namespaceId, queueId)
err := queueNs.PushItem(queue.NewQueueItem(uuid.New(), priority))
item, err := q.PopItem()
```
For network server, see [server](https://github.com/kokaq/server).

## 🧪 Running Tests

```bash
go test ./...
go test -race ./...
```
To run benchmarks:
```bash
go test -bench=. ./pkg/queue
```

## 🧱 Contributing

Contributions welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for code style and testing requirements.

## 📜 License

[MIT](./LICENSE) — open-source and production-ready.

