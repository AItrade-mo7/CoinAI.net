package config

import (
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mStr"
)

// 系统通知的邮箱
var SysEmail = "meichangliang@outlook.com"

var LeverOpt = []int{1, 2, 3, 4, 5, 6}

var AppEnv dbType.AppEnvType

var MainUser dbType.UserTable

var NoticeEmail = []string{}

// 计价的锚定货币
var Unit = "USDT"

var SPOT_suffix = mStr.Join("-", Unit)

var SWAP_suffix = mStr.Join(SPOT_suffix, "-SWAP")
