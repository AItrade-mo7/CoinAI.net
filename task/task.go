package main

import (
	_ "embed"

	"CoinAI.net/server/global"
	"CoinAI.net/task/taskStart"
	"CoinAI.net/task/testHunter"
	"github.com/EasyGolang/goTools/mTime"
)

var ResultBasePath = "/root/AItrade/CoinAI.net/task/analyConfig/最近8个月2"

func main() {
	// 初始化系统参数
	global.Start()

	BackAnaly()
}

var EmaPArr = []int{
	// 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98,
	// 102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128,
	// 130, 132, 134, 136, 138, 140, 142, 144, 146, 148, 150, 152, 154, 156, 158,
	// 160, 162, 164, 166, 168, 170, 172, 174, 176, 178, 180, 182, 184, 186, 188,
	// 190, 192, 194, 196, 198, 200, 202, 204, 206, 208, 210, 212, 214, 216, 218,
	// 220, 222, 224, 226, 228, 230, 232, 234, 236, 238, 240, 242, 244, 246, 248,
	// 250, 252, 254, 256, 258, 260, 262, 264, 266, 268, 270, 272, 274, 276, 278,
	// 280, 282, 284, 286, 288, 290, 292, 294, 296, 298, 300, 302, 304, 306, 308,
	// 310, 312, 314, 316, 318, 320, 322, 324, 326, 328, 330, 332, 334, 336, 338,
	// 340, 342, 344, 346, 348, 350, 352, 354, 356, 358, 360, 362, 364, 366, 368,
	// 370, 372, 374, 376, 378, 380, 382, 384, 386, 388, 390, 392, 394, 396, 398,
	// 400, 402, 404, 406, 408, 410, 412, 414, 416, 418, 420, 422, 424, 426, 428,
	// 430, 432, 434, 436, 438, 440, 442, 444, 446, 448, 450, 452, 454, 456, 458,
	// 460, 462, 464, 466, 468, 470, 472, 474, 476, 478, 480, 482, 484, 486, 488,
	// 490, 492, 494, 496, 498, 500, 502, 504, 506, 508, 510, 512, 514, 516, 518,
	// 520, 522, 524, 526, 528, 530, 532, 534, 536, 538, 540, 542, 544, 546, 548,
	// 550, 552, 554, 556, 558, 560, 562, 564, 566, 568, 570, 572, 574, 576, 578,
	// 580, 582, 584, 586, 588, 590, 592,
	60, 108, 344, 590,
}

var (
	CAPArr = []int{
		2,
		3,
		// 4,
		// 5,
		// 6,
	}
	LevelArr = []int{1}
	CAPMax   = []string{
		"0.2",
		//  "0.4",
		//  "0.6",
		//  "0.8",
		"1",
	}
)

func BackAnaly() {
	EndTime := mTime.GetUnixInt64()
	StartTime := EndTime - (mTime.UnixTimeInt64.Day * 260)

	InstID := "BTC-USDT"
	BTCResult := taskStart.BackTest(taskStart.BackOpt{
		StartTime: StartTime,
		EndTime:   EndTime,
		InstID:    InstID,
		GetConfigOpt: testHunter.GetConfigOpt{
			EmaPArr:  EmaPArr,
			CAPArr:   CAPArr,
			LevelArr: LevelArr,
			CAPMax:   CAPMax,
		},
	})
	BTCResult.ResultBasePath = ResultBasePath
	// analyConfig.GetWinArr(
	// 	BTCResult,
	// 	// taskStart.BackReturn{
	// 	// 	InstID:         InstID,
	// 	// 	BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
	// 	// 	ResultBasePath: ResultBasePath,
	// 	// },
	// )

	// InstID = "ETH-USDT"
	// ETHResult := taskStart.BackTest(taskStart.BackOpt{
	// 	StartTime: StartTime,
	// 	EndTime:   EndTime,
	// 	InstID:    InstID,
	// 	GetConfigOpt: testHunter.GetConfigOpt{
	// 		EmaPArr:  EmaPArr,
	// 		CAPArr:   CAPArr,
	// 		LevelArr: LevelArr,
	// 		CAPMax:   CAPMax,
	// 	},
	// })
	// ETHResult.ResultBasePath = ResultBasePath
	// analyConfig.GetWinArr(
	// ETHResult,
	// taskStart.BackReturn{
	// 	InstID:         InstID,
	// 	BillingPath:    mStr.Join(ResultBasePath, "/", InstID, "-BillingArr.json"),
	// 	ResultBasePath: ResultBasePath,
	// },
	// )
}
