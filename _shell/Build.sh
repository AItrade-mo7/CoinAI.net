#!/bin/bash
source "./_shell/init.sh"
#############

echo " =========== 正在进行编译 aarch64 =========== "
armName="${BuildName}_aarch64"

go mod tidy &&
  GOOS=linux GOARCH=arm64 GOARM=7 go build -o "${armName}"
echo "编译 完成"

echo " =========== 正在进行编译 x86_64 =========== "
amdName="${BuildName}_x86_64"

go mod tidy &&
  GOOS=linux GOARCH=amd64 go build -o "${amdName}"
echo "编译 完成"

echo " =========== 开始进行 文件整理 =========== "

echo "清理并创建 dist 目录"
rm -rf "${OutPutPath}"
mkdir "${OutPutPath}" &&
  echo "移动 goRun 文件"
mv "${armName}" "${OutPutPath}" &&
  mv "${amdName}" "${OutPutPath}" &&
  cp -r "${NowPath}""/README.md" "${OutPutPath}" &&
  cp -r "${NowPath}""/package.json" "${OutPutPath}" &&
  exit 0
