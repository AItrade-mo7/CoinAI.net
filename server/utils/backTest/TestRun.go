package backTest

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/hunter"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 任意时间点之间的回测函数

type TestOpt struct {
	StartTime int64
	EndTime   int64
	CcyName   string
}

type TestObj struct {
	StartTime int64
	EndTime   int64
	TableName string
	InstID    string
	KdataList []mOKX.TypeKd
}

func NewTest(opt TestOpt) *TestObj {
	obj := TestObj{}
	now := mTime.GetUnixInt64()

	if opt.EndTime < dbType.MinTime { // 如果太小，则变成当前
		opt.EndTime = now
	}

	if opt.StartTime < dbType.MinTime { // 如果太小，则自动变成过去一个月
		opt.StartTime = mTime.UnixTimeInt64.Day * 30
	}

	if opt.StartTime < dbType.DBKdataStart {
		opt.StartTime = dbType.DBKdataStart
	}

	if opt.StartTime > opt.EndTime {
		global.LogErr("backTest.NewTest 开始时间不可以大于结束时间")
	}

	obj.StartTime = opt.StartTime
	obj.EndTime = opt.EndTime

	obj.TableName = opt.CcyName + "USDT"
	obj.InstID = opt.CcyName + "-USDT"

	return &obj
}

func (_this *TestObj) GetDBKdata() *TestObj {
	total := (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Hour

	if total < 1 {
		return nil
	}

	Timeout := int(total) * 10
	if Timeout < 100 {
		Timeout = 100
	}

	db := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AITrade",
		Timeout:  Timeout,
	}).Connect().Collection(_this.TableName)
	defer global.RunLog.Println("关闭数据库连接", _this.TableName)
	defer db.Close()

	findOpt := options.Find()
	findOpt.SetSort(map[string]int{
		"TimeUnix": 1,
	})
	findOpt.SetAllowDiskUse(true)

	FK := bson.D{}
	FK = append(FK, bson.E{
		Key: "TimeUnix",
		Value: bson.D{
			{
				Key:   "$gte", // 大于或等于
				Value: _this.StartTime,
			}, {
				Key:   "$lte", // 小于或等于
				Value: _this.EndTime,
			},
		},
	})

	cur, err := db.Table.Find(db.Ctx, FK, findOpt)
	if err != nil {
		db.Close()
		return nil
	}

	_this.KdataList = []mOKX.TypeKd{}

	for cur.Next(db.Ctx) {
		var result mOKX.TypeKd
		cur.Decode(&result)
		_this.KdataList = append(_this.KdataList, result)
	}

	return _this
}

func (_this *TestObj) CheckKdataList() (resErr error) {
	resErr = nil

	if len(_this.KdataList) < 1 {
		resErr = fmt.Errorf("KdataList 长度不正确")
		return
	}

	for key, val := range _this.KdataList {
		preIndex := key - 1
		if preIndex < 0 {
			preIndex = 0
		}
		preItem := _this.KdataList[preIndex]
		nowItem := _this.KdataList[key]
		if key > 0 {
			if nowItem.TimeUnix-preItem.TimeUnix != mTime.UnixTimeInt64.Hour {
				resErr = fmt.Errorf("数据检查出错, %+v", nowItem.TimeUnix-preItem.TimeUnix)
				global.LogErr("数据检查出错 backTest.CheckKdataList", val.InstID, val.TimeStr, key)
				break
			}
		}
	}

	return
}

var (
	EMA_Arr        = []string{}
	MA_Arr         = []string{}
	TradeKdataList = []okxInfo.TradeKdType{}
	FormatEnd      = []mOKX.TypeKd{}
)

