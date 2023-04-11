package analyConfig

import (
	"CoinAI.net/server/okxInfo"
	"CoinAI.net/task/taskStart"
	"CoinAI.net/task/testHunter"
)

type ViolenceOpt struct {
	StartTime int64
	EndTime   int64
	InstID    string
	EmaPArr   []int
	CAPArr    []int
	LevelArr  []int
	CAPMax    []string
	ConfArr   []okxInfo.TradeKdataOpt
	OutPutDir string
}

// 暴力求值
func Violence(opt ViolenceOpt) {
	taskStart.BackTest(taskStart.BackOpt{
		StartTime: opt.StartTime,
		EndTime:   opt.EndTime,
		InstID:    opt.InstID,
		OutPutDir: opt.OutPutDir,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  opt.EmaPArr,
			CAPArr:   opt.CAPArr,
			LevelArr: opt.LevelArr,
			CAPMax:   opt.CAPMax,
			ConfArr:  opt.ConfArr,
		},
	})
}
