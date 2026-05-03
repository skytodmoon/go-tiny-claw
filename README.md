# Go Tiny Claw

A lightweight AI assistant framework for Go.

## Features

- Multi-provider support (Claude, Zhipu)
- Context management with token monitoring
- Tool registry with built-in tools
- File-based memory storage
- Feishu integration

## Getting Started

```bash
go run cmd/claw/main.go
```

## Project Structure

```
go-tiny-claw/
├── cmd/
│   └── claw/
│       └── main.go
├── internal/
│   ├── engine/      # MainLoop core implementation
│   ├── provider/    # LLM provider adapters
│   ├── context/     # Context management
│   ├── tools/       # Tool registry and built-in tools
│   ├── memory/      # File-based memory storage
│   └── feishu/      # Feishu bot integration
├── go.mod
└── README.md
```