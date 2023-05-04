package hunter

import (
	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/global/dbType"
	"CoinAI.net/server/okxApi"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mCount"
	"github.com/EasyGolang/goTools/mEncrypt"
	"github.com/EasyGolang/goTools/mFile"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mStruct"
	"github.com/EasyGolang/goTools/mTime"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// // 下单  参数：dir 下单方向 NowKdata : 当前市场行情
func (_this *HunterObj) OnOrder(dir int) {
	NowKTradeData := _this.TradeKdataList[len(_this.TradeKdataList)-1]

	// 结算本期持仓
	_this.BillingFun()

	// 同步持仓状态, 相当于下单了
	if dir > 0 {
		// 开多
		_this.NowVirtualPosition.NowDir = 1
	}
	if dir < 0 {
		// 开空
		_this.NowVirtualPosition.NowDir = -1
	}
	// 同步下单价格
	_this.NowVirtualPosition.OpenAvgPx = NowKTradeData.C
	_this.NowVirtualPosition.OpenTimeStr = NowKTradeData.TimeStr
	_this.NowVirtualPosition.OpenTime = mTime.GetTime().TimeUnix

	// 同步平仓状态
	if dir == 0 {
		_this.NowVirtualPosition.NowDir = 0
		_this.NowVirtualPosition.OpenAvgPx = ""
		_this.NowVirtualPosition.OpenTimeStr = ""
		_this.NowVirtualPosition.OpenTime = 0
	}
	// 平仓后未实现盈亏重置为 0
	_this.NowVirtualPosition.NowUplRatio = "0"

	// 在这里执行下单
	_this.OrderOpen() // 这里的结果 要么是 1 要么是 0  要么 是 -1 没有第三种了
	global.TradeLog.Println(_this.HunterName, "下单一次", mJson.ToStr(_this.NowVirtualPosition))
	_this.OrderArr = append(_this.OrderArr, _this.NowVirtualPosition)
	mFile.Write(_this.OutPutDirectory+"/OrderArr.json", mJson.ToStr(_this.OrderArr))

	_this.SyncAllApiKey()
}

func (_this *HunterObj) BillingFun() {
	// 这里是下单之前的结算周期
	Upl := mCount.Div(_this.NowVirtualPosition.NowUplRatio, "100")     // 格式化收益率
	ChargeUpl := mCount.Div(_this.NowVirtualPosition.ChargeUpl, "100") // 格式化手续费率

	Money := _this.NowVirtualPosition.Money // 提取 Money
	makeMoney := mCount.Mul(Money, Upl)     // 当前盈利的金钱
	Money = mCount.Add(Money, makeMoney)    // 相加得出当账户剩余资金

	nowCharge := mCount.Mul(Money, ChargeUpl) // 当前产生的手续费
	Money = mCount.Sub(Money, nowCharge)      // 减去手续费
	Money = mCount.CentRound(Money, 3)        // 四舍五入保留三位小数
	_this.NowVirtualPosition.Money = Money    // 保存结果到当前持仓
	global.Run.Println("结算一次", mJson.ToStr(_this.NowVirtualPosition))
}

func (_this *HunterObj) OrderOpen() {
	// 在这里进行 下单存储。
	global.Run.Println("下单", mJson.ToStr(_this.NowVirtualPosition))
	if _this.NowVirtualPosition.NowDir > 0 {
		_this.SetOrderDB("Buy")
	}
	if _this.NowVirtualPosition.NowDir < 0 {
		_this.SetOrderDB("Sell")
	}

	if _this.NowVirtualPosition.NowDir == 0 {
		_this.SetOrderDB("Close")
	}
}

func (_this *HunterObj) SetOrderDB(Type string) {
	var orderData dbType.CoinOrderTable
	jsoniter.Unmarshal(mJson.ToJson(_this.NowVirtualPosition), &orderData)
	orderData.CreateTime = mTime.GetUnixInt64()
	orderData.Type = Type
	orderData.ServeID = config.AppEnv.ServeID
	orderData.TimeID = mOKX.GetTimeID(orderData.NowTime)
	orderData.OrderID = mEncrypt.GetUUID()

	db, err := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "AIServe",
	}).Connect()
	if err != nil {
		global.LogErr("hunter.SetOrderDB 数据库连接失败", _this.HunterName, err)
		return
	}
	defer db.Close()
	db.Collection("CoinOrder")

	_, err = db.Table.InsertOne(db.Ctx, orderData)
	if err != nil {
		global.LogErr("hunter.SetOrderDB 数据存储失败", _this.HunterName, err)
	}
}

