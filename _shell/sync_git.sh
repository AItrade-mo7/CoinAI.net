#!/bin/bash
# 加载变量
source "./_shell/init.sh"
#############

desc=$1

if [ -a ${desc} ]; then
  desc="exit-push"
fi

git pull &&
  git add . &&
  git commit -m "${desc}" &&
  git push &&
  echo "同步完成"
exit
