# Go Tiny Claw API 文档

## 工具注册表

框架提供了一套工具注册与调用机制，以下是内置工具的使用示例。

### 1. 文件读取工具 (ReadFile)

```json
{
  "tool": "read_file",
  "parameters": {
    "path": "README.md"
  }
}
```

### 2. 文件写入工具 (WriteFile)

```json
{
  "tool": "write_file",
  "parameters": {
    "path": "docs/example.txt",
    "content": "Hello, world!"
  }
}
```

### 3. 文件编辑工具 (EditFile)

```json
{
  "tool": "edit_file",
  "parameters": {
    "path": "docs/example.txt",
    "old_str": "Hello",
    "new_str": "Hi"
  }
}
```

### 4. Bash 命令工具 (Bash)

```json
{
  "tool": "bash",
  "parameters": {
    "command": "ls -la"
  }
}
```

## 工具扩展

如需扩展自定义工具，请参考 `internal/tools/registry.go` 实现工具接口。