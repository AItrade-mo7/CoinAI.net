package ready

import (
	"fmt"

	"CoinAI.net/server/utils/reqDataCenter"
	jsoniter "github.com/json-iterator/go"
)

type PingDataType struct {
	Code int `json:"Code"`
	Data struct {
		IP string `json:"IP"`
	} `json:"Data"`
	Msg string `json:"Msg"`
}

func GetIP() {
	resData, _ := reqDataCenter.NewRest(reqDataCenter.RestOpt{
		Origin: "https://trade-api.mo7.cc",
		Path:   "/ping",
		Method: "GET",
	})

	var PingData PingDataType
	jsoniter.Unmarshal(resData, &PingData)

	LocalIP := PingData.Data.IP

	fmt.Println(LocalIP)
}
