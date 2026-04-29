#!/bin/bash

# 进入脚本所在目录
cd "$(dirname "$0")"

echo "正在安装依赖..."
go mod tidy

echo "正在构建应用..."
go build -o school-app cmd/server/main.go

echo "构建完成！可执行文件：school-app"
