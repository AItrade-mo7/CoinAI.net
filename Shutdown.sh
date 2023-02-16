#!/bin/bash
startFilePath="/root/AItrade/CoinAI.net/jsonData/Shutdown.txt"
sudo cat >${startFilePath} <<END
#!/bin/bash
time3=$(date "+%Y-%m-%d %H:%M:%S")

END
