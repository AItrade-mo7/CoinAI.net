#!/bin/bash
# 加载变量
source "./_shell/init.sh"
#############

echo "清理目录"
rm -rf ./logs
rm -rf ./jsonData

echo "整理 mod"
go mod tidy

echo " ========== 开始运行 ========== "
go run main.go
