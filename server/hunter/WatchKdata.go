package hunter

import "github.com/EasyGolang/goTools/mOKX"

func WatchKdata() {
	GetBaseKdata()
}

func GetBaseKdata() []mOKX.TypeKd {
	KdataList := []mOKX.TypeKd{}

	for i := 2; i >= 0; i-- {
		// fmt.Println(i)
		// List := okxApi.GetKdata(okxApi.GetKdataOpt{
		// 	InstID:  Inst.InstID,
		// 	After:   StartTime,
		// 	Current: i,
		// })

		// KdataList = append(KdataList, List...)
	}

	return KdataList
}
