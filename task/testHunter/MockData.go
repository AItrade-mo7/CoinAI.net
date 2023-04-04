package testHunter

/*

// 模拟数据流动并执行分析交易
func (_this *TestObj) MockData(MockOpt BillingType, TradeKdataOpt hunter.TradeKdataOpt) {
	// 收益结算
	Billing = BillingType{}
	Billing.MockName = MockOpt.MockName
	Billing.InitMoney = MockOpt.InitMoney // 设定初始资金
	Billing.Money = MockOpt.InitMoney     // 设定当前账户资金
	Billing.Level = MockOpt.Level
	Billing.Charge = MockOpt.Charge
	Billing.InstID = _this.KdataList[0].InstID
	Billing.Days = (_this.EndTime - _this.StartTime) / mTime.UnixTimeInt64.Day

	Billing.MinMoney.Value = MockOpt.InitMoney
	Billing.MaxMoney.Value = MockOpt.InitMoney
	// 交易信息清空
	NowPosition = PositionType{}   // 清空持仓
	PositionArr = []PositionType{} // 清空持仓数组
	OrderArr = []OrderType{}       // 下单数组清空

	global.Run.Println("新建Mock数据",
		mJson.Format(map[string]any{
			"参数组名称":   Billing.MockName,
			"初始资金":    Billing.InitMoney,
			"杠杆倍率":    Billing.Level,
			"手续费率(%)": Billing.Charge,
		}),
		mJson.Format(TradeKdataOpt),
	)

	global.TradeLog.Println(" ============== 开始分析和交易 ============== ", Billing.MockName)

	// 清理 TradeKdataList
	hunter.TradeKdataList = []hunter.TradeKdType{}
	hunter.EMA_Arr = []string{}
	hunter.MA_Arr = []string{}
	hunter.RSI_Arr = []string{}

	FormatEnd := []mOKX.TypeKd{}
	for _, Kdata := range _this.KdataList {
		FormatEnd = append(FormatEnd, Kdata)

		if len(FormatEnd) < TradeKdataOpt.MA_Period {
			continue
		}

		TradeKdata := hunter.NewTradeKdata(FormatEnd, TradeKdataOpt)
		hunter.TradeKdataList = append(hunter.TradeKdataList, TradeKdata)

		// 开始执行分析交易
		Analy()

		if mCount.Le(NowPosition.UplRatio, "-45") < 0 {
			global.Log.Println("爆仓！", Billing.MockName, Kdata.TimeStr)
			break
		}
	}

	if len(hunter.TradeKdataList) > 0 {
		Billing.EndTime = hunter.TradeKdataList[len(hunter.TradeKdataList)-1].TimeStr
	}

	// 搜集和整理结果
	ResultCollect()
}

func Analy() {
	NowKdata := hunter.TradeKdataList[len(hunter.TradeKdataList)-1]
	AnalyDir := 0 // 分析的方向，默认为 0 不开仓

	if mCount.Le(NowKdata.CAP_EMA, "0") > 0 { // 大于 0 则开多
		AnalyDir = 1
	}

	if mCount.Le(NowKdata.CAP_EMA, "0") < 0 { // 小于 0 则开空
		AnalyDir = -1
	}

	// 更新持仓状态
	if NowPosition.Dir != 0 { // 当前为持持仓状态，则计算收益率
		UplRatio := mCount.RoseCent(NowKdata.C, NowPosition.OpenAvgPx)
		if NowPosition.Dir < 0 { // 当前为持 空 仓状态 则翻转该值
			UplRatio = mCount.Sub("0", UplRatio)
		}
		NowPosition.UplRatio = mCount.Mul(UplRatio, Billing.Level) // 乘以杠杆倍数
	}
	NowPosition.NowTimeStr = NowKdata.TimeStr
	NowPosition.NowC = NowKdata.C
	PositionArr = append(PositionArr, NowPosition)

	if mCount.Le(NowPosition.UplRatio, Billing.PositionMinRatio.Value) < 0 {
		Billing.PositionMinRatio.Value = NowPosition.UplRatio
		Billing.PositionMinRatio.TimeStr = NowPosition.NowTimeStr
	}
	if mCount.Le(NowPosition.UplRatio, Billing.PositionMaxRatio.Value) > 0 {
		Billing.PositionMaxRatio.Value = NowPosition.UplRatio
		Billing.PositionMaxRatio.TimeStr = NowPosition.NowTimeStr
	}

	// 当前持仓与 判断方向不符合时，执行一次下单操作
	if NowPosition.Dir != AnalyDir {
		OnOrder(AnalyDir, NowKdata)
	}

	global.KdataLog.Println(Billing.MockName, NowKdata.TimeStr, AnalyDir)
}

// 根据下单结果进行模拟持仓
func BillingFun(NowKdata hunter.TradeKdType) {
	fmt.Println(Billing.MockName, "下单总结一次",
		NowKdata.TimeStr,
		"持仓方向", NowPosition.Dir,
		"收益率", NowPosition.UplRatio,
	)

	if NowPosition.Dir == 0 {
		Billing.NilNum++ // 空仓计数
	} else {
		if len(Billing.StartTime) == 0 {
			Billing.StartTime = NowPosition.OpenTimeStr
		}
	}

	if NowPosition.Dir < 0 {
		Billing.SellNum++ // 开空 计数
	}
	if NowPosition.Dir > 0 {
		Billing.BuyNum++ // 开多 计数
	}

	if mCount.Le(NowPosition.UplRatio, "0") > 0 {
		Billing.Win++                                                         // 盈利次数计数
		Billing.WinRatio = mCount.Add(NowPosition.UplRatio, Billing.WinRatio) // 盈利比例相加
	}

	if mCount.Le(NowPosition.UplRatio, "0") < 0 {
		Billing.Lose++                                                          // 亏损次数计数
		Billing.LoseRatio = mCount.Add(NowPosition.UplRatio, Billing.LoseRatio) // 盈利比例相加
	}
	// 单次最大亏损和单次最大盈利
	if mCount.Le(NowPosition.UplRatio, Billing.MaxRatio.Value) > 0 {
		Billing.MaxRatio.Value = NowPosition.UplRatio
		Billing.MaxRatio.TimeStr = NowPosition.NowTimeStr
	}
	if mCount.Le(NowPosition.UplRatio, Billing.MinRatio.Value) < 0 {
		Billing.MinRatio.Value = NowPosition.UplRatio
		Billing.MinRatio.TimeStr = NowPosition.NowTimeStr
	}

	Upl := mCount.Div(NowPosition.UplRatio, "100") // 格式化收益率
	ChargeUpl := mCount.Div(Billing.Charge, "100") // 格式化手续费率

	makeMoney := mCount.Mul(Billing.Money, Upl)          // 当前盈利的金钱
	Billing.Money = mCount.Add(Billing.Money, makeMoney) // 相加得出当账户总资金量

	nowCharge := mCount.Mul(Billing.Money, ChargeUpl)            // 当前产生的手续费
	Billing.Money = mCount.Sub(Billing.Money, nowCharge)         // 减去手续费
	Billing.ChargeAll = mCount.Add(Billing.ChargeAll, nowCharge) // 记录一下手续费

	Billing.Money = mCount.CentRound(Billing.Money, 3)         // 四舍五入保留两位小数
	Billing.ChargeAll = mCount.CentRound(Billing.ChargeAll, 3) // 四舍五入保留两位小数

	if mCount.Le(Billing.Money, Billing.MinMoney.Value) < 0 {
		Billing.MinMoney.Value = Billing.Money
		Billing.MinMoney.TimeStr = NowPosition.NowTimeStr
	}

	if mCount.Le(Billing.Money, Billing.MaxMoney.Value) > 0 {
		Billing.MaxMoney.Value = Billing.Money
		Billing.MaxMoney.TimeStr = NowPosition.NowTimeStr
	}

	Billing.AllNum++ // 记录一下总交易次数
}

// 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func OnOrder(dir int, NowKdata hunter.TradeKdType) {
	BillingFun(NowKdata) // 下单之前 计算一次收益

	if dir > 0 { // 开多
		// 下订单
		OrderArr = append(OrderArr, OrderType{
			Type:    "Buy",            // 下多单
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新持仓状态
		NowPosition = PositionType{
			Dir:         1,          // 持仓多方向
			OpenAvgPx:   NowKdata.C, // 持仓价格
			NowTimeStr:  NowKdata.TimeStr,
			UplRatio:    "0", // 当前收益率
			NowC:        NowKdata.C,
			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
			InstID:      NowKdata.InstID,  // 开仓币种
		}
	}

	if dir < 0 { // 开空
		// 下订单
		OrderArr = append(OrderArr, OrderType{
			Type:    "Sell",           // 下空单
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新持仓状态
		NowPosition = PositionType{
			Dir:         -1,         // 持仓空方向
			OpenAvgPx:   NowKdata.C, // 持仓价格
			NowTimeStr:  NowKdata.TimeStr,
			UplRatio:    "0", // 当前收益率
			NowC:        NowKdata.C,
			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
			InstID:      NowKdata.InstID,  // 开仓币种
		}
	}

	if dir == 0 { // 平仓
		// 下订单
		OrderArr = append(OrderArr, OrderType{
			Type:    "Close",          // 平仓
			InstID:  NowKdata.InstID,  // 下单币种
			AvgPx:   NowKdata.C,       // 记录下单价格
			TimeStr: NowKdata.TimeStr, // 记录下单时间
		})
		// 更新为空仓状态
		NowPosition = PositionType{
			Dir:         0,  // 持仓空方向
			OpenAvgPx:   "", // 持仓价格
			NowTimeStr:  NowKdata.TimeStr,
			UplRatio:    "0", // 当前收益率
			NowC:        NowKdata.C,
			OpenTimeStr: NowKdata.TimeStr, // 开仓时间
			InstID:      NowKdata.InstID,  // 开仓币种
		}
	}
}

func ResultCollect() {
	// 记录 整理好的数组
	TradeKdataList_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-TradeKdataList.json")
	mFile.Write(TradeKdataList_Path, string(mJson.ToJson(hunter.TradeKdataList)))
	global.Run.Println("TradeKdataList: ", TradeKdataList_Path)

	// 记录 持仓数组
	PositionArr_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-PositionArr.json")
	mFile.Write(PositionArr_Path, string(mJson.ToJson(PositionArr)))
	global.Run.Println("PositionArr: ", PositionArr_Path)

	// 记录 下单数组
	OrderArr_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-OrderArr.json")
	mFile.Write(OrderArr_Path, string(mJson.ToJson(OrderArr)))
	global.Run.Println("OrderArr: ", OrderArr_Path)

	// 记录 交易结果
	Billing_Path := mStr.Join(config.Dir.JsonData, "/", Billing.MockName, "-Billing.json")
	mFile.Write(Billing_Path, string(mJson.ToJson(Billing)))
	global.Run.Println("Billing: ", Billing_Path)

	Tmp := `交易结果:
InstID: ${InstID}
第一次持仓时间: ${StartTime}
数据结束时间: ${EndTime}
总天数: ${Days}
空仓次数: ${NilNum}
开空次数: ${SellNum}
开多次数: ${BuyNum}
总开仓次数: ${AllNum}
盈利次数: ${Win}
总盈利比率: ${WinRatio}
亏损次数: ${Lose}
总亏损比率: ${LoseRatio}
平仓后单笔最大盈利比率: ${MaxRatio}
平仓后单笔最小盈利比率: ${MinRatio}
手续费率: ${Charge}
总手续费: ${ChargeAll}
参数名称: ${MockName}
初始金钱: ${InitMoney}
账户当前余额: ${Money}
平仓后历史最低余额: ${MinMoney}
平仓后历史最高余额: ${MaxMoney}
持仓过程中最低盈利比率: ${PositionMinRatio}
持仓过程中最高盈利比率: ${PositionMaxRatio}
杠杆倍数: ${Level}
`
	Data := map[string]string{
		"InstID":           Billing.InstID,
		"StartTime":        Billing.StartTime, // 开始时间
		"EndTime":          Billing.EndTime,   // 结束时间
		"Days":             mStr.ToStr(Billing.Days),
		"NilNum":           mStr.ToStr(Billing.NilNum),
		"SellNum":          mStr.ToStr(Billing.SellNum),
		"BuyNum":           mStr.ToStr(Billing.BuyNum),
		"AllNum":           mStr.ToStr(Billing.AllNum),
		"Win":              mStr.ToStr(Billing.Win),
		"WinRatio":         Billing.WinRatio,
		"Lose":             mStr.ToStr(Billing.Lose),
		"LoseRatio":        Billing.LoseRatio,
		"MaxRatio":         mJson.ToStr(Billing.MaxRatio),
		"MinRatio":         mJson.ToStr(Billing.MinRatio),
		"Charge":           Billing.Charge,
		"ChargeAll":        Billing.ChargeAll,
		"MockName":         Billing.MockName,
		"InitMoney":        Billing.InitMoney,
		"Money":            Billing.Money,
		"MinMoney":         mJson.ToStr(Billing.MinMoney),
		"MaxMoney":         mJson.ToStr(Billing.MaxMoney),
		"PositionMinRatio": mJson.ToStr(Billing.PositionMinRatio),
		"PositionMaxRatio": mJson.ToStr(Billing.PositionMaxRatio),
		"Level":            Billing.Level,
	}

	global.TradeLog.Println(mStr.Temp(Tmp, Data))
}

*/