type ErrObj struct {
	Name string
	Err  string //
}

type SettlementType struct {
	OkxPositions    dbType.PositionsData
	OKXBalance      []dbType.AccountBalance
	OkxKey          dbType.OkxKeyType
	VirtualPosition dbType.VirtualPositionType // 当前的虚拟持仓 数据库 OrderArr 最后一位
}

func (_this *HunterObj) SyncAllApiKey() {
	DirText := "保持空仓"
	if _this.NowVirtualPosition.NowDir > 0 {
		DirText = "买多看涨"
	}
	if _this.NowVirtualPosition.NowDir < 0 {
		DirText = "买空看跌"
	}

	global.TradeLog.Println(_this.HunterName, "开始执行所有的ApiKey")
	ApiKeyList := []dbType.OkxKeyType{}

	for _, item := range config.AppEnv.ApiKeyList {
		if item.Hunter == _this.HunterName {
			ApiKeyList = append(ApiKeyList, item)

			// 发送邮件通知
			tmplStr := `
		<br />
		策略名称： ${HunterName}  <br />
		当前持仓建议： ${NowDir}  <br />
		ApiKey Name: ${Name} <br />
		`
			lMap := map[string]string{
				"HunterName": _this.HunterName,
				"NowDir":     DirText,
				"Name":       item.Name,
			}

			Content := mStr.Temp(tmplStr, lMap)
			taskPush.SysEmail(taskPush.SysEmailOpt{
				From:        config.SysName,
				To:          []string{item.UserID},
				Subject:     "市场方向已改变",
				Title:       "市场方向已改变,系统将在30秒内同步您的持仓",
				Content:     Content,
				Description: "单个用户同步持仓邮件",
			})

		}
	}

	if len(ApiKeyList) < 1 {
		return
	}

	RightAccount := []dbType.OkxKeyType{}

	AccountSettlement := []SettlementType{} // 用于存储当前用户的平仓订单

	var ErrList []ErrObj
	for _, OkxKey := range ApiKeyList {

		// 新建账户对象
		OKXAccount, err := okxApi.NewAccount(okxApi.AccountParam{
			OkxKey: OkxKey,
		})
		if err != nil {
			ErrList = append(ErrList, ErrObj{
				Err:  mStr.Join("创建用户对象失败:", err),
				Name: OkxKey.Name,
			})
			continue
		}
		// 获取 Hunter
		err = OKXAccount.GetHunter()
		if err != nil {
			ErrList = append(ErrList, ErrObj{
				Err:  mStr.Join("获取Hunter失败:", err),
				Name: OkxKey.Name,
			})
			continue
		}
		// 读取当前持仓
		err = OKXAccount.GetPositions()
		if err != nil {
			ErrList = append(ErrList, ErrObj{
				Err:  mStr.Join("读取持仓失败:", err),
				Name: OkxKey.Name,
			})
			continue
		}

		// 存储当前 Hunter 持仓
		var NowAccountPos struct {
			Dir          int
			InstID       string
			OkxPositions dbType.PositionsData
		}
		for _, Positions := range OKXAccount.Positions {
			if OKXAccount.NowHunter.TradeInst.InstID == Positions.InstID {
				NowAccountPos.InstID = Positions.InstID
				NowAccountPos.Dir = mCount.Le(Positions.Pos, "0")
				NowAccountPos.OkxPositions = Positions
			}
		}

		// 判断持仓是否一致, 一致则无需操作
		if NowAccountPos.Dir == _this.NowVirtualPosition.NowDir {
			ErrList = append(ErrList, ErrObj{
				Err:  "当前账户持仓已经与策略保持一致,无需下单。",
				Name: OkxKey.Name,
			})
			continue
		}

		// 执行平仓操作
		err = OKXAccount.Close()
		if err != nil {
			ErrList = append(ErrList, ErrObj{
				Err:  mStr.Join("平仓失败:", err),
				Name: OkxKey.Name,
			})
			continue
		}

		// 此时读取账户余额
		err = OKXAccount.GetBalance()
		if err != nil {
			ErrList = append(ErrList, ErrObj{
				Err:  mStr.Join("读取余额失败:", err),
				Name: OkxKey.Name,
			})
			continue
		}

		// 平仓成功后，记录 okx 的持仓、账户余额、okxKey、虚拟持仓,也就是开仓。
		AccountSettlement = append(AccountSettlement, SettlementType{
			OkxPositions:    NowAccountPos.OkxPositions,
			OKXBalance:      OKXAccount.Balance,
			OkxKey:          OkxKey,
			VirtualPosition: _this.NowVirtualPosition,
		})

		// 根据情况开仓
		if _this.NowVirtualPosition.NowDir > 0 {
			err = OKXAccount.Buy()
		}
		if _this.NowVirtualPosition.NowDir < 0 {
			err = OKXAccount.Sell()
		}

		if err != nil {
			ErrList = append(ErrList, ErrObj{
				Err:  mStr.Join("交易所下单失败:", err),
				Name: OkxKey.Name,
			})
			continue
		}
		// 记录走完流程的账户
		RightAccount = append(RightAccount, OkxKey)
	}

	tmplStr := `
<br />
策略名称： ${HunterName}  <br />
当前持仓建议： ${NowDir}  <br />
需要同步账户数量：${TradeAccountNum}  <br />
已同步完成：${RightAccountNum}  <br />
报错账户信息: ${ErrList}  <br />
`

	lMap := map[string]string{
		"HunterName":      _this.HunterName,
		"NowDir":          DirText,
		"TradeAccountNum": mStr.ToStr(len(ApiKeyList)),
		"RightAccountNum": mStr.ToStr(len(RightAccount)),
		"ErrList":         mJson.ToStr(ErrList),
	}

	Content := mStr.Temp(tmplStr, lMap)
	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "市场方向已改变",
		Title:       "市场方向已改变,所有账户均以同步持仓",
		Content:     Content,
		Description: "同步持仓邮件",
	})

	global.TradeLog.Println(_this.HunterName, "交易失败列表", ErrList)

	_this.CloseOrderSettlement(AccountSettlement)
}

