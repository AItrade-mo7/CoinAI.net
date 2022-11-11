package hunter

import (
	"github.com/EasyGolang/goTools/mOKX"
)

type NewOpt struct {
	StartTime int64 // 开始的时间节点
	EndTime   int64 // 结束的时间节点
}

type NewObj struct {
	StartTime int64 // 开始的时间节点
	BaseList  []mOKX.TypeKd
	RunList   []mOKX.TypeKd
}

func New(opt NewOpt) *NewObj {
	obj := NewObj{}

	if opt.StartTime < 0 { // 如果不存在 StartTime， 则 StartTime = 当前时间
		// BaseList 为当前时间点向前 N 条
	} else {
		// BaseList 为 StartTime 时间点向前 N 条
	}

	if opt.EndTime > 0 { // 如果存在 EndTime ， 则表示历史时间为固定字段
		// 需要在这里回填历史 Start 至 EndTime
	}

	return &obj
}
