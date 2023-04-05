package hunter

import (
	"fmt"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mStr"
	"github.com/EasyGolang/goTools/mTime"
)

func (_this *HunterObj) Start() {
	go func() {
		for ok := range okxInfo.ReadyTicking {
			global.TradeLog.Println(_this.HunterName, " ===== hunter.Start 执行 ===== ", mTime.UnixFormat(mTime.GetUnixInt64()), ok)
			_this.Running()
		}
	}()
}

func (_this *HunterObj) Running() {
	global.TradeLog.Println(" === hunter.Running === ", _this.KdataInst.InstID)

	if len(_this.KdataInst.InstID) < 2 || len(_this.TradeInst.InstID) < 2 {
		global.LogErr(_this.HunterName, "hunter.Running", "okxInfo.TradeInst.InstID 或 KdataInst 为空")
		return
	}

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

	// err = FormatTradeKdata(TradeKdataOpt{
	// 	MA_Period:      108,
	// 	RSI_Period:     18,
	// 	RSI_EMA_Period: 14,
	// 	CAP_Period:     3,
	// })
	// if err != nil { // 这里参数出了问题
	// 	global.LogErr(err)
	// 	_this.Running() // 立即重新执行一次 Running
	// 	return
	// }

	// Analy()
}

func (_this *HunterObj) FileBaseKdata() error {
	Page := 4 // 如果数组为空，则填充 500 条进去
	if len(_this.NowKdataList) < 100 {
		// 回填历史数据 1 组
		for i := Page; i >= 0; i-- {
			time.Sleep(time.Second / 3)
			List := mOKX.GetKdata(mOKX.GetKdataOpt{
				InstID: _this.KdataInst.InstID,
				Page:   i,
				After:  mTime.GetUnixInt64(),
			})
			_this.NowKdataList = append(_this.NowKdataList, List...)
		}
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
