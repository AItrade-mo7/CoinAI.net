package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mMongo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (_this *HunterObj) Start() {
	go _this.Running()
}

func (_this *HunterObj) Running() {
	global.TradeLog.Println(_this.HunterName, " === hunter.Running === ", _this.KdataInst.InstID)

	// RoundNum := mCount.GetRound(0, 60) // 延迟随机秒数
	// time.Sleep(time.Second * time.Duration(RoundNum))

	// 选取K线和合约信息
	if len(_this.KdataInst.InstID) < 2 || len(_this.TradeInst.InstID) < 2 {
		err := _this.SetTradeInst(_this.InstID)
		if err != nil {
			global.LogErr(err)
			return
		}
		_this.Running()
		return
	}

	// 在这里填充基础数据 走 mongodb
	err := _this.FileBaseKdata()
	if err != nil { // 在这里切换了币种，重新执行
		_this.Running() // 立即重新执行一次 Running
		return          // 阻断后面的执行
	}

	err = _this.SetNowKdata()
	if err != nil { // 在这里检查数据出了问题
		global.LogErr(err)
		_this.Running() // 立即重新执行一次 Running
		return
	}

	err = _this.SetTradeConfig()
	if err != nil { // 在这里检查数据出了问题
		global.LogErr(err)
		_this.Running() // 立即重新执行一次 Running
		return
	}

	_this.SyncInfoToGlobal() // 同步一次数据

	// 策略执行的核心
	err = _this.FormatTradeKdata()
	if err != nil { // 这里参数出了问题
		global.LogErr(err)
		_this.Running() // 立即重新执行一次 Running
		return
	}
	_this.Analy()
	// 策略执行的核心模块

	_this.SyncInfoToGlobal() // 同步一次数据
}

func (_this *HunterObj) FileBaseKdata() error {
	if len(_this.NowKdataList) < 100 {
		// 回填历史数据 1 组
		db := mMongo.New(mMongo.Opt{
			UserName: config.SysEnv.MongoUserName,
			Password: config.SysEnv.MongoPassword,
			Address:  config.SysEnv.MongoAddress,
			DBName:   "CoinMarket",
			Timeout:  _this.MaxLen,
		}).Connect().Collection(_this.KdataInst.InstID)
		defer db.Close()
		findOpt := options.Find()
		findOpt.SetSort(map[string]int{
			"TimeUnix": -1,
		})
		findOpt.SetAllowDiskUse(true)
		findOpt.SetLimit(int64(_this.MaxLen))
		cur, err := db.Table.Find(db.Ctx, bson.D{}, findOpt)
		if err != nil {
			db.Close()
			return err
		}
		AllList := []mOKX.TypeKd{}
		for cur.Next(db.Ctx) {
			var result mOKX.TypeKd
			cur.Decode(&result)
			AllList = append(AllList, result)
		}
		db.Close()

		KdataList := []mOKX.TypeKd{}
		for i := len(AllList) - 1; i >= 0; i-- {
			el := AllList[i]
			KdataList = append(KdataList, el)
		}
		_this.NowKdataList = KdataList

		Last := _this.NowKdataList[len(_this.NowKdataList)-1]
		global.TradeLog.Println(_this.HunterName, "基础数据回填完毕", len(_this.NowKdataList), Last.TimeStr, Last.InstID)
		return nil
	} else { // 如果不为空 则检查当前的数组和持仓币种的关系
		// 在这里执行重启
		if _this.KdataInst.InstID != _this.NowKdataList[len(_this.NowKdataList)-1].InstID {
			_this.NowKdataList = []mOKX.TypeKd{} // 清空历史数据
			warnStr := _this.HunterName + "即将切换监听币种为: " + _this.KdataInst.InstID
			global.TradeLog.Println(warnStr)
			_this.SendEmail(warnStr)
			return fmt.Errorf(warnStr)
		}
	}

	return nil
}

func (_this *HunterObj) SendEmail(Message string) {
	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     mStr.Join(_this.HunterName, "币种监听切换通知"),
		Title:       mStr.Join(_this.HunterName, "币种监听切换"),
		Content:     Message,
		Description: "监听切换通知",
	})
}

func (_this *HunterObj) SyncInfoToGlobal() {
	Name := _this.HunterName

	HunterData := okxInfo.HunterData{
		HunterName:         _this.HunterName,         // 策略的名字
		Describe:           _this.Describe,           // 描述
		InstID:             _this.InstID,             // 当前策略主打币种
		TradeInst:          _this.TradeInst,          // 交易的 InstID SWAP
		KdataInst:          _this.KdataInst,          // K线的 InstID SPOT
		NowKdataList:       _this.NowKdataList,       // 现货的原始K线
		TradeKdataList:     _this.TradeKdataList,     // 计算好各种指标之后的K线
		TradeKdataOpt:      _this.TradeKdataOpt,      // 当前参数
		NowVirtualPosition: _this.NowVirtualPosition, // 当前的虚拟持仓
	}
	okxInfo.NowHunterData[Name] = HunterData
}
