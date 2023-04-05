package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTalib"
	jsoniter "github.com/json-iterator/go"
)

func (_this *HunterObj) FormatTradeKdata() error {
	if len(_this.NowKdataList) < _this.TradeKdataOpt.MA_Period {
		err := fmt.Errorf(_this.HunterName, "hunter.FormatTradeKdata 数据不足")
		return err
	}

	if _this.TradeKdataOpt.MA_Period == 0 ||
		_this.TradeKdataOpt.RSI_Period == 0 ||
		_this.TradeKdataOpt.RSI_EMA_Period == 0 ||
		_this.TradeKdataOpt.CAP_Period == 0 {
		err := fmt.Errorf(_this.HunterName, "hunter.FormatTradeKdata2 参数不正确 %+v", _this.TradeKdataOpt)
		return err
	}

	// 清理 TradeKdataList
	_this.TradeKdataList = []TradeKdType{}

	TradeObj := NewTradeKdataObj(_this.TradeKdataOpt)

	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range _this.NowKdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := TradeObj.NewTradeKdata(FormatEnd)
		_this.TradeKdataList = append(_this.TradeKdataList, TradeKdata)
	}

	WriteFilePath := config.Dir.JsonData + "/" + _this.HunterName + "_TradeKdataList.json"
	mFile.Write(WriteFilePath, string(mJson.ToJson(_this.TradeKdataList)))
	global.TradeLog.Println("数据整理完毕,已写入", len(_this.TradeKdataList), WriteFilePath)
	return nil
}

// ============================================================================

type TradeKdataObj struct {
	EMA_Arr []string
	MA_Arr  []string
	RSI_Arr []string
	Opt     TradeKdataOpt
}

type TradeKdType struct {
	mOKX.TypeKd
	EMA          string // EMA 值
	MA           string // MA 值
	RSI          string // RSI 的值
	RSI_EMA      string // RSI 的值
	CAP_EMA      string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	CAP_MA       string // 基于 EMA 的 平滑点数 0-100 的浮点类型
	CAPIdx       int
	RsiEmaRegion int // 整型 Rsi 的震荡区域  -3 -2 -1 0 1 2 3
	Opt          TradeKdataOpt
}

func NewTradeKdataObj(opt TradeKdataOpt) *TradeKdataObj {
	obj := TradeKdataObj{}
	obj.EMA_Arr = []string{}
	obj.MA_Arr = []string{}
	obj.RSI_Arr = []string{}
	obj.Opt = opt

	return &obj
}

func (_this *TradeKdataObj) NewTradeKdata(TradeKdataList []mOKX.TypeKd) (TradeKdata TradeKdType) {
	TradeKdata = TradeKdType{}
	jsonByte := mJson.ToJson(TradeKdataList[len(TradeKdataList)-1])
	jsoniter.Unmarshal(jsonByte, &TradeKdata)

	TradeKdata.Opt = _this.Opt // 在这里把设置打印出来

	// EMA
	TradeKdata.EMA = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: _this.Opt.MA_Period,
	}).EMA().ToStr()
	_this.EMA_Arr = append(_this.EMA_Arr, TradeKdata.EMA)

	// MA
	TradeKdata.MA = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: _this.Opt.MA_Period,
	}).MA().ToStr()
	_this.MA_Arr = append(_this.MA_Arr, TradeKdata.MA)

	// RSI
	TradeKdata.RSI = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: TradeKdataList,
		Period: _this.Opt.RSI_Period,
	}).RSI().ToStr()
	_this.RSI_Arr = append(_this.RSI_Arr, TradeKdata.RSI)

	// RSI_EMA
	TradeKdata.RSI_EMA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  _this.RSI_Arr,
		Period: _this.Opt.RSI_EMA_Period,
	}).EMA().ToStr()

	// CAP_EMA
	TradeKdata.CAP_EMA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  _this.EMA_Arr,
		Period: _this.Opt.CAP_Period,
	}).CAP().ToStr()
	// CAP_MA
	TradeKdata.CAP_MA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  _this.MA_Arr,
		Period: _this.Opt.CAP_Period,
	}).CAP().ToStr()

	// CAPIdx 计算
	TradeKdata.CAPIdx = GetCAPIdx(TradeKdata)

	// 区域计算
	TradeKdata.RsiEmaRegion = GetRsiRegion(TradeKdata)

	return
}
