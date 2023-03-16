#!/bin/bash
##WebHook:~ 发布 CoinAIPackage ~

## 设置并加载变量
source "./_shell/init.sh"
NowPath=${NowPath}
OutPutPath=${OutPutPath}
DeployPath=${DeployPath}
DeployLocalPath=${DeployLocalPath}

## 判断参数
desc=$1
if [ -z "${desc}" ]; then
  echo -e "\033[31m Err:需要发布说明 \033[0m"
  exit 1
fi
echo "git commit: ${desc}"

## 清理目录 dist 和 git 仓库
rm -rf "${OutPutPath}"
rm -rf "${DeployLocalPath}"

## 开始编译
npm run build &&
  echo "编译结束"

# 开始进行发布步骤 OutPutPath -> 远程 DeployPath
## 拉取远程仓库
git clone "${DeployPath}"

## git 本地仓库中 .git 移动到 dist 目录中
mv "${DeployLocalPath}/.git" "${OutPutPath}/.git"

## 删除本地 git 仓库
rm -rf "${DeployLocalPath}"

## dist 更名为 本地 git 仓库名
mv "${OutPutPath}" "${DeployLocalPath}"

# 将仓库推送到远端
cd "${DeployLocalPath}" || exit

git add . &&
  git commit -m "${desc}" &&
  git push &&
  ## 清理 dist
  rm -rf ".git" &&
  mv "${DeployLocalPath}" "${OutPutPath}" &&
  echo "已推送至${DeployPath}"

exit 0
