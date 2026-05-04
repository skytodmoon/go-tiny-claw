# Level 3 测试指南

## 目录

- [快速开始](#快速开始)
- [测试架构](#测试架构)
- [测试类型](#测试类型)
- [运行测试](#运行测试)
- [测试用例](#测试用例)

## 快速开始

### 使用 Makefile（推荐）

```bash
# 运行所有测试
make test

# 只运行单元测试
make test-unit

# 只运行集成测试
make test-integration

# 运行端到端测试
make test-e2e

# 运行带覆盖率的测试
make cover
```

### 手动运行

```bash
# 单元测试
go test -v ./internal/engine ./internal/tools

# 集成测试
go test -v ./integration

# Level 3 验证测试
go run cmd/level3-test/main.go
```

## 测试架构

```
go-tiny-claw/
├── cmd/
│   └── level3-test/        # Level 3 验证测试
│       └── main.go
├── internal/
│   ├── engine/
│   │   ├── loop.go         # 主引擎
│   │   └── loop_test.go    # 引擎单元测试
│   ├── tools/
│   │   ├── read_file.go
│   │   ├── write_file.go
│   │   ├── edit_file.go
│   │   └── tools_test.go   # 工具单元测试
│   └── provider/
├── integration/
│   └── level3_integration_test.go  # Level 3 集成测试
├── testutils/              # 测试工具包
│   ├── mock_provider.go    # Mock Provider
│   └── test_helpers.go     # 测试辅助函数
├── scripts/
│   └── test-e2e.sh         # 端到端测试脚本
└── Makefile
```

## 测试类型

### 1. 单元测试 (`test-unit`)

**位置**: `internal/engine/loop_test.go`, `internal/tools/tools_test.go`

**测试内容**:
- `truncate()` 函数的边界测试
- `ReadFileTool` 的读写测试
- `WriteFileTool` 的写入测试
- `EditFileTool` 的编辑测试

**运行**:
```bash
go test -v ./internal/engine ./internal/tools
```

### 2. 集成测试 (`test-integration`)

**位置**: `integration/level3_integration_test.go`

**测试内容**:
- Level 3 四阶段循环（Thinking → Acting → Observation → Re-thinking）
- 多工具调用顺序
- 读写编辑组合流程
- 上下文管理

**运行**:
```bash
go test -v ./integration
```

### 3. Level 3 验证测试 (`level3-test`)

**位置**: `cmd/level3-test/main.go`

**测试内容**:
- 完整 Level 3 功能验证
- 真实 LLM Provider 调用
- 文件创建和读取
- 执行摘要统计

**运行**:
```bash
export SILICONFLOW_API_KEY="your-api-key"
go run cmd/level3-test/main.go
```

### 4. 端到端测试 (`test-e2e`)

**位置**: `scripts/test-e2e.sh`

**测试内容**:
- 真实环境完整流程
- 自定义 Prompt 测试
- 集成验证

**运行**:
```bash
export SILICONFLOW_API_KEY="your-api-key"
make test-e2e
```

## 运行测试

### 环境要求

- Go 1.21+
- `SILICONFLOW_API_KEY` 环境变量（用于 e2e 测试）

### 环境变量

```bash
export SILICONFLOW_API_KEY="your-api-key"
```

## 测试用例

### 集成测试用例

| 测试名称 | 描述 |
|---------|------|
| `TestLevel3ReadWriteCycle` | 测试读取 -> 写入 -> 完成循环 |
| `TestLevel3EditCycle` | 测试读取 -> 编辑 -> 完成循环 |
| `TestLevel3MultipleTools` | 测试多工具调用序列 |

### Level 3 验证用例

| 用例名称 | 描述 |
|---------|------|
| Read and Write | 读取 README.md 并创建 SUMMARY.md |

## 测试覆盖率

```bash
make cover
```

生成的 HTML 报告将显示代码覆盖率情况。

## 开发流程

1. 修改代码
2. 运行 `make fmt vet` 验证
3. 运行 `make test` 执行测试
4. 运行 `make cover` 检查覆盖率

## 故障排除

### 测试超时
```bash
go test -timeout 30s ./...
```

### 跳过特定测试
```bash
go test -skip TestLevel3EditCycle ./integration
```

### 并行测试
```bash
go test -parallel 4 ./...
```
