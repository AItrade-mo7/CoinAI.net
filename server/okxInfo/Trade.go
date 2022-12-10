package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

var TradeInst mOKX.TypeTicker

var Ticking = make(chan string, 2)
