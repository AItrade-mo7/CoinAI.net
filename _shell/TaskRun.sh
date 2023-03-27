#!/bin/bash

## 设置并加载变量
source "./_shell/init.sh"
OutPutPath=${OutPutPath}
StartName="Task-Run"
BuildName=${StartName}
DeployPath="/root/ProdProject/"${StartName}

echo "清理目录"
rm -rf ./logs
rm -rf ./jsonData

echo "整理 mod"
go mod tidy

echo " ========== 开始运行 ========== "
go run task/task.go
