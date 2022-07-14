#!/bin/bash
source "./_shell/init.sh"
#############

echo " =========== 正在进行编译 aarch64 =========== "

go mod tidy &&
  GOOS=linux GOARCH=arm64 GOARM=7 go build -o ${buildName}"_aarch64"
echo "编译 完成"

echo " =========== 正在进行编译 x86_64 =========== "

set GOARCH=amd64
go mod tidy &&
  GOOS=linux GOARCH=amd64 go build -o ${buildName}"_x86_64"
echo "编译 完成"

echo " =========== 开始进行 文件整理 =========== "

echo "清理并创建 dist 目录"
rm -rf ${outPutPath}
mkdir ${outPutPath} &&
  echo "移动 goRun 文件"
mv ${buildName}"_aarch64" ${outPutPath} &&
  mv ${buildName}"_x86_64" ${outPutPath} &&
  cp -r ${path}"/README.md" ${outPutPath}
cp -r ${path}"/package.json" ${outPutPath}
