#!/bin/bash

startName="CoinAI.net-{{.Port}}"
Path="{{.Path}}-{{.Port}}"

pm2 delete ${startName}
rm -rf ${Path}
