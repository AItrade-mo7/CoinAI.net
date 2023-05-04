package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/dbType"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func (_this *HunterObj) Analy() {
	if len(_this.TradeKdataList) < 100 {
		global.LogErr("hunter.Analy 数据长度错误", len(_this.TradeKdataList))
		return
	}
	// 控制持仓的最大数量 防止内存爆炸
	if len(_this.PositionArr)-_this.MaxLen > 0 {
		_this.PositionArr = _this.PositionArr[len(_this.PositionArr)-_this.MaxLen:]
	}

	// 更新持仓状态
	_this.CountPosition()
	_this.PositionArr = append(_this.PositionArr, _this.NowVirtualPosition)

	AnalyDir := GetAnalyDir(_this.NowVirtualPosition)

	// 打印日志和文件写入
	global.TradeLog.Println(_this.HunterName, "更新持仓状态", AnalyDir, mJson.ToStr(_this.NowVirtualPosition))
	mFile.Write(_this.OutPutDirectory+"/PositionArr.json", mJson.ToStr(_this.PositionArr))

	// 当前持仓方向不符合计算结果时，执行一次下单操作
	if _this.NowVirtualPosition.NowDir != AnalyDir {
		//  1, 16, 31, 46 每15分钟执行一次分析和换仓
		if IsAnalyTimeScale(mTime.GetUnixInt64()) {
			_this.OnOrder(AnalyDir)
		}
	}
}

func (_this *HunterObj) CountPosition() {
	NowTradeKdata := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	_this.NowVirtualPosition.InstID = NowTradeKdata.InstID
	_this.NowVirtualPosition.HunterName = _this.HunterName
	_this.NowVirtualPosition.NowTimeStr = NowTradeKdata.TimeStr
	_this.NowVirtualPosition.NowTime = mTime.GetTime().TimeUnix
	_this.NowVirtualPosition.NowC = NowTradeKdata.C
	_this.NowVirtualPosition.CAP_EMA = NowTradeKdata.CAP_EMA
	_this.NowVirtualPosition.EMA = NowTradeKdata.EMA
	_this.NowVirtualPosition.HunterConfig = NowTradeKdata.Opt

	if _this.NowVirtualPosition.NowDir != 0 { // 当前为持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowTradeKdata.C, _this.NowVirtualPosition.OpenAvgPx)
		if _this.NowVirtualPosition.NowDir < 0 { // 当前为持空仓状态则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		Level := _this.NowVirtualPosition.HunterConfig.MaxTradeLever
		_this.NowVirtualPosition.NowUplRatio = mCount.Mul(UplRatio, mStr.ToStr(Level)) // 乘以杠杆倍数

		// 表示当前正在持仓，持仓状态下收窄 边界值 80%
		_this.NowVirtualPosition.HunterConfig.CAP_Max = mCount.Mul(_this.NowVirtualPosition.HunterConfig.CAP_Max, "0.8")
		_this.NowVirtualPosition.HunterConfig.CAP_Min = mCount.Mul(_this.NowVirtualPosition.HunterConfig.CAP_Min, "0.8")
	}
}

func GetAnalyDir(NowPosition dbType.VirtualPositionType) int {
	CAP_EMA := NowPosition.CAP_EMA              // 当前CAP
	CAP_Max := NowPosition.HunterConfig.CAP_Max // Max 边界值
	CAP_Min := NowPosition.HunterConfig.CAP_Min // Min 边界值

	CountDir := 0
	if mCount.Le(CAP_EMA, CAP_Max) > 0 { // 大于 Max 则向上
		CountDir = 1
	}
	if mCount.Le(CAP_EMA, CAP_Min) < 0 { // 小于 Min 则向下
		CountDir = -1
	}

	return CountDir
}
