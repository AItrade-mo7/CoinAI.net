#!/bin/bash

## 设置并加载变量
source "./_shell/init.sh"
NowPath=${NowPath}
OutPutPath="${NowPath}/TaskRunner"
StartName="goRun-Task"
BuildName=${StartName}

echo "开始打包" &&
  go mod tidy &&
  go build -o "${BuildName}" task/task.go &&
  echo " =========== 开始进行文件整合 =========== "

pm2 delete "${StartName}"

rm -rf "${OutPutPath}"
mkdir "${OutPutPath}"

echo "移动 go build 文件"
mv "${BuildName}" "${OutPutPath}/" &&
  echo "启动 pm2 服务"

cd "${OutPutPath}" || exit

pm2 start "./${BuildName}" --name "${StartName}" --output "${OutPutPath}/out.log" --error "${OutPutPath}/err.log" --no-autorestart &&
  exit 0
