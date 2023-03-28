package hunter

import (
	"fmt"
	"time"

	"CoinAI.net/server/global"
	"CoinAI.net/server/global/config"
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/server/utils/taskPush"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

func Start() {
	for ok := range okxInfo.HunterTicking {
		global.TradeLog.Println(" ===== hunter.Start 执行 ===== ", mTime.UnixFormat(mTime.GetUnixInt64()), ok)
		Running()
	}
}

func Running() {
	global.TradeLog.Println(" === hunter.Running === ", okxInfo.KdataInst.InstID)

	err := SetTradeInst() // 设置一下
	if err != nil {
		global.LogErr(err)
		return
	}

	if len(okxInfo.KdataInst.InstID) < 2 || len(okxInfo.TradeInst.InstID) < 2 {
		global.LogErr("hunter.Running", "okxInfo.TradeInst.InstID 或 KdataInst 为空")
		return
	}

	err = FileBaseKdata()
	if err != nil { // 在这里切换了币种，重新执行
		Running() // 立即重新执行一次 Running
		return
	}

	// SetNowKdata()

	// FormatTradeKdata()

	// Analy()
}

func FileBaseKdata() error {
	Page := 4 // 如果数组为空，则填充 500 条进去
	if len(okxInfo.NowKdataList) < 100 {
		// 回填历史数据 1 组
		for i := Page; i >= 0; i-- {
			time.Sleep(time.Second / 3)
			List := mOKX.GetKdata(mOKX.GetKdataOpt{
				InstID: okxInfo.KdataInst.InstID,
				Page:   i,
				After:  mTime.GetUnixInt64(),
			})
			okxInfo.NowKdataList = append(okxInfo.NowKdataList, List...)
		}
		global.TradeLog.Println("基础数据回填完毕", len(okxInfo.NowKdataList))
		return nil
	} else { // 如果不为空 则检查当前的数组和持仓币种的关系
		// 在这里执行重启
		if okxInfo.KdataInst.InstID != okxInfo.NowKdataList[len(okxInfo.NowKdataList)-1].InstID {
			okxInfo.NowKdataList = []mOKX.TypeKd{} // 清空历史数据
			warnStr := "即将切换监听币种为: " + okxInfo.KdataInst.InstID
			global.TradeLog.Println(warnStr)
			SendEmail(warnStr)
			return fmt.Errorf(warnStr)
		}
	}

	return nil
}

func SendEmail(Message string) {
	taskPush.SysEmail(taskPush.SysEmailOpt{
		From:        config.SysName,
		To:          config.NoticeEmail,
		Subject:     "监听切换通知",
		Title:       config.SysName + " 币种监听切换",
		Content:     Message,
		Description: "监听切换通知",
	})
}
