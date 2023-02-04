#!/bin/bash
##WebHook:~ 发布 ShellHooks.net ~
source "./_shell/init.sh"
# 加载变量

rm -rf "${OutPutPath}"

npm run build &&
  cd "${OutPutPath}" || exit

nowTime=$(date +%Y-%m-%d\T%H:%M:%S)

git init
git add .
git commit -m "${nowTime}"
git remote add origin "${DeployPath}"
git push -f --set-upstream origin master:main

echo "同步完成"
rm -rf "${OutPutPath}/.git"

exit 0
