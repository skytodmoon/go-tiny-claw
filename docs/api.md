# Go Tiny Claw - API 文档模板

## 工具调用

### 文件工具
- **读取文件**：`read_file(path: string)`
- **写入文件**：`write_file(path: string, content: string)`
- **编辑文件**：`edit_file(path: string, old_str: string, new_str: string)`

### 命令行工具
- **执行 Bash 命令**：`bash(command: string)`

## 核心模块

### 循环引擎
- `MainLoop()`

### 工具注册表
- `RegisterTool(name: string, tool: Tool)`

### 上下文管理
- `AddContext(key: string, value: string)`
- `GetContext(key: string)`