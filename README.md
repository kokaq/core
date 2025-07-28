<div align="center">
  <img height="300" src="https://github.com/kokaq/.github/blob/main/kokaq-core.png?raw=true" alt="cute quokka as kokaq logo"/>
</div>

`core` is the heart of `kokaq`, a distributed, cloud-native priority queue implementation. It provides the foundational logic and data structures for enabling true, weight-based prioritization across distributed systems.

[![Go Reference](https://pkg.go.dev/badge/github.com/kokaq/core.svg)](https://pkg.go.dev/github.com/kokaq/core)
[![Tests](https://github.com/kokaq/core/actions/workflows/go.yml/badge.svg)](https://github.com/kokaq/core/actions/workflows/go.yml)


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
ns := queue.NewNamespace("./data/db", queue.NamespaceConfig{
 NamespaceName: "data-db",
 NamespaceId:   1,
})

logger.ConsoleLog("INFO", "Namespace created: #", ns.Id, ": ", ns.Name)

qConfig := queue.QueueConfiguration{
 QueueName:       "test-queue",
 QueueId:         1,
 EnableDLQ:       true,
 EnableInvisible: true,
}
var q *queue.Queue
var err error

q, err = ns.AddQueue(&qConfig)
var qi = &queue.QueueItem{
 MessageId: uuid.New(),
 Priority:  1,
}

err = q.Enqueue(qi)
qi, err := q.Peek()
qi, err := q.Dequeue()
ns.DeleteQueue(1)
```

For network server, see [server](https://github.com/kokaq/server).

## 🧪 Running Tests

```bash
go test ./...
```

## 🧱 Contributing

Contributions welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for code style and testing requirements.

## 📜 License

[MIT](./LICENSE) — open-source and production-ready.
