# Go Tiny Claw 使用示例

## 快速开始

### 1. 设置 API Key

选择以下任一 Provider 并设置对应的 API Key：

```bash
export SILICONFLOW_API_KEY="your-api-key"
export ZHIPU_API_KEY="your-api-key"
export NVIDIA_API_KEY="your-api-key"
```

### 2. 运行程序

使用以下命令启动 Go Tiny Claw：

```bash
go run cmd/claw/main.go
```

### 3. 运行示例任务

框架提供了 15 种不同的测试任务，可以通过以下命令运行测试套件：

```bash
go run cmd/task-runner/main.go
```

## 示例任务说明

以任务 1（文件读写基础）为例：
1. 读取 `README.md` 文件。
2. 创建项目说明文档并保存到 `docs/` 目录。

执行后，检查生成的文档是否包含预期的内容。