func (_this *TestObj) MockData() {
	// 在这里开始执行模拟数据流的流动
	// 执行一次清理
	TradeKdataList = []okxInfo.TradeKdType{}
	EMA_Arr = []string{}
	MA_Arr = []string{}
	FormatEnd = []mOKX.TypeKd{}

	for _, Kdata := range _this.KdataList {
		// 开始执行整理
		FormatEnd = append(FormatEnd, Kdata)
		TradeKdata := hunter.NewTradeKdata(Kdata, FormatEnd)
		TradeKdataList = append(TradeKdataList, TradeKdata)

		if len(TradeKdataList) >= 100 {
			// 开始执行分析
			Analy()
		}
	}
	WriteFilePath := config.Dir.JsonData + "/TestRun.json"
	mFile.Write(WriteFilePath, string(mJson.ToJson(TradeKdataList)))

	PrintResult()
}

type TypeOpen struct {
	Dir         int    // 开仓方向
	AvgPx       string // 开仓价格
	UplRatio    string // 未实现收益率
	OpenTimeStr string // 开仓时间
}

var (
	RSIMax  = okxInfo.TradeKdType{}
	NowOpen TypeOpen
	OpenArr []TypeOpen
)

// 打印结结果
func PrintResult() {
	NilNum := 0     // 空仓次数
	SellNum := 0    // 开空次数
	BuyNum := 0     // 开多次数
	AllNum := 0     // 总开仓次数
	Win := 0        // 盈利次数
	Lose := 0       // 亏损次数
	WinRatio := ""  // 盈利的比例
	LoseRatio := "" //  亏损的比例
	MaxWin := ""    //  最大盈利
	MaxLose := ""   //  最大亏损
	StartTime := OpenArr[0].OpenTimeStr
	EndTime := OpenArr[len(OpenArr)-1].OpenTimeStr

	Money := "1000"

	Lever := "1"

	Charge := mCount.Div("0.02", "100")
	ChargeAll := "0"

	for _, val := range OpenArr {

		if val.Dir == 0 {
			nowCharge := mCount.Mul(Money, Charge)
			Money = mCount.Sub(Money, nowCharge)
			Money = mCount.Cent(Money, 2)

			ChargeAll = mCount.Add(ChargeAll, nowCharge)
			NilNum++
			continue
		}

		if len(StartTime) < 1 {
			StartTime = val.OpenTimeStr
		}

		if len(val.OpenTimeStr) > 1 {
			EndTime = val.OpenTimeStr
		}

		if val.Dir > 0 {
			BuyNum++
		}
		if val.Dir < 0 {
			SellNum++
		}

		if mCount.Le(val.UplRatio, "0") > 0 {
			Win++
			WinRatio = mCount.Add(val.UplRatio, WinRatio)
		}

		if mCount.Le(val.UplRatio, "0") < 0 {
			Lose++
			LoseRatio = mCount.Add(val.UplRatio, LoseRatio)
		}

		if mCount.Le(val.UplRatio, MaxLose) < 0 {
			MaxLose = val.UplRatio
		}

		if mCount.Le(val.UplRatio, MaxWin) > 0 {
			MaxWin = val.UplRatio
		}

		Upl := mCount.Div(val.UplRatio, "100")
		LeverUpl := mCount.Mul(Upl, Lever)

		nowMoney := mCount.Mul(Money, LeverUpl)
		Money = mCount.Add(Money, nowMoney)

		nowCharge := mCount.Mul(Money, Charge)
		Money = mCount.Sub(Money, nowCharge)

		Money = mCount.Cent(Money, 2)

		ChargeAll = mCount.Add(ChargeAll, nowCharge)

		AllNum++
	}

	fmt.Printf(
		`空仓次数: %+v;
开空次数：%+v;
开多次数：%+v;
总开仓次数: %+v;
盈利次数: %+v;
亏损次数: %+v;
总盈利比例: %+v;
总亏损的比例: %+v;
开始时间: %+v;
结束时间: %+v;
最大单次盈利: %+v;
最大单次亏损: %+v;
1000 扣除手续费后结余: %+v;
总手续费: %+v;
`, NilNum, SellNum, BuyNum, AllNum, Win, Lose, WinRatio, LoseRatio,
		StartTime, EndTime,
		MaxWin, MaxLose, Money, ChargeAll,
	)

	mFile.Write(config.Dir.JsonData+"/Open.json", string(mJson.ToJson(OpenArr)))
}

