#!/bin/bash
startName="CoinAI.net-9453"
Path="/root/AITrade/CoinAI.net-9453"

echo ""===== 下载可执行文件 "====="

SystemType=$(arch)

downLoadPath="https://raw.githubusercontent.com/mo7static/CoinAI.net/main/CoinAI.net_x86_64"

if [[ ${SystemType} =~ "aarch64" ]]; then
  downLoadPath="https://raw.githubusercontent.com/mo7static/CoinAI.net/main/CoinAI.net_aarch64"
fi

cd ${Path}

pm2 delete ${startName}

rm -rf ${startName} &&
  wget -O ${startName} ${downLoadPath}

sudo chmod 777 ${startName}

echo "===== 启动服务 ====="

pm2 start ${startName} --name ${startName}