package okxInfo

import "github.com/EasyGolang/goTools/mOKX"

// 15 分钟线
type KD struct {
	mOKX.TypeKd // KDData 基础结构
	MA_26       string
	EMA_10      string
	EMA_60      string
}

/*

交易方向: 来自于 CoinMarket

开仓策略：
EMA_10 上穿 MA_26 或者 EMA_60 开多
EMA_10 下穿 MA_26 或者 EMA_60 开空

平仓策略
EMA_10 下穿 EMA_60 平多仓  (保险)
EMA_10 上穿 EMA_60 平空仓  (保险)

交易对象退出排行榜单 平仓


分两种交易策略。
1: 均线，挑选振幅比 居中币种
2: 裸K ，挑选振动幅比较大币种，看收盘价 角度 进行交易

*/
