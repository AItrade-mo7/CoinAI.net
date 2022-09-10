package ready

import (
	"fmt"

	"CoinAI.net/server/analy"
	"CoinAI.net/server/global"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mClock"
	"github.com/EasyGolang/goTools/mJson"
	"github.com/EasyGolang/goTools/mOKX"
)

func Start() {
	SetMarket()
	go mClock.New(mClock.OptType{
		Func: SetMarket,
		Spec: "30 0,15,30,45 * * * ? ",
	})
}

func SetMarket() {
	err := CheckAccount()
	if err != nil {
		return
	}
	// 获取产品信息
	GetSWAPInst()

	// 获取市场行情
	GetCoinMarket()

	// 筛选最近币种的信息
	RecentTickerList := analy.RecentTicker()

	// 币种历史数据
	okxInfo.AnalyKdata = make(map[string][]mOKX.TypeKd)
	AnalyKdata := make(map[string][]mOKX.TypeKd)
	if len(RecentTickerList) > 3 {
		for _, item := range RecentTickerList {
			// 开始设置 SWAP
			SwapInst := mOKX.TypeInst{}
			for _, SWAP := range okxInfo.SWAP_inst {
				if SWAP.Uly == item.InstID {
					SwapInst = SWAP
					break
				}
			}
			if len(SwapInst.InstID) < 3 {
				continue
			}

			list := GetCoinAnalyKdata(SwapInst.InstID)
			if len(list) == 300 {
				AnalyKdata[SwapInst.InstID] = list
			}
		}
	}
	okxInfo.AnalyKdata = AnalyKdata

	// 根据 振幅 筛选并排序
	okxInfo.HLAnalySelect = []okxInfo.HLAnalySelectType{}
	okxInfo.HLAnalySelect = analy.HLAnalySelect(AnalyKdata)

	// 根据方向涨跌幅度筛选并排序
	analy.DirAnalySelect()

	okxInfo.SetHunterInstID("11") // 暂时写死
}

// 用户信息检查
func CheckAccount() (resErr error) {
	GetUserInfo()
	GetOkxKey()

	mJson.Println(okxInfo.CoinServe)
	mJson.Println(okxInfo.UserInfo)

	resErr = nil
	if len(okxInfo.CoinServe.OkxKeyID) < 10 {
		resErr = fmt.Errorf("读取 dbData.CoinServe 失败 %+v", okxInfo.CoinServe)
		global.LogErr(resErr)
		return
	}

	if len(okxInfo.UserInfo.OkxKeyList) < 1 {
		resErr = fmt.Errorf("读取 dbData.UserInfo 失败 %+v", okxInfo.UserInfo)
		global.LogErr(resErr)
		return
	}

	for _, val := range okxInfo.UserInfo.OkxKeyList {
		if okxInfo.CoinServe.OkxKeyID == val.OkxKeyID {
			okxInfo.OkxKey = val
			break
		}
	}

	if len(okxInfo.OkxKey.OkxKeyID) < 10 {
		resErr = fmt.Errorf("读取 dbData.OkxKey 失败 %+v", okxInfo.OkxKey)
		global.LogErr(resErr)
		return
	}

	return
}
