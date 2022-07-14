#!/bin/bash
# 加载变量
source "./_shell/init.sh"
#############

echo "更新依赖"
go mod tidy
## run
echo " ========== 开始运行 ========== "
go run main.go
