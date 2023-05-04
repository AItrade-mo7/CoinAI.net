package hunter

import (
	"github.com/EasyGolang/goTools/mTime"
)

func IsAnalyTimeScale(KTime int64) bool {
	nowTimeObj := mTime.MsToTime(KTime, "0")

	Minute := nowTimeObj.Minute()

	isIn := false
	timeScale := []int{1, 16, 31, 46}
	for _, val := range timeScale {
		if Minute-val == 0 {
			isIn = true
			break
		}
	}
	return isIn
}
