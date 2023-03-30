package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var HunterTicking = make(chan string, 2) // 计算频率

var KdataInst mOKX.TypeInst // 这里一定为现货

var TradeInst mOKX.TypeInst // 这里一定为合约

var MaxLen = 900

var NowKdataList = []mOKX.TypeKd{}

var HLPerLeVel = 3 // 涨跌幅等级  按照涨跌幅度排名，  数字越大越不稳定 为 2 时 只剩下 BTC 和 ETH 俩货
