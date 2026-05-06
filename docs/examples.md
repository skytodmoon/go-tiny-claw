# Go Tiny Claw - 使用示例

## 快速开始

```bash
# 设置 API Key
export SILICONFLOW_API_KEY="your-api-key"

# 运行程序
go run cmd/claw/main.go
```

## 示例任务

### 1. 读取并写入文件

- 目标：读取 `README.md`，生成项目概述文档。
- 执行步骤：
  ```bash
  go run cmd/claw/main.go --task="read-and-summarize"
  ```

### 2. 多步执行

- 目标：创建多个文件并编辑内容。
- 执行步骤：
  ```bash
  go run cmd/claw/main.go --task="multi-step-file-edit"
  ```

### 3. 测试套件

```bash
# 运行内置测试任务
go run cmd/task-runner/main.go
```