<div align="center">
  <img height="300" src="https://github.com/kokaq/.github/blob/main/kokaq-core.png?raw=true" alt="cute quokka as kokaq logo"/>
</div>

`kokaq-core` is the heart of `kokaq`, a distributed, cloud-native priority queue implementation. It provides the foundational logic and data structures for enabling true, weight-based prioritization across distributed systems.


<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/kokaq/kokaq-core.svg)](https://pkg.go.dev/github.com/kokaq/kokaq-core) -->
<!-- [![Tests](https://github.com/kokaq/kokaq-core/actions/workflows/test.yml/badge.svg)](https://github.com/kokaq/kokaq-core/actions/workflows/test.yml) -->

## 🔍 What is kokaq-core?

`kokaq-core` contains all domain logic for the `kokaq` platform:
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

## 📁 Package Structure

```
pkg/
├── queue/ # Enqueue/dequeue interfaces and backends
├── scheduler/ # Priority/fairness strategies (Pending)
├── wireframe/ # Task models and serialization logic (Pending)
└── utils/ # Utility methods 

internal/
├── profiler/ # CPU, mem, IO profilers
└── metrics/ # Internal metrics helpers (no exporters) (Pending)
```

## 🚀 Getting Started

```bash
go get github.com/kokaq/core
```
Import and use in your server:

```go
import (
    "github.com/kokaq/kokaq-core/pkg/queue"
)
queueNs, _ := queue.NewKokaq(namespaceId, queueId)
err := queueNs.PushItem(queue.NewQueueItem(uuid.New(), priority))
item, err := q.PopItem()
```
For network server, see [kokaq-server](https://github.com/kokaq/kokaq-server).

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

