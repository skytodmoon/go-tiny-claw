#!/bin/bash

# 批量重命名项目中的测试文件
# 用法示例: ./batch_rename.sh /path/to/project_name

if [ $# -ne 1 ]; then
    echo "Usage: $0 <项目路径>"
    exit 1
fi

PROJECT_DIR=$1

# 在项目中查找所有测试文件（假设测试文件以 _test.go 结尾）
find "$PROJECT_DIR" -name "*_test.go" | while read -r testfile; do
    BASE_NAME=$(basename "$testfile" "_test.go")
    DIR_NAME=$(dirname "$testfile")
    NEW_NAME="${DIR_NAME}/${BASE_NAME}_spec.go"

    # 检查新文件名是否已存在
    if [ -f "$NEW_NAME" ]; then
        echo "警告: 文件 $NEW_NAME 已存在，跳过重命名 $testfile"
    else
        mv "$testfile" "$NEW_NAME"
        echo "已重命名: $testfile -> $NEW_NAME"
    fi
done