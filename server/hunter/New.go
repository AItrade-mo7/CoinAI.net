package hunter

import (
	"fmt"

	"CoinAI.net/server/global"
	"CoinAI.net/server/okxApi"
	"CoinAI.net/server/okxInfo"
	"github.com/EasyGolang/goTools/mOKX"
	"github.com/EasyGolang/goTools/mTime"
)

type NewOpt struct {
	StartTime int64 // 开始的时间节点
	EndTime   int64 // 结束的时间节点
	InstID    string
}

type NewObj struct {
	StartTime int64 // 开始的时间节点
	BaseList  []mOKX.TypeKd
	RunList   []mOKX.TypeKd
	Inst      mOKX.TypeInst
}

func New(opt NewOpt) *NewObj {
	obj := NewObj{}

	if len(opt.InstID) < 3 {
		global.LogErr("hunter.New", obj.Inst.InstID)
		return nil
	}

	obj.Inst = okxInfo.Inst[opt.InstID]
	if len(obj.Inst.InstID) < 3 {
		global.LogErr("hunter.New", obj.Inst.InstID)
		return nil
	}

	if opt.StartTime > 0 { // 如果不存在 StartTime， 则 StartTime = 当前时间
		obj.StartTime = opt.StartTime
	} else {
		obj.StartTime = mTime.GetUnixInt64()
	}
	obj.GetBaseKdata()

	if opt.EndTime > 0 { // 如果存在 EndTime ， 则表示历史时间为固定字段
		// 需要在这里回填历史 Start 至 EndTime
	} else {
		// 需要在这里启动 wss 监听
	}

	return &obj
}

func (_this *NewObj) GetBaseKdata() {
	List := okxApi.GetKdata(okxApi.GetKdataOpt{
		InstID: _this.Inst.InstID,
		After:  _this.StartTime,
	})

	fmt.Println(List)

	if len(List) > 0 {
		fmt.Println(len(List),
			List[0].TimeStr,
			List[len(List)-1].TimeStr,
		)
	}
}

func GetHistoryKdata() {
}