func (_this *HunterObj) CloseOrderSettlement(Settlement []SettlementType) {
	UserOrderArr := []dbType.UserOrderTable{}

	resErr := []string{}
	for _, item := range Settlement {
		if len(item.OkxKey.UserID) > 10 {
			UserOrderArr = append(UserOrderArr, dbType.UserOrderTable{
				OkxPositions:    item.OkxPositions,
				OKXBalance:      item.OKXBalance,
				OkxKey:          item.OkxKey,
				VirtualPosition: item.VirtualPosition,
				UserID:          item.OkxKey.UserID,
				OrderID:         mEncrypt.GetUUID(),
				CreateTime:      mTime.GetUnixInt64(),
			})
		} else {
			resErr = append(resErr, mStr.Join(
				"hunter.CloseOrderSettlement UserID为空",
				"ServeID:", config.AppEnv.ServeID, "<br />",
				"HunterName", _this.HunterName, "<br />",
				"OkxKeyName", item.OkxKey.Name, "<br />",
			))
		}
	}

	if len(resErr) > 0 {
		taskPush.SysEmail(taskPush.SysEmailOpt{
			From:        config.SysName,
			To:          config.NoticeEmail,
			Subject:     "存储用户订单出错",
			Title:       "存储用户订单出错,错误信息如下",
			Content:     mJson.ToStr(resErr),
			Description: "同步持仓邮件",
		})
	}

	db, err := mMongo.New(mMongo.Opt{
		UserName: config.SysEnv.MongoUserName,
		Password: config.SysEnv.MongoPassword,
		Address:  config.SysEnv.MongoAddress,
		DBName:   "Account",
		Timeout:  len(UserOrderArr) * 20,
	}).Connect()
	if err != nil {
		global.LogErr("hunter.CloseOrderSettlement 数据库连接失败", err)
		return
	}
	defer db.Close()
	db.Collection("OkxOrder")

	for _, Kd := range UserOrderArr {
		FK := bson.D{{
			Key:   "CreateTime",
			Value: Kd.CreateTime,
		}}
		UK := bson.D{}
		mStruct.Traverse(Kd, func(key string, val any) {
			UK = append(UK, bson.E{
				Key: "$set",
				Value: bson.D{
					{
						Key:   key,
						Value: val,
					},
				},
			})
		})

		upOpt := options.Update()
		upOpt.SetUpsert(true)
		_, err := db.Table.UpdateOne(db.Ctx, FK, UK, upOpt)
		if err != nil {
			global.LogErr("数据更插失败  %+v", mStr.Join(
				"hunter.CloseOrderSettlement UserID为空",
				"ServeID:", config.AppEnv.ServeID, "<br />",
				"HunterName", _this.HunterName, "<br />",
				"OkxKeyName", Kd.OkxKey.Name, "<br />",
			))
		}
	}
}
