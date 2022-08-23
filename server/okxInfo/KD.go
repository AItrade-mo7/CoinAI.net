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

策略：
EMA_10 上穿 MA_26 开多
EMA_10 下穿 MA_26 开空

EMA_10 下穿 EMA_60 平多仓
EMA_10 上穿 EMA_60 平空仓

*/
