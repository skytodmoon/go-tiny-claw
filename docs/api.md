# API 参考文档

## 核心接口

### 工具调用

#### 描述
通过工具注册表调用内置工具（如文件读写、Bash 命令等）。

#### 请求示例

```json
{
  "tool": "read_file",
  "params": {
    "path": "README.md"
  }
}
```

#### 响应示例

```json
{
  "status": "success",
  "result": "文件内容"
}
```

### 四阶段循环

#### 描述
Agent 的工作循环：Thinking → Acting → Observation → Re-thinking。

#### 请求示例

```json
{
  "task": "读取文件并生成概述",
  "phases": ["Thinking", "Acting", "Observation", "Re-thinking"]
}
```

#### 响应示例

```json
{
  "status": "success",
  "summary": "任务完成总结"
}
```

---

## 工具列表

| 工具名称 | 描述 |
|----------|------|
| read_file | 读取文件内容 |
| write_file | 写入文件 |
| edit_file | 编辑文件内容 |
| bash | 执行 Bash 命令 |