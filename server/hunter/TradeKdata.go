package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTalib"
	jsoniter "github.com/json-iterator/go"
)

func (_this *HunterObj) FormatTradeKdata() error {
	if len(_this.NowKdataList) < _this.TradeKdataOpt.EMA_Period+1 {
		err := fmt.Errorf(_this.HunterName, "hunter.FormatTradeKdata 数据不足")
		return err
	}

	if _this.TradeKdataOpt.EMA_Period == 0 ||
		_this.TradeKdataOpt.CAP_Period == 0 {
		err := fmt.Errorf(_this.HunterName, "hunter.FormatTradeKdata2 参数不正确 %+v", _this.TradeKdataOpt)
		return err
	}

	// 清理 TradeKdataList
	_this.TradeKdataList = []okxInfo.TradeKdType{}

	TradeKlineObj := NewTradeKdataObj(_this.TradeKdataOpt)

	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range _this.NowKdataList {
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := TradeKlineObj.NewTradeKdata(FormatEnd)
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
	Opt     okxInfo.TradeKdataOpt
}

func NewTradeKdataObj(opt okxInfo.TradeKdataOpt) *TradeKdataObj {
	obj := TradeKdataObj{}
	obj.EMA_Arr = []string{}
	obj.Opt = opt

	if obj.Opt.EMA_Period < 0 {
		obj.Opt.EMA_Period = 171
		global.LogErr("obj.Opt.MA_Period 参数为空，已设置为默认")
	}

	if obj.Opt.CAP_Period < 0 {
		obj.Opt.CAP_Period = 4
		global.LogErr("obj.Opt.CAP_Period 参数为空，已设置为默认")
	}

	return &obj
}

func (_this *TradeKdataObj) NewTradeKdata(KdataList []mOKX.TypeKd) (TradeKdata okxInfo.TradeKdType) {
	TradeKdata = okxInfo.TradeKdType{}
	jsonByte := mJson.ToJson(KdataList[len(KdataList)-1])
	jsoniter.Unmarshal(jsonByte, &TradeKdata)

	TradeKdata.Opt = _this.Opt // 在这里把设置打印出来

	// EMA
	TradeKdata.EMA = mTalib.ClistNew(mTalib.ClistOpt{
		KDList: KdataList,
		Period: _this.Opt.EMA_Period,
	}).EMA().ToStr()
	_this.EMA_Arr = append(_this.EMA_Arr, TradeKdata.EMA)

	// CAP_EMA
	TradeKdata.CAP_EMA = mTalib.ClistNew(mTalib.ClistOpt{
		CList:  _this.EMA_Arr,
		Period: _this.Opt.CAP_Period,
	}).CAP().ToStr()

	return
}
