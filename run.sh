#!/bin/bash

# 进入脚本所在目录
cd "$(dirname "$0")"

echo "正在安装依赖..."
go mod tidy

echo "正在启动应用..."
go run cmd/server/main.go