package config

import "github.com/EasyGolang/goTools/mOKX"

var LeverOpt = []int{2, 3, 4, 5, 6, 7, 8, 9, 10}

func GetOKXKey(Index int) mOKX.TypeOkxKey {
	ReturnKey := mOKX.TypeOkxKey{}
	for key, val := range AppEnv.ApiKeyList {
		if key == Index {
			ReturnKey = val
			break
		}
	}
	return ReturnKey
}