func Analy() {
	// Pre := TradeKdataList[len(TradeKdataList)-2]
	Now := TradeKdataList[len(TradeKdataList)-1]

	// preIdx := hunter.CAPIdxToText(Pre.CAPIdx)
	// nowIdx := hunter.CAPIdxToText(Now.CAPIdx)

	PreList5 := TradeKdataList[len(TradeKdataList)-6:]
	RsiRegion_Down := hunter.Is_RsiRegion_GoDown(PreList5)
	RsiRegion_Up := hunter.Is_RsiRegion_GoUp(PreList5)
	RsiRegion_Gte2 := hunter.Is_RsiRegion_Gte2(PreList5)

	Open := 0
	// 副调 RSI 超买超卖
	// if len(RSIMax.RSI_18) > 1 {
	// 	if mCount.Le(RSIMax.RSI_18, "35") < 0 && Now.CAPIdx > 0 {
	// 		Open = 1
	// 	}
	// 	if mCount.Le(RSIMax.RSI_18, "65") > 0 && Now.CAPIdx < 0 {
	// 		Open = -1
	// 	}
	// }

	// 主调  Last.CAPIdx
	// if nowIdx != preIdx {
	if Now.CAPIdx > 0 { // Buy
		if len(RsiRegion_Up) > 1 {
			if RsiRegion_Gte2 {
				Open = 1
			}
		}
	}

	if Now.CAPIdx < 0 { // sell
		if len(RsiRegion_Down) > 1 {
			if RsiRegion_Gte2 {
				Open = -1
			}
		}
	}
	// }

	PrintLnResult := func() {
		global.TradeLog.Printf(
			"%v %6v RSI:%2v %8v CAP_EMA:%7v Upl:%10v Gte2:%v RsiDown:%+v RsiUp:%+v \n",
			Now.TimeStr, fmt.Sprint(Open)+hunter.CAPIdxToText(Open)+fmt.Sprint(Now.CAPIdx),
			Now.RsiRegion, Now.RSI_18,
			Now.CAP_EMA,
			NowOpen.UplRatio+","+fmt.Sprint(NowOpen.Dir),
			RsiRegion_Gte2,
			RsiRegion_Down,
			RsiRegion_Up,
		)

		if Open != NowOpen.Dir {
			OpenArr = append(OpenArr, NowOpen) // 记录平仓收益

			// 开仓
			NowOpen.Dir = Open
			NowOpen.AvgPx = Now.C
			NowOpen.UplRatio = ""
			NowOpen.OpenTimeStr = Now.TimeStr

		}
	}

	// 计算收益
	if len(NowOpen.AvgPx) > 0 {
		NowOpen.UplRatio = mCount.RoseCent(Now.C, NowOpen.AvgPx)
		if NowOpen.Dir < 0 {
			NowOpen.UplRatio = mCount.Sub("0", NowOpen.UplRatio)
		}
	}

	if Open > 0 { // buy
		PrintLnResult()
		return
	}
	if Open < 0 { // sell
		PrintLnResult()
		return
	}

	global.TradeLog.Printf(
		"%v %6v RSI:%2v %8v CAP_EMA:%7v Upl:%10v RsiDown:%+v RsiUp:%+v     \n",
		Now.TimeStr, fmt.Sprint(Now.CAPIdx),
		Now.RsiRegion, Now.RSI_18,
		Now.CAP_EMA,
		NowOpen.UplRatio+","+fmt.Sprint(NowOpen.Dir),
		RsiRegion_Down,
		RsiRegion_Up,
	)
}
