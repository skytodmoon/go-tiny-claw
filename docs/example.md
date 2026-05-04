# 使用示例

## 目标

读取 `README.md` 并生成项目概述文档。

## 步骤

1. **准备任务**: 定义任务目标。
   ```go
   task := "读取 README.md 并创建项目概述文档"
   ```

2. **调用工具**: 使用 `read_file` 工具读取文件。
   ```json
   {
     "tool": "read_file",
     "params": {
       "path": "README.md"
     }
   }
   ```

3. **处理内容**: 提取核心理念和功能列表。
   ```go
   summary := extractCoreInfo(fileContent)
   ```

4. **生成文档**: 使用 `write_file` 工具保存概述文档。
   ```json
   {
     "tool": "write_file",
     "params": {
       "path": "docs/overview.md",
       "content": "项目概述内容"
     }
   }
   ```

5. **验证结果**: 检查文件是否生成成功。
   ```bash
   cat docs/overview.md
   ```

## 完整示例代码

```go
package main

func main() {
    // 1. 定义任务
    task := "读取 README.md 并创建项目概述文档"

    // 2. 调用工具读取文件
    fileContent := callTool("read_file", map[string]string{"path": "README.md"})

    // 3. 提取关键信息
    summary := extractCoreInfo(fileContent)

    // 4. 生成文档
    callTool("write_file", map[string]string{
        "path":    "docs/overview.md",
        "content": summary,
    })

    // 5. 验证结果
    fmt.Println("任务完成！")
}
